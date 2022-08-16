package team

import (
	"time"
	"unsafe"

	"github.com/AhmedShaef/wakt/business/core/team/db"
)

// Team represents an individual Team.
type Team struct {
	ID          string    `json:"team_id"`
	PID         string    `json:"pid"`
	UID         string    `json:"uid"`
	WID         string    `json:"wid"`
	Manager     bool      `json:"manager"`
	Rate        float64   `json:"rate"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewTeam contains information needed to create a new Team.
type NewTeam struct {
	PID     string  `json:"pid" validate:"required"`
	UID     string  `json:"uid" validate:"required"`
	WID     string  `json:"wid"`
	Manager bool    `json:"manager"`
	Rate    float64 `json:"rate"`
	Puis    string  `json:"puis"`
}

// UpdateTeam defines what information may be provided to modify an existing
// project user. All fields are optional so project users can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateTeam struct {
	Rate    *float64 `json:"rate"`
	Manager *bool    `json:"manager"`
}

// =============================================================================

func toTeam(dbTeam db.Team) Team {
	pu := (*Team)(unsafe.Pointer(&dbTeam))
	return *pu
}

func toTeamSlice(dbTeam []db.Team) []Team {
	teams := make([]Team, len(dbTeam))
	for i, dbTeam := range dbTeam {
		teams[i] = toTeam(dbTeam)
	}
	return teams
}
