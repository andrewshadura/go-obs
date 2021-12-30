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
	Context("when a group directory is marshalled", func() {
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
			Expect(string(data)).To(Equal(Unindent(`
				<directory>
						<entry name="foo"></entry>
						<entry name="bar"></entry>
						<entry name="baz"></entry>
				</directory>`)))
		})
	})

	Context("when a group object is marshalled", func() {
		It("should produce a valid XML", func() {
			g := group{Group{
				ID: "foo",
				Members: []UserRef{
					{"foo-member"},
					{"bar-member"},
					{"baz-member"},
				}},
			}
			data, err := xml.Marshal(g)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(Unindent(`
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
