// Copyright (C) 2021, Andrej Shadura
// Copyright (C) 2021, Collabora Limited
//
// SPDX-License-Identifier: Apache-2.0

package obs_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

// TODO: Replace this with --junit-report= when migrating to ginkgo v2:
func configureJUnitXML() []Reporter {
	if junit_xml, ok := os.LookupEnv("JUNIT_XML"); ok {
		if !strings.Contains(junit_xml, "%d") {
			if config.GinkgoConfig.ParallelTotal == 1 {
				return []Reporter{
					reporters.NewJUnitReporter(junit_xml),
				}
			} else {
				parts := strings.Split(junit_xml, ".")
				if len(parts) < 2 {
					parts[len(parts)-1] = parts[len(parts)-1] + "_%d"
				} else {
					parts[len(parts)-2] = parts[len(parts)-2] + "_%d"
				}
				junit_xml = strings.Join(parts, ".")
			}
		}
		return []Reporter{
			reporters.NewJUnitReporter(fmt.Sprintf(junit_xml, config.GinkgoConfig.ParallelNode)),
		}
	}
	return []Reporter{}
}

func TestObs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t, "Obs Suite", configureJUnitXML())
}
