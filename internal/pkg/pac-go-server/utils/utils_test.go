package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type SampleRequest struct {
	Name   string            `json:"name" binding:"required,max=5"`
	Age    int               `json:"age" binding:"gte=18,lte=60"`
	Active bool              `json:"active"`
	Tags   []string          `json:"tags"`
	Meta   map[string]string `json:"meta"`
	ID     string            `json:"id" binding:"uuid"`
	Note   chan int          `json:"note"` // unsupported type to hit getExpectedType "value"
}

func performRequest(body string) (*httptest.ResponseRecorder, bool, error) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var req SampleRequest
	ok, err := BindAndValidate(c, &req)
	return w, ok, err
}

func TestBindAndValidate_AllBranches(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name           string
		body           string
		expectOK       bool
		expectHTTPCode int
		expectField    string
		expectError    string
	}{
		// --- success case ---
		{
			name:           "success",
			body:           `{"name":"John","age":30,"active":true,"tags":["go"],"meta":{"k":"v"},"id":"` + validUUID + `"}`,
			expectOK:       true,
			expectHTTPCode: 200,
		},
		// --- translateValidationError branches ---
		{
			name:           "required field missing",
			body:           `{"age":30,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "Name",
			expectError:    "is required",
		},
		{
			name:           "max length exceeded",
			body:           `{"name":"TooLongName","age":30,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "Name",
			expectError:    "must not be longer than",
		},
		{
			name:           "age less than gte",
			body:           `{"name":"John","age":10,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "Age",
			expectError:    "greater than or equal",
		},
		{
			name:           "age greater than lte",
			body:           `{"name":"John","age":100,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "Age",
			expectError:    "less than or equal",
		},
		{
			name:           "invalid UUID",
			body:           `{"name":"John","age":30,"id":"not-a-uuid"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "ID",
			expectError:    "must be a valid UUID",
		},
		// --- getExpectedType branches ---
		{
			name:           "type mismatch number",
			body:           `{"name":"John","age":"not-a-number","id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "age",
			expectError:    "must be number",
		},
		{
			name:           "type mismatch boolean",
			body:           `{"name":"John","age":30,"active":"yes","id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "active",
			expectError:    "must be boolean",
		},
		{
			name:           "type mismatch string",
			body:           `{"name":123,"age":30,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "name",
			expectError:    "must be string",
		},
		{
			name:           "type mismatch array",
			body:           `{"name":"John","age":30,"tags":"not-an-array","id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "tags",
			expectError:    "must be array",
		},
		{
			name:           "type mismatch object",
			body:           `{"name":"John","age":30,"meta":"not-an-object","id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "meta",
			expectError:    "must be object",
		},
		{
			name:           "type mismatch unsupported kind",
			body:           `{"name":"John","age":30,"note":123,"id":"` + validUUID + `"}`,
			expectOK:       false,
			expectHTTPCode: 400,
			expectField:    "note",
			expectError:    "must be value",
		},
		// --- malformed JSON fallback ---
		{
			name:           "malformed JSON",
			body:           `{"name":"John","age":30,"id":"` + validUUID + `"`, // missing brace
			expectOK:       false,
			expectHTTPCode: 400,
			expectError:    "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, ok, err := performRequest(tt.body)

			assert.Equal(t, tt.expectOK, ok)
			if tt.expectOK {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectHTTPCode, w.Code)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, tt.expectHTTPCode, w.Code)

			var resp map[string]any
			_ = json.Unmarshal(w.Body.Bytes(), &resp)

			if tt.expectField != "" {
				assert.Equal(t, tt.expectField, resp["field"])
			}
			if tt.expectError != "" {
				assert.Contains(t, resp["error"], tt.expectError)
			}
		})
	}
}
