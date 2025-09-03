package tnc

import (
	"fmt"
	"testing"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTnC(t *testing.T) {
	// Run the Ginkgo tests
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting TnC e2e test suite\n")
	RunSpecs(t, "TnC e2e Suite")
}

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, ">>> Entered TnC E2E BeforeSuite <<<\n")
	Expect(config.LoadConfig("../../config/config.yaml", "test")).To(Succeed())
})
