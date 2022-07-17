package workspace

import (
	"github.com/AhmedShaef/wakt/business/core/workspace/db"
	"time"
	"unsafe"
)

// Workspace represents an individual Group.
type Workspace struct {
	ID                         string    `json:"id"`
	Name                       string    `json:"name"`
	UID                        string    `json:"uid"`
	DefaultHourlyRate          float32   `json:"default_hourly_rate"`
	DefaultCurrency            string    `json:"default_currency"`
	OnlyAdminMayCreateProjects bool      `json:"only_admin_may_create_projects"`
	OnlyAdminSeeBillableRates  bool      `json:"only_admin_see_billable_rates"`
	OnlyAdminSeeTeamDashboard  bool      `json:"only_admin_see_team_dashboard"`
	Rounding                   int       `json:"rounding"`
	RoundingMinutes            int       `json:"rounding_minutes"`
	DateCreated                time.Time `json:"date_created"`
	DateUpdated                time.Time `json:"date_updated"`
	LogoURL                    string    `json:"logo_url"`
}

// NewWorkspace contains information needed to create a new Group.
type NewWorkspace struct {
	Name string `json:"name" validate:"required"`
	UID  string `json:"uid" validate:"required"`
}

// UpdateWorkspace defines what information may be provided to modify an existing
// group. All fields are optional so group can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateWorkspace struct {
	Name                       *string  `json:"name"`
	DefaultHourlyRate          *float32 `json:"default_hourly_rate"`
	DefaultCurrency            *string  `json:"default_currency"`
	OnlyAdminMayCreateProjects *bool    `json:"only_admin_may_create_projects"`
	OnlyAdminSeeBillableRates  *bool    `json:"only_admin_see_billable_rates"`
	OnlyAdminSeeTeamDashboard  *bool    `json:"only_admin_see_team_dashboard"`
	Rounding                   *int     `json:"rounding" validate:"omitempty,eq=0|eq=1|eq=-1"`
	RoundingMinutes            *int     `json:"rounding_minutes"`
	LogoURL                    string   `json:"logo_url"`
}

// =============================================================================

func toWorkspace(dbWorkspace db.Workspace) Workspace {
	pu := (*Workspace)(unsafe.Pointer(&dbWorkspace))
	return *pu
}

func toWorkspaceSlice(dbWorkspace []db.Workspace) []Workspace {
	workspaces := make([]Workspace, len(dbWorkspace))
	for i, dbUsr := range dbWorkspace {
		workspaces[i] = toWorkspace(dbUsr)
	}
	return workspaces
}
