package tag

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

func TestTag(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testtag")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with tag records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single tag.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nt := NewTag{
				Name: "Ahmed Shaef",
				WID:  "7da3ca14-6366-47cf-b953-f706226567d8",
			}

			tag, err := core.Create(ctx, nt, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create tag.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, tag.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tag by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tag by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(tag, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same tag. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same tag.", dbtest.Success, testID)

			upd := UpdateTag{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, tag.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update tag.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, tag.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete tag.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, tag.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve tag.", dbtest.Success, testID)
		}
	}
}

func TestPagingTag(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through tag records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 tags.", testID)
		{
			ctx := context.Background()

			tags1, err := core.QueryWorkspaceTags(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tags for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tags for page 1.", dbtest.Success, testID)

			if len(tags1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single tag.", dbtest.Success, testID)

			tags2, err := core.QueryWorkspaceTags(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tags for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tags for page 2.", dbtest.Success, testID)

			if len(tags2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single tag : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single tag.", dbtest.Success, testID)

			if tags1[0].ID == tags2[0].ID {
				t.Logf("\t\tTest %d:\ttag1: %v", testID, tags1[0].ID)
				t.Logf("\t\tTest %d:\ttag2: %v", testID, tags2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different tags : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different tags.", dbtest.Success, testID)
		}
	}
}
