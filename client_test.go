// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	const (
		username = "user"
		password = "p4s5w07d"
	)

	var (
		server *ghttp.Server
		c      *Client
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		c, _ = NewClient(username, password, WithBaseURL(server.URL()))
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when response cannot be unmarshalled", func() {
		It("should return an error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/group"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWithJSONEncoded(http.StatusOK, "foo"),
				),
			)
			gg, err := c.GetGroups()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("EOF"))
			Expect(gg).To(BeNil())
		})
	})

	Context("when group list is requested", func() {
		It("should return group names and no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/group"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `
						<directory count="3">
							<entry name="foo"/>
							<entry name="bar"/>
							<entry name="baz"/>
						</directory>`),
				),
			)
			gg, err := c.GetGroups()
			Expect(err).ToNot(HaveOccurred())
			Expect(gg).To(Equal([]string{
				"foo",
				"bar",
				"baz",
			}))
		})
	})

	Context("when foo group is requested", func() {
		It("should return group members and no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/group/foo"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `
						<group>
							<title>foo</title>
							<person>
								<person userid="foo-member" />
								<person userid="bar-member" />
								<person userid="baz-member" />
							</person>
						</group>`),
				),
			)
			g, err := c.GetGroup("foo")
			Expect(err).ToNot(HaveOccurred())
			Expect(g.Members).To(Equal([]GroupMember{
				{"foo-member"},
				{"bar-member"},
				{"baz-member"},
			}))
		})
	})

	Context("when foo@bar.org user is being looked up", func() {
		It("should return user for this email and no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/search/person", "match=@email='foo@bar.org'"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `
						<collection matches="1">
							<person>
									<login>foo-member</login>
									<email>foo@bar.org</email>
									<realname>Foo Bar</realname>
									<state>confirmed</state>
									<watchlist>
										<project name="project-1" />
										<project name="project-2" />
									</watchlist>
							</person>
						</collection>`),
				),
			)
			u, err := c.GetUserByEmail("foo@bar.org")
			Expect(err).ToNot(HaveOccurred())
			Expect(u.ID).To(Equal("foo-member"))
			Expect(u.Watchlist).To(Equal([]WatchlistEntry{
				{"project-1"},
				{"project-2"},
			}))
		})
	})

	Context("when non-existing user email is being looked up", func() {
		It("should return error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/search/person", "match=@email='qux@bar.org'"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `
						<collection matches="0">
						</collection>`),
				),
			)
			u, err := c.GetUserByEmail("qux@bar.org")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("no users found with email"))
			Expect(u).To(BeNil())
		})
	})

	Context("when a new group is being created", func() {
		It("should return no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPut, "/group/test-group"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.VerifyBody([]byte("<group><title>test-group</title><person></person></group>")),
					ghttp.RespondWith(http.StatusOK, `
						<status code="ok">
							<summary>Ok</summary>
						</status>`),
				),
			)
			err := c.NewGroup("test-group")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when an existing group is being deleted", func() {
		It("should return no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodDelete, "/group/test-group"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `<status code="ok">
							<summary>Ok</summary>
						</status>`),
				),
			)
			err := c.DeleteGroup("test-group")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when a non-existing group is being deleted", func() {
		It("should return error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodDelete, "/group/test-non-group"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusNotFound, `
						<status code="not_found">
							<summary>Couldn't find Group 'test-non-group'</summary>
						</status>`),
				),
			)
			err := c.DeleteGroup("test-non-group")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HaveSuffix("Couldn't find Group 'test-non-group'"))
		})
	})

	Context("when a user is being added to a group", func() {
		It("should return no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/group/test-group", "cmd=add_user&userid=test-user"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.RespondWith(http.StatusOK, `<status code="ok">
							<summary>Ok</summary>
						</status>`),
				),
			)
			err := c.AddGroupMember("test-group", "test-user")
			Expect(err).ToNot(HaveOccurred())
		})
	})


})
