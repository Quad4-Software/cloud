// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}
