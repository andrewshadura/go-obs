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
	Context("when a user profile is marshalled", func() {
		It("should produce a valid XML", func() {
			d := User{
				ID:       "user",
				Email:    "user@host.co.uk",
				Realname: "User Person",
				State:    "confirmed",
				Watchlist: []ProjectRef{
					{"project-1"},
					{"project-2"},
				},
			}
			data, err := xml.Marshal(d)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(unindent(`
				<person>
						<login>user</login>
						<email>user@host.co.uk</email>
						<realname>User Person</realname>
						<state>confirmed</state>
						<watchlist>
							<project name="project-1"></project>
							<project name="project-2"></project>
						</watchlist>
				</person>`)))
		})
	})
})

var _ = Describe("Users", func() {
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

	When("foo@bar.org user is being looked up", func() {
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
			Expect(u.Watchlist).To(Equal([]ProjectRef{
				{"project-1"},
				{"project-2"},
			}))
		})
	})

	When("non-existing user email is being looked up", func() {
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

})
