package common

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
	"github.com/PDeXchange/pac/test/e2e/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping Service", func() {
	Context("GET /ping", func() {
		It("should return status code 200", func() {
			// parse the url
			fullURL, err := helpers.GetParsedURL(config.Current.PacGoServer + "/ping")
			Expect(err).ToNot(HaveOccurred())

			fmt.Println(fullURL.String())

			// Make http GET call to the required url path
			resp, err := http.Get(fullURL.String())
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Read the body
			body, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			// Check that the body contains the expected fields
			Expect(string(body)).To(MatchJSON(`{
        		"message": "pong"
    		}`))
		})
	})

	Context("GET /pong", func() {
		It("should return status code 404", func() {
			// parse the url
			fullURL, err := helpers.GetParsedURL(config.Current.PacGoServer + "/pong")
			Expect(err).ToNot(HaveOccurred())

			// Make http GET call to the required url path
			resp, err := http.Get(fullURL.String())
			Expect(err).ToNot(HaveOccurred())

			// Get 404 status code
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
