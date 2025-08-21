package common

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	// Run the Ginkgo tests
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting Common e2e test suite\n")
	RunSpecs(t, "Common Test e2e Suite")
}
