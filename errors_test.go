// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"encoding/xml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshalling", func() {
	When("an error object is marshalled", func() {
		It("should produce a valid XML", func() {
			g := ErrorResponse{
				Message: "Couldn't find group 'group'",
				Code:    "not_found",
			}
			data, err := xml.Marshal(g)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(data)).To(Equal(unindent(`
			<status code="not_found">
					<summary>Couldn&#39;t find group &#39;group&#39;</summary>
			</status>`)))
		})
	})
})
