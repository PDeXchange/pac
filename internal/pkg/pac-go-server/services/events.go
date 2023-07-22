package services

import (
	"net/http"
	"strconv"

	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/models"
	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetEvents returns all events
func GetEvents(c *gin.Context) {
	kc := utils.NewKeyClockClient(c.Request.Context())

	page := c.DefaultQuery("page", "1")         // Get the page number from the query parameter
	perPage := c.DefaultQuery("per_page", "10") // Get the number of items per page from the query parameter

	// Convert the page and perPage values to integers
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	perPageInt, _ := strconv.ParseInt(perPage, 10, 64)

	// Calculate the starting index and ending index for the current page
	startIndex := (pageInt - 1) * perPageInt

	// MongoDB options to handle paging
	findOptions := options.Find()
	findOptions.SetSkip(startIndex)
	findOptions.SetLimit(perPageInt)

	var userID string
	if !kc.IsRole(utils.ManagerRole) {
		// Get authenticated user's ID
		userID = kc.GetUserID()
	}
	events, totalCount, err := dbCon.GetEventsByUserID(userID, startIndex, perPageInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate the total number of pages based on the perPage value
	totalPages := totalCount / perPageInt
	if totalCount%perPageInt != 0 {
		totalPages++
	}

	var response = struct {
		TotalPages int64          `json:"total_pages"`
		TotalItems int64          `json:"total_items"`
		Events     []models.Event `json:"events"`
		Links      struct {
			Self string `json:"self"`
			Next string `json:"next,omitempty"`
			Last string `json:"last,omitempty"`
		} `json:"links"`
	}{
		TotalPages: totalPages,
		TotalItems: totalCount,
		Events:     events,
		Links: struct {
			Self string `json:"self"`
			Next string `json:"next,omitempty"`
			Last string `json:"last,omitempty"`
		}{
			Self: c.Request.URL.String(),
			// Get the next page link
			Next: getNextPageLink(c, pageInt, totalPages),
			// Get the last page link
			Last: getLastPageLink(c, totalPages),
		},
	}

	c.JSON(http.StatusOK, response)
}

func getNextPageLink(c *gin.Context, currentPage, totalPages int64) string {
	if currentPage >= totalPages {
		return ""
	}
	nextPage := currentPage + 1
	return getPaginationLink(c, nextPage)
}

func getLastPageLink(c *gin.Context, totalPages int64) string {
	if totalPages <= 1 {
		return ""
	}
	return getPaginationLink(c, totalPages)
}

func getPaginationLink(c *gin.Context, page int64) string {
	queryParams := c.Request.URL.Query()
	queryParams.Set("page", strconv.FormatInt(page, 10))
	return c.Request.URL.Path + "?" + queryParams.Encode()
}
