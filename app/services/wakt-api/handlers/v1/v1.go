// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/clientgrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/groupgrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/projectgrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/projectusergrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/taggrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/taskgrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/timeentrygrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/workspacegrp"
	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/workspaceusergrp"
	"github.com/AhmedShaef/wakt/business/core/client"
	"github.com/AhmedShaef/wakt/business/core/group"
	"github.com/AhmedShaef/wakt/business/core/project"
	"github.com/AhmedShaef/wakt/business/core/project_user"
	"github.com/AhmedShaef/wakt/business/core/tag"
	"github.com/AhmedShaef/wakt/business/core/task"
	"github.com/AhmedShaef/wakt/business/core/time_entry"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	"net/http"

	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers/v1/usergrp"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	"github.com/AhmedShaef/wakt/business/web/v1/mid"
	"github.com/AhmedShaef/wakt/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)

	// Register client management endpoints.
	cgh := clientgrp.Handlers{
		Client:    client.NewCore(cfg.Log, cfg.DB),
		Workspace: workspace.NewCore(cfg.Log, cfg.DB),
		Project:   project.NewCore(cfg.Log, cfg.DB),
		User:      user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/client", cgh.Create, authen)
	app.Handle(http.MethodPut, version, "/client/:id", cgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/client/:id", cgh.Delete, authen)
	app.Handle(http.MethodGet, version, "/client/:id", cgh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/client/:page/:rows", cgh.Query, authen)
	app.Handle(http.MethodGet, version, "/client/:id/project/:page/:rows", cgh.QueryClientProjects, authen)

	// Register group management endpoints.
	ggh := groupgrp.Handlers{
		Group:     group.NewCore(cfg.Log, cfg.DB),
		Workspace: workspace.NewCore(cfg.Log, cfg.DB),
		User:      user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/group", ggh.Create, authen)
	app.Handle(http.MethodPut, version, "/group/:id", ggh.Update, authen)
	app.Handle(http.MethodDelete, version, "/group/:id", ggh.Delete, authen)

	// Register project management endpoints.
	pgh := projectgrp.Handlers{
		Project:       project.NewCore(cfg.Log, cfg.DB),
		Workspace:     workspace.NewCore(cfg.Log, cfg.DB),
		User:          user.NewCore(cfg.Log, cfg.DB),
		Task:          task.NewCore(cfg.Log, cfg.DB),
		WorkspaceUser: workspace_user.NewCore(cfg.Log, cfg.DB),
		ProjectUser:   project_user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/project", pgh.Create, authen)
	app.Handle(http.MethodPut, version, "/project/:id", pgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/project/:id", pgh.BulkDelete, authen)
	app.Handle(http.MethodGet, version, "/project/:id", pgh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/project/:id/task/:page/:rows", pgh.QueryProjectTasks, authen)

	// Register project user management endpoints.
	pugh := projectusergrp.Handlers{
		ProjectUser: project_user.NewCore(cfg.Log, cfg.DB),
		Workspace:   workspace.NewCore(cfg.Log, cfg.DB),
		User:        user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/team", pugh.Add, authen)
	app.Handle(http.MethodPut, version, "/team/:id", pugh.BulkUpdate, authen)
	app.Handle(http.MethodDelete, version, "/team/:id", pugh.BulkDelete, authen)

	// Register tag management endpoints.
	tgh := taggrp.Handlers{
		Tag:       tag.NewCore(cfg.Log, cfg.DB),
		Workspace: workspace.NewCore(cfg.Log, cfg.DB),
		User:      user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/tag", tgh.Create, authen)
	app.Handle(http.MethodPut, version, "/tag/:id", tgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/tag/:id", tgh.Delete, authen)

	// Register task management endpoints.
	tkgh := taskgrp.Handlers{
		Task:      task.NewCore(cfg.Log, cfg.DB),
		Workspace: workspace.NewCore(cfg.Log, cfg.DB),
		User:      user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/task", tkgh.Create, authen)
	app.Handle(http.MethodPut, version, "/task/:id", tkgh.BulkUpdate, authen)
	app.Handle(http.MethodDelete, version, "/task/:id", tkgh.BulkDelete, authen)
	app.Handle(http.MethodGet, version, "/task/:id", tkgh.QueryByID, authen)

	// Register time entry management endpoints.
	tegh := timeentrygrp.Handlers{
		TimeEntry: time_entry.NewCore(cfg.Log, cfg.DB),
		Workspace: workspace.NewCore(cfg.Log, cfg.DB),
		User:      user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/time_entry", tegh.Create, authen)
	app.Handle(http.MethodPost, version, "/time_entry/start", tegh.Start, authen)
	app.Handle(http.MethodPut, version, "/time_entry/:id/stop", tegh.Stop, authen)
	app.Handle(http.MethodGet, version, "/time_entry/:id", tegh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/time_entry/running/:page/:rows", tegh.QueryRunning, authen)
	app.Handle(http.MethodPut, version, "/time_entry/update/:id", tegh.Update, authen)
	app.Handle(http.MethodPut, version, "/time_entry/tags/:id", tegh.UpdateTags, authen)
	app.Handle(http.MethodDelete, version, "/time_entry/delete/:id", tegh.Delete, authen)
	app.Handle(http.MethodGet, version, "/time_entry/:page/:rows", tegh.QueryRange, authen)
	app.Handle(http.MethodGet, version, "/dashboard", tegh.QueryDash, authen)

	// Register user management and authentication endpoints.
	ugh := usergrp.Handlers{
		User:          user.NewCore(cfg.Log, cfg.DB),
		Auth:          cfg.Auth,
		Workspace:     workspace.NewCore(cfg.Log, cfg.DB),
		WorkspaceUser: workspace_user.NewCore(cfg.Log, cfg.DB),
		Project:       project.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodPost, version, "/signup", ugh.SignUp)
	app.Handle(http.MethodGet, version, "/token", ugh.NewToken)
	app.Handle(http.MethodPost, version, "/image", ugh.UpdateImage, authen)
	app.Handle(http.MethodGet, version, "/me", ugh.QueryByID, authen)
	app.Handle(http.MethodPut, version, "/me", ugh.Update, authen)
	app.Handle(http.MethodPut, version, "/change_password", ugh.ChangePassword, authen)
	app.Handle(http.MethodPut, version, "/change_image", ugh.UpdateImage, authen)
	app.Handle(http.MethodGet, version, "/user_projects/:page/:rows", ugh.QueryUserProjects, authen)

	// Register workspace management endpoints.
	wgh := workspacegrp.Handlers{
		Workspace:     workspace.NewCore(cfg.Log, cfg.DB),
		User:          user.NewCore(cfg.Log, cfg.DB),
		Client:        client.NewCore(cfg.Log, cfg.DB),
		Group:         group.NewCore(cfg.Log, cfg.DB),
		Project:       project.NewCore(cfg.Log, cfg.DB),
		Task:          task.NewCore(cfg.Log, cfg.DB),
		Tag:           tag.NewCore(cfg.Log, cfg.DB),
		ProjectUser:   project_user.NewCore(cfg.Log, cfg.DB),
		WorkspaceUser: workspace_user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/workspace", wgh.Create, authen)
	app.Handle(http.MethodPut, version, "/workspace/:id", wgh.Update, authen)
	app.Handle(http.MethodPut, version, "/workspace/logo/:id", wgh.UpdateLogo, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id", wgh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/workspace/:page/:rows", wgh.Query, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/users/:page/:rows", wgh.QueryWorkspaceUsers, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/clients/:page/:rows", wgh.QueryWorkspaceClients, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/groups/:page/:rows", wgh.QueryWorkspaceGroups, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/projects/:page/:rows", wgh.QueryWorkspaceProjects, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/tasks/:page/:rows", wgh.QueryWorkspaceTasks, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/tags/:page/:rows", wgh.QueryWorkspaceTags, authen)
	app.Handle(http.MethodGet, version, "/workspace/:id/project_users/:page/:rows", wgh.QueryWorkspaceProjectUsers, authen)

	// Register workspace user management endpoints.
	wugh := workspaceusergrp.Handlers{
		WorkspaceUser: workspace_user.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/workspace_user", wugh.Invite, authen)
	app.Handle(http.MethodPut, version, "/workspace_user/:id", wugh.Update, authen)
	app.Handle(http.MethodDelete, version, "/workspace_user/:id", wugh.Delete, authen)

}
