package db

import (
	"time"
)

// Workspace represent the structure we need for moving data
// between the app and the database.
type Workspace struct {
	ID                         string    `db:"workspace_id"`
	Name                       string    `db:"name"`
	Uid                        string    `db:"uid"`
	DefaultHourlyRate          float32   `db:"default_hourly_rate"`
	DefaultCurrency            string    `db:"default_currency"`
	OnlyAdminMayCreateProjects bool      `db:"only_admin_may_create_projects"`
	OnlyAdminSeeBillableRates  bool      `db:"only_admin_see_billable_rates"`
	OnlyAdminSeeTeamDashboard  bool      `db:"only_admin_see_team_dashboard"`
	Rounding                   int       `db:"rounding"`
	RoundingMinutes            int       `db:"rounding_minutes"`
	DateCreated                time.Time `db:"date_created"`
	DateUpdated                time.Time `db:"date_updated"`
	LogoURL                    string    `db:"logo_url"`
}
