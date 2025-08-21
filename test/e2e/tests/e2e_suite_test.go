package tests

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	// Run the Ginkgo tests
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting E2E test suite\n")
	RunSpecs(t, "E2E Suite")
}
