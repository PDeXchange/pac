package common

import (
	"fmt"
	"testing"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommon(t *testing.T) {
	// Run the Ginkgo tests
	RegisterFailHandler(Fail)
	fmt.Fprintf(GinkgoWriter, "Starting Common e2e test suite\n")
	RunSpecs(t, "Common e2e Suite")
}

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, ">>> Entered Common E2E BeforeSuite <<<\n")
	Expect(config.LoadConfig("../../config/config.yaml", "test")).To(Succeed())
})
