package project

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

func TestProject(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testproject")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with project records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single project.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nc := NewProject{
				Name: "Ahmed Shaef",
				Wid:  "7da3ca14-6366-47cf-b953-f706226567d8",
				Cid:  "c78db68e-e004-44f5-895b-ba562dc53d9d",
			}

			project, err := core.Create(ctx, nc, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create project.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, project.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve project by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve project by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(project, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same Project. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same project.", dbtest.Success, testID)

			upd := UpdateProject{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, project.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update project.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, project.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete project.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, project.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve project.", dbtest.Success, testID)
		}
	}
}

func TestPagingProject(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through project records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 projects.", testID)
		{
			ctx := context.Background()

			projects1, err := core.QueryWorkspaceProjects(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve projects for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve projects for page 1.", dbtest.Success, testID)

			if len(projects1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single project.", dbtest.Success, testID)

			projects2, err := core.QueryWorkspaceProjects(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve projects for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve projects for page 2.", dbtest.Success, testID)

			if len(projects2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single project : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single project.", dbtest.Success, testID)

			if projects1[0].ID == projects2[0].ID {
				t.Logf("\t\tTest %d:\tproject1: %v", testID, projects1[0].ID)
				t.Logf("\t\tTest %d:\tproject2: %v", testID, projects2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different projects : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different projects.", dbtest.Success, testID)
		}
	}
}