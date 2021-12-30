// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/xml"

	. "github.com/andrewshadura/go-obs/unindent"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshalling", func() {
	Context("when a user profile is marshalled", func() {
		It("should produce a valid XML", func() {
			d := User{
				ID:       "user",
				Email:    "user@host.co.uk",
				Realname: "User Person",
				State:    "confirmed",
				Watchlist: []WatchlistEntry{
					{"project-1"},
					{"project-2"},
				},
			}
			data, err := xml.Marshal(d)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(Unindent(`
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
