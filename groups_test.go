// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/xml"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Marshalling", func() {
	When("a group directory is marshalled", func() {
		It("should produce a valid XML", func() {
			d := directory{
				Entries: []directoryEntry{
					{"foo"},
					{"bar"},
					{"baz"},
				},
			}
			data, err := xml.Marshal(d)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(unindent(`
				<directory>
						<entry name="foo"></entry>
						<entry name="bar"></entry>
						<entry name="baz"></entry>
				</directory>`)))
		})
	})

	When("a group object is marshalled", func() {
		It("should produce a valid XML", func() {
			g := Group{
				ID: "foo",
				Members: []UserRef{
					{"foo-member"},
					{"bar-member"},
					{"baz-member"},
				},
			}
			data, err := xml.Marshal(g)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(unindent(`
				<group>
						<title>foo</title>
						<person>
							<person userid="foo-member"></person>
							<person userid="bar-member"></person>
							<person userid="baz-member"></person>
						</person>
				</group>`)))
		})
	})
})

var _ = Describe("Groups", func() {
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

	Describe("listing", func() {
		When("a response cannot be unmarshalled", func() {
			It("should return an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/group"),
						ghttp.VerifyBasicAuth(username, password),
						ghttp.RespondWithJSONEncoded(http.StatusOK, "foo"),
					),
				)
				gg, err := c.ListGroups()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("EOF"))
				Expect(gg).To(BeNil())
			})
		})

		When("an existing group is requested", func() {
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
				gg, err := c.ListGroups()
				Expect(err).ToNot(HaveOccurred())
				Expect(gg).To(Equal([]string{
					"foo",
					"bar",
					"baz",
				}))
			})
		})

	})

	When("group 'foo' is requested", func() {
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
			Expect(g.Members).To(Equal([]UserRef{
				{"foo-member"},
				{"bar-member"},
				{"baz-member"},
			}))
		})
	})

	When("a new group is being created", func() {
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

	When("a group is being updated", func() {
		It("should return no error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPut, "/group/foo"),
					ghttp.VerifyBasicAuth(username, password),
					ghttp.VerifyBody([]byte(`<group><title>foo</title><person><person userid="foo-member"></person><person userid="bar-member"></person><person userid="baz-member"></person></person></group>`)),
					ghttp.RespondWith(http.StatusOK, `<status code="ok">
							<summary>Ok</summary>
						</status>`),
				),
			)
			g := Group{
				ID: "foo",
				Members: []UserRef{
					{"foo-member"},
					{"bar-member"},
					{"baz-member"},
				},
			}
			err := c.UpdateGroup(&g)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("an existing group is being deleted", func() {
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

	When("a non-existing group is being deleted", func() {
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

	When("a user is being added to a group", func() {
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
