package common

import (
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping Service", func() {
	Context("GET /ping", func() {
		It("should return status code 200", func() {
			resp, err := http.Get("http://localhost:8000/ping")
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

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
			resp, err := http.Get("http://localhost:8000/pong")
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
