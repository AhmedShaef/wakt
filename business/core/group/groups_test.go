package group

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

func TestGroup(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testgroup")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with group records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single group.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

			ng := NewGroup{
				Name: "Ahmed Shaef",
				WID:  "7da3ca14-6366-47cf-b953-f706226567d8",
			}

			group, err := core.Create(ctx, "5cf37266-3473-4006-984f-9325122678b7", ng, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create group.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, group.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve group by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve group by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(group, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same group. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same group.", dbtest.Success, testID)

			upd := UpdateGroup{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, group.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update group.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, group.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete group.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, group.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve group.", dbtest.Success, testID)
		}
	}
}

func TestPagingGroup(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through group records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 groups.", testID)
		{
			ctx := context.Background()

			groups1, err := core.QueryWorkspaceGroups(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve groups for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve groups for page 1.", dbtest.Success, testID)

			if len(groups1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single group.", dbtest.Success, testID)

			groups2, err := core.QueryWorkspaceGroups(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve groups for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve groups for page 2.", dbtest.Success, testID)

			if len(groups2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single group : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single group.", dbtest.Success, testID)

			if groups1[0].ID == groups2[0].ID {
				t.Logf("\t\tTest %d:\tgroup1: %v", testID, groups1[0].ID)
				t.Logf("\t\tTest %d:\tgroup2: %v", testID, groups2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different groups : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different groups.", dbtest.Success, testID)
		}
	}
}
