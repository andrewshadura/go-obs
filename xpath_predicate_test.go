// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Escaping according to XPath rules", func() {
	When("a value without apostrophes is being escaped", func() {
		It("is enclosed in apostrophes", func() {
			p := XPathAttrEquals("key", "value")
			Expect(p.String()).To(Equal("@key='value'"))
		})
	})

	When("a value with apostrophes is being escaped", func() {
		It("is enclosed in double quotes", func() {
			p := XPathAttrEquals("key", "va'lue")
			Expect(p.String()).To(Equal("@key=\"va'lue\""))
		})
	})
})
