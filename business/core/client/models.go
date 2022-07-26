package client

import (
	"time"
	"unsafe"

	"github.com/AhmedShaef/wakt/business/core/client/db"
)

// Client represents an individual client.
type Client struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	UID         string    `json:"uid"`
	WID         string    `json:"wid"`
	Notes       string    `json:"notes"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewClient contains information needed to create a new client.
type NewClient struct {
	Name  string `json:"name" validate:"required"`
	WID   string `json:"wid"`
	Notes string `json:"notes"`
}

// UpdateClient defines what information may be provided to modify an existing
// client. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateClient struct {
	Name  *string `json:"name"`
	Notes *string `json:"notes"`
}

// =============================================================================

func toClient(dbClient db.Client) Client {
	pu := (*Client)(unsafe.Pointer(&dbClient))
	return *pu
}

func toClientsSlice(dbClient []db.Client) []Client {
	clients := make([]Client, len(dbClient))
	for i, dbclint := range dbClient {
		clients[i] = toClient(dbclint)
	}
	return clients
}
