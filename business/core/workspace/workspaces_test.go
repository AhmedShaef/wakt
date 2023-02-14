package workspace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AhmedShaef/wakt/business/data/dbschema"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/foundation/docker"
	"github.com/google/go-cmp/cmp"
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

func TestWorkspace(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testworkspace")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with workspace records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single workspace.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nt := NewWorkspace{
				Name: "Ahmed Shaef",
				UID:  "5cf37266-3473-4006-984f-9325122678b7",
			}

			workspace, err := core.Create(ctx, nt, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create workspace : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create workspace.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, workspace.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(workspace, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same workspace. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same workspace.", dbtest.Success, testID)

			upd := UpdateWorkspace{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, workspace.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update workspace : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update workspace.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}
		}
	}
}

func TestPagingWorkspace(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through workspace records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 workspaces.", testID)
		{
			ctx := context.Background()

			workspaces1, err := core.List(ctx, "5cf37266-3473-4006-984f-9325122678b7", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspaces for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspaces for page 1.", dbtest.Success, testID)

			if len(workspaces1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single workspace : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single workspace.", dbtest.Success, testID)

			workspaces2, err := core.List(ctx, "5cf37266-3473-4006-984f-9325122678b7", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspaces for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspaces for page 2.", dbtest.Success, testID)

			if len(workspaces2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single workspace : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single workspace.", dbtest.Success, testID)

			if workspaces1[0].ID == workspaces2[0].ID {
				t.Logf("\t\tTest %d:\tworkspace1: %v", testID, workspaces1[0].ID)
				t.Logf("\t\tTest %d:\tworkspace2: %v", testID, workspaces2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different workspaces : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different workspaces.", dbtest.Success, testID)

		}
	}
}
