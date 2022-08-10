package workspaceuser

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/data/dbschema"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/foundation/docker"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestWorkspaceUser(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testworkspaceuser")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with workspace_user records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single workspace_user.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			workspaceUser, err := core.Create(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", "5cf37266-3473-4006-984f-9325122678b7", now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create workspace_user.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, workspaceUser.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace_user by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace_user by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(workspaceUser, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same workspace_user. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same workspace_user.", dbtest.Success, testID)

			iu := InviteUsers{
				Emails:    []string{"hasAccount@example.com", "hasntaccount@example.com"},
				InviterID: "32c1494f-1c1f-4981-857f-b0526cb654ec",
			}

			workspaceUsers, err := core.InviteUser(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", iu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create workspace_user.", dbtest.Success, testID)

			saved2, err := core.QueryByID(ctx, workspaceUsers[0].ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace_user by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace_user by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(workspaceUsers[0], saved2); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same workspace_user. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same workspace_user.", dbtest.Success, testID)

			upd := UpdateWorkspaceUser{
				Active: dbtest.BoolPointer(true),
				Admin:  dbtest.BoolPointer(true),
			}

			if err := core.Update(ctx, workspaceUsers[0].ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update workspace_user.", dbtest.Success, testID)

			if saved2.Active == *upd.Active {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved2.Active)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Active)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, workspaceUsers[0].ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete workspace_user.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, workspaceUsers[0].ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve workspace_user.", dbtest.Success, testID)

			_, err = core.QueryByuIDwID(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", "5cf37266-3473-4006-984f-9325122678b7")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve workspace_user.", dbtest.Success, testID)
		}
	}
}

func TestPagingWorkspaceUser(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through workspace_user records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 workspace_users.", testID)
		{
			ctx := context.Background()

			workspaceUsers1, err := core.Query(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace_users for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace_users for page 1.", dbtest.Success, testID)

			if len(workspaceUsers1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single workspace_user.", dbtest.Success, testID)

			workspaceUsers2, err := core.Query(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace_users for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace_users for page 2.", dbtest.Success, testID)

			if len(workspaceUsers2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single workspace_user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single workspace_user.", dbtest.Success, testID)

			if workspaceUsers1[0].ID == workspaceUsers2[0].ID {
				t.Logf("\t\tTest %d:\tworkspace_user1: %v", testID, workspaceUsers1[0].ID)
				t.Logf("\t\tTest %d:\tworkspace_user2: %v", testID, workspaceUsers2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different workspace_users : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different workspace_users.", dbtest.Success, testID)
		}
	}
}
