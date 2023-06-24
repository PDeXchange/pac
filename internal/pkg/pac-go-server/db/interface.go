package db

import (
	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/models"
)

type DB interface {
	Connect() error
	Disconnect() error

	GetRequestsByUserID(id, requestType string) ([]models.Request, error)
	NewRequest(request *models.Request) error
	GetRequestByGroupIDAndUserID(groupID, userID string) ([]models.Request, error)
	GetRequestByID(string) (*models.Request, error)
	DeleteRequest(string) error
	UpdateRequestState(id string, state models.RequestStateType) error
	GetRequestByServiceName(string) ([]models.Request, error)

	GetKeyByID(id string) (*models.Key, error)
	GetKeyByUserID(userid string) ([]models.Key, error)
	CreateKey(key *models.Key) error
	DeleteKey(string) error

	// Implementations for group quota.
	NewQuota(*models.Quota) error
	UpdateQuota(*models.Quota) error
	DeleteQuota(string) error
	GetQuotaForGroupID(string) (*models.Quota, error)
	GetGroupsQuota([]string) ([]models.Quota, error)
}
