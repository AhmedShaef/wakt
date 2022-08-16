package team

import (
	"context"
	"errors"
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

func TestTeam(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testprojectuser")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with Team records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Team.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			npu := NewTeam{
				UID: "5cf37266-3473-4006-984f-9325122678b7",
				PID: "45cf87a3-5915-4079-a9af-6c559239ddbf",
				WID: "7da3ca14-6366-47cf-b953-f706226567d8",
			}

			projectUser, err := core.Create(ctx, npu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create Team.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, projectUser[0].ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Team by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Team by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(projectUser, saved); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same Team. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same Team.", dbtest.Success, testID)

			upd := UpdateTeam{
				Manager: dbtest.BoolPointer(true),
			}

			if err := core.Update(ctx, projectUser[0].ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update Team.", dbtest.Success, testID)

			if saved.Manager == *upd.Manager {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Manager)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Manager)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, projectUser[0].ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete Team.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, projectUser[0].ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve Team.", dbtest.Success, testID)
		}
	}
}

func TestPagingTeam(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through Team records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 Teams.", testID)
		{
			ctx := context.Background()

			Teams1, err := core.QueryWorkspaceTeams(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Teams for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Teams for page 1.", dbtest.Success, testID)

			if len(Teams1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single Team.", dbtest.Success, testID)

			Teams2, err := core.QueryWorkspaceTeams(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Teams for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Teams for page 2.", dbtest.Success, testID)

			if len(Teams2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single Team : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single Team.", dbtest.Success, testID)

			if Teams1[0].ID == Teams2[0].ID {
				t.Logf("\t\tTest %d:\tTeam1: %v", testID, Teams1[0].ID)
				t.Logf("\t\tTest %d:\tTeam2: %v", testID, Teams2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different Teams : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different Teams.", dbtest.Success, testID)

		}
	}
}
