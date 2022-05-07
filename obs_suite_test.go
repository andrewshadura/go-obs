// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestObs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Obs Suite")
}
