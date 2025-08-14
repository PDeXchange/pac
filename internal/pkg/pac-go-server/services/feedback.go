package services

import (
	"fmt"
	"net/http"

	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/client"
	log "github.com/PDeXchange/pac/internal/pkg/pac-go-server/logger"
	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/models"
	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetFeedback		godoc
// @Summary			Get feedbacks submitted by users
// @Description		Get feedback resource
// @Tags			feedbacks
// @Accept			json
// @Produce			json
// @Param			Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success			200
// @Router			/api/v1/feedbacks [get]
func GetFeedback(c *gin.Context) {
	logger := log.GetLogger()
	config := client.GetConfigFromContext(c.Request.Context())
	kc := client.NewKeyCloakClient(config, c.Request.Context())

	pageInt, perPageInt := utils.GetCurrentPageAndPageCount(c)
	startIndex := (pageInt - 1) * perPageInt
	var userId string
	if !kc.IsRole(utils.ManagerRole) {
		userId = kc.GetUserID()
	}
	feedbackFiler := models.FeedbacksFilter{
		UserID: userId,
	}
	feedback, totalCount, err := dbCon.GetFeedbacks(feedbackFiler, startIndex, perPageInt)
	if err != nil {
		logger.Error("failed to get feedbacks from db", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", feedback)})
	}
	// Calculate the total number of pages based on the perPage value
	totalPages := utils.GetTotalPages(totalCount, perPageInt)
	var response = models.FeedbackResponse{
		TotalPages: totalPages,
		TotalItems: totalCount,
		Feedbacks:  feedback,
		Links: models.Links{
			Self: c.Request.URL.String(),
			Next: getNextPageLink(c, pageInt, totalPages),
			Last: getLastPageLink(c, totalPages),
		},
	}
	c.JSON(http.StatusOK, response)
}
