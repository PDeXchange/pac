package tnc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
	"github.com/PDeXchange/pac/test/e2e/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	dbCollectionName = "tnc"
)

var _ = Describe("GET TnC (/api/v1/tnc)", func() {
	var mongoClient *helpers.MongoClient
	var token, userID string
	var err error
	var url *url.URL

	BeforeEach(func() {
		// parse the url
		url, err = helpers.GetParsedURL(config.Current.PacGoServer + "/api/v1/tnc")
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

		// get mongo client object
		mongoClient, err = helpers.GetMongoClient()
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(mongoClient).NotTo(BeNil())

		// drop the collection: tnc if already exists
		Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())

		// fetch the token for the test user
		token, err = helpers.GetTestUserToken()
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(token).NotTo(BeEmpty())

		// fetch the id of the user from the token
		userID, err = helpers.GetUserID(token)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(userID).NotTo(BeEmpty())
	})

	AfterEach(func() {
		// Disconnect the mongoClient
		Expect(mongoClient.Disconnect()).To(Succeed())
	})

	When("TnC exists in DB", func() {
		BeforeEach(func() {
			// Insert a record into tnc collection
			Expect(mongoClient.InsertOne(dbCollectionName, map[string]any{
				"user_id":     userID,
				"accepted":    true,
				"accepted_at": time.Now(),
			})).To(Succeed())
		})

		AfterEach(func() {
			// Drop tnc collection
			Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())
		})

		It("should return status code 200 with user tnc info", func() {
			client := &http.Client{}

			// make GET request
			req, err := http.NewRequest("GET", url.String(), nil)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Send request
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
			defer resp.Body.Close()

			// Read response body
			respBody, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// unmarshal response body to result map
			var result map[string]any
			Expect(json.Unmarshal(respBody, &result)).To(Succeed())

			// verify each fields in the response body
			Expect(result["user_id"]).To(Equal(userID))
			Expect(result["accepted"]).To(BeTrue())
			Expect(result["accepted_at"]).NotTo(BeEmpty())
		})
	})

	When("TnC does not exists in DB", func() {
		AfterEach(func() {
			// Drop tnc collection
			Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())
		})

		It("should return status code 200 with user tnc info containing only userID", func() {
			client := &http.Client{}

			// make GET request
			req, err := http.NewRequest("GET", url.String(), nil)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Send request
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
			defer resp.Body.Close()

			// Read the response body
			respBody, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// unmarshal response body to result map
			var result map[string]any
			Expect(json.Unmarshal(respBody, &result)).To(Succeed())

			// verify each fields in the response body
			Expect(result["user_id"]).To(Equal(userID))
			Expect(result["accepted"]).To(BeFalse())
		})
	})
})

var _ = Describe("Accept TnC (/api/v1/tnc)", func() {
	var mongoClient *helpers.MongoClient
	var token, userID string
	var err error
	var url *url.URL

	BeforeEach(func() {
		// parse the url
		url, err = helpers.GetParsedURL(config.Current.PacGoServer + "/api/v1/tnc")
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

		// get mongo client object
		mongoClient, err = helpers.GetMongoClient()
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(mongoClient).NotTo(BeNil())

		// drop the collection: tnc if already exists
		Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())

		// fetch the token for the test user
		token, err = helpers.GetTestUserToken()
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(token).NotTo(BeEmpty())

		// fetch the id of the user from the token
		userID, err = helpers.GetUserID(token)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
		Expect(userID).NotTo(BeEmpty())
	})

	AfterEach(func() {
		// Disconnect the mongoClient
		Expect(mongoClient.Disconnect()).ToNot(HaveOccurred())
	})

	When("TnC already present in the DB", func() {
		Context("TnC is not Accepted", func() {
			BeforeEach(func() {
				// Insert the record having only userID into tnc collection
				Expect(mongoClient.InsertOne(dbCollectionName, map[string]any{
					"user_id": userID,
				})).To(Succeed())
			})

			AfterEach(func() {
				// Drop tnc collection
				Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())
			})

			It("should return status code 201 on successful acceptance", func() {
				client := &http.Client{}

				// Do POST request
				req, err := http.NewRequest("POST", url.String(), nil)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

				// Set headers
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				// Send request
				resp, err := client.Do(req)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				// Fetch the record from tnc collection
				dbRecord, err := mongoClient.FindOne(dbCollectionName)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

				// Check that the db record contains the expected fields
				Expect(dbRecord["user_id"]).To(Equal(userID))
				Expect(dbRecord["accepted"]).To(Equal(true))
				Expect(dbRecord["accepted_at"]).NotTo(BeNil())
			})
		})

		Context("TnC is Accepted", func() {
			BeforeEach(func() {
				// Insert the record into tnc collection
				Expect(mongoClient.InsertOne(dbCollectionName, map[string]any{
					"user_id":     userID,
					"accepted":    true,
					"accepted_at": time.Now(),
				})).To(Succeed())
			})

			AfterEach(func() {
				// Drop tnc collection
				Expect(mongoClient.DropCollection(dbCollectionName)).To(Succeed())
			})

			It("should return status code 400", func() {
				client := &http.Client{}

				// Do POST request
				req, err := http.NewRequest("POST", url.String(), nil)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

				// Set headers
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				// Send request
				resp, err := client.Do(req)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
				defer resp.Body.Close()

				respBody, err := io.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				// unmarshal response body to result map
				var result map[string]any
				Expect(json.Unmarshal(respBody, &result)).To(Succeed())

				// Check that the response body contains the error field
				Expect(result["error"]).To(Equal("terms and conditions already accepted"))
			})
		})

	})

	When("TnC does not exists in DB", func() {
		AfterEach(func() {
			Expect(mongoClient.DropCollection(dbCollectionName)).ToNot(HaveOccurred())
		})

		It("should return status code 200 with user tnc stored in DB", func() {
			client := &http.Client{}

			req, err := http.NewRequest("POST", url.String(), nil)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Send request
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))
			defer resp.Body.Close()

			dbRecord, err := mongoClient.FindOne(dbCollectionName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("unexpected error: %v", err))

			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			// Check that the body contains the expected fields
			Expect(dbRecord["user_id"]).To(Equal(userID))
			Expect(dbRecord["accepted"]).To(Equal(true))
			Expect(dbRecord["accepted_at"]).NotTo(BeNil())
		})
	})
})
