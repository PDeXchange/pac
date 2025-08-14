package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetFeedback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, mockDBClient, mockKCClient, tearDown := setUp(t)
	defer tearDown()

	testcases := []struct {
		name           string
		mockFunc       func()
		requestContext testContext
		httpStatus     int
	}{
		{
			name: "Get feedback successfully",
			mockFunc: func() {
				mockKCClient.EXPECT().IsRole(gomock.Any()).Return(true).Times(1)
				mockDBClient.EXPECT().GetFeedbacks(gomock.Any(), gomock.Any(), gomock.Any()).Return(getResource("get-feedbacks", nil).([]models.Feedback), int64(1), nil).Times(1)
			},
			requestContext: formContext(customValues{
				"userid":                "12345",
				"keycloak_hostname":     "127.0.0.1",
				"keycloak_access_token": "Bearer test-token",
				"keycloak_realm":        "test-pac",
			}),
			httpStatus: http.StatusOK,
		},
		{
			name: "Get feedback with InternalServerError",
			mockFunc: func() {
				mockKCClient.EXPECT().IsRole(gomock.Any()).Return(true).Times(1)
				mockDBClient.EXPECT().GetFeedbacks(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(1), (errors.New(""))).Times(1)
			},
			requestContext: formContext(customValues{
				"userid":                "12345",
				"keycloak_hostname":     "127.0.0.1",
				"keycloak_access_token": "Bearer test-token",
				"keycloak_realm":        "test-pac",
			}),
			httpStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, err := http.NewRequest(http.MethodGet, "/feedbacks", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := getContext(tc.requestContext)
			c.Request = req.WithContext(ctx)
			dbCon = mockDBClient
			GetFeedback(c)
			assert.Equal(t, tc.httpStatus, c.Writer.Status())
		})
	}
}
