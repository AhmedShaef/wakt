package db

import (
	"time"
)

// Workspace represent the structure we need for moving data
// between the app and the database.
type Workspace struct {
	ID                         string    `db:"workspace_id"`
	Name                       string    `db:"name"`
	Profile                    int       `db:"profile"`
	Premium                    bool      `db:"premium"`
	Admin                      bool      `db:"admin"`
	DefaultHourlyRate          float32   `db:"default_hourly_rate"`
	DefaultCurrency            string    `db:"default_currency"`
	OnlyAdminMayCreateProjects bool      `db:"only_admin_may_create_projects"`
	OnlyAdminSeeBillableRates  bool      `db:"only_admin_see_billable_rates"`
	OnlyAdminSeeTeamDashboard  bool      `db:"only_admin_see_team_dashboard"`
	ProjectBillableByDefault   bool      `db:"project_billable_by_default"`
	Rounding                   int       `db:"rounding"`
	RoundingMinutes            int       `db:"rounding_minutes"`
	DateUpdated                time.Time `db:"date_updated"` //TODO date created
	LogoURL                    string    `db:"logo_url"`
	IcalURL                    string    `db:"ical_url"`
	IcalEnabled                bool      `db:"ical_enabled"`
}
