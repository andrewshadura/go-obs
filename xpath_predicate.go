// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"strings"
)

type XPathPredicate struct {
	path     string
	operator string
	value    string
}

func XPathAttrEquals(name string, value string) *XPathPredicate {
	return &XPathPredicate{
		path:     "@" + name,
		operator: "=",
		value:    value,
	}
}

func (p *XPathPredicate) String() string {
	if strings.ContainsAny(p.value, "'") {
		return p.path + p.operator + "\"" + p.value + "\""
	} else {
		return p.path + p.operator + "'" + p.value + "'"
	}
}
