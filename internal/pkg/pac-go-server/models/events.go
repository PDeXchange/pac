package models

import (
	"bytes"
	"context"
	"html/template"
	"time"

	"github.com/PDeXchange/pac/internal/pkg/pac-go-server/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string
type EventLogLevel string

const (
	EventGroupJoinRequest     EventType = "GROUP_JOIN_REQUEST"
	EventServiceExpiryRequest EventType = "SERVICE_EXPIRY_REQUEST"
	EventGroupExitRequest     EventType = "GROUP_EXIT_REQUEST"

	EventTypeRequestApproved EventType = "REQUEST_APPROVED"
	EventTypeRequestRejected EventType = "REQUEST_REJECTED"
	EventTypeRequestDeleted  EventType = "REQUEST_DELETED"

	EventCatalogCreate EventType = "CATALOG_CREATE"
	EventCatalogUpdate EventType = "CATALOG_UPDATE"
	EventCatalogDelete EventType = "CATALOG_DELETE"
	EventCatalogRetire EventType = "CATALOG_RETIRE"

	EventServiceCreate EventType = "SERVICE_CREATE"
	EventServiceUpdate EventType = "SERVICE_UPDATE"
	EventServiceDelete EventType = "SERVICE_DELETE"

	EventLogLevelINFO  EventLogLevel = "INFO"
	EventLogLevelERROR EventLogLevel = "ERROR"
)

type Event struct {
	// ID is the event identifier
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// EventType is the event type
	Type EventType `json:"type" bson:"type"`
	// CreatedAt is the time the event was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// Originator is the user that created the event
	Originator string `json:"originator" bson:"originator"`
	// UserID is the user that the event is about
	UserID string `json:"user_id" bson:"user_id"`
	// UserEmail is the user email that the event is about
	UserEmail string `json:"user_email" bson:"user_email"`
	// Notify is a flag to indicate if the user should be notified
	Notify bool `json:"notify" bson:"notify"`
	// NotifyAdmin is a flag to indicate if the admin should be notified
	NotifyAdmin bool `json:"notify_admin" bson:"notify_admin"`
	// Notified is a flag to indicate if the user has been notified
	Notified bool `json:"notified" bson:"notified"`
	// Log contains all the event information
	Log EventLog `json:"log" bson:"log"`
}

type EventLog struct {
	Level   EventLogLevel `json:"level" bson:"level"`
	Message string        `json:"message" bson:"message"`
}

type EventResponse struct {
	// TotalPages is the total number of pages
	TotalPages int64 `json:"total_pages"`
	// TotalItems is the total number of items
	TotalItems int64 `json:"total_items"`
	// Events is the list of events
	Events []Event `json:"events"`
	// Links contains the links for the current page, next page and last page
	Links Links `json:"links"`
}

func NewEvent(ctx context.Context, userid, originator string, typ EventType) (*Event, error) {
	email := ""
	// if originator is not the same as userid, then the event is created by an admin, hence get the user info for the supplied userid
	if originator != userid {
		user, err := client.NewKeyClockClient(ctx).GetUser(userid)
		if err != nil {
			return nil, err
		}
		email = *user.Email
	} else {
		// Get user info for the token in the context for the non-admin user
		usr, err := client.NewKeyClockClient(ctx).GetUserInfo()
		if err != nil {
			return nil, err
		}
		email = *usr.Email
	}

	return &Event{
		Type:      typ,
		CreatedAt: time.Now(),

		UserID:     userid,
		UserEmail:  email,
		Originator: originator,
	}, nil
}

// SetType sets the event type
func (e *Event) SetType(t EventType) {
	e.Type = t
}

// SetUserID sets the user id
func (e *Event) SetUserID(id string) {
	e.UserID = id
}

// SetNotify sets the notify flag
func (e *Event) SetNotify() {
	e.Notify = true
}

// SetNotifyAdmin sets the notify admin flag
func (e *Event) SetNotifyAdmin() {
	e.NotifyAdmin = true
}

// SetNotifiyBoth sets both Notify and NotifyAdmin to true
func (e *Event) SetNotifiyBoth() {
	e.Notify = true
	e.NotifyAdmin = true
}

// SetNotified sets the notified flag
func (e *Event) SetNotified(b bool) {
	e.Notified = b
}

func (e *Event) SetLog(level EventLogLevel, message string) {
	e.Log = EventLog{
		Level:   EventLogLevelINFO,
		Message: message,
	}
}

var bodyTemplate = `
{{- if .Log.Message -}}
Hi,

{{ .Log.Message }}

Please visit the PAC portal for more details.

Note: This is an auto-generated email. Please do not reply to this email.
Generated at: {{ .CreatedAt.Format "Jan 02, 2006 15:04:05 UTC" }}

Thanks,
PAC Support.
{{- end -}}
`

func (e *Event) ComposeMailBody() (string, error) {
	tmpl, err := template.New("pac").Parse(bodyTemplate)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	tmpl.Execute(&tpl, e)
	return tpl.String(), nil
}
