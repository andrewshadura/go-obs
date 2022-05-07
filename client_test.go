// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/ghttp"
)

const (
	username = "user"
	password = "p4s5w07d"
)

func unindent(s string) string {
        lines := strings.Split(s, "\n")
        output := strings.Builder{}
        for _, l := range lines {
                output.WriteString(strings.TrimLeft(l, " \t"))
        }
        return output.String()
}

var _ = Describe("Client", func() {
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
		_ = c
	})

})
