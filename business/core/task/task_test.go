package task

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

func TestTask(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testtask")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with task records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single task.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nt := NewTask{
				Name: "Ahmed Shaef",
				PID:  "45cf87a3-5915-4079-a9af-6c559239ddbf",
				Wid:  "7da3ca14-6366-47cf-b953-f706226567d8",
			}

			tsk, err := core.Create(ctx, "5cf37266-3473-4006-984f-9325122678b7", nt, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create task.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, tsk.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve task by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve task by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(tsk, saved); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same task. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same task.", dbtest.Success, testID)

			upd := UpdateTask{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, tsk.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update task.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, tsk.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete task.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, tsk.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve task.", dbtest.Success, testID)
		}
	}
}

func TestPagingTask(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through task records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 tasks.", testID)
		{
			ctx := context.Background()

			tasks1, err := core.QueryProjectTasks(ctx, "45cf87a3-5915-4079-a9af-6c559239ddbf", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 1.", dbtest.Success, testID)

			if len(tasks1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single task.", dbtest.Success, testID)

			tasks2, err := core.QueryProjectTasks(ctx, "45cf87a3-5915-4079-a9af-6c559239ddbf", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 2.", dbtest.Success, testID)

			if len(tasks2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single task.", dbtest.Success, testID)

			if tasks1[0].ID == tasks2[0].ID {
				t.Logf("\t\tTest %d:\ttask1: %v", testID, tasks1[0].ID)
				t.Logf("\t\tTest %d:\ttask2: %v", testID, tasks2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different tasks : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different tasks.", dbtest.Success, testID)

			//========================================================================================================================

			tasks3, err := core.QueryWorkspaceTasks(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 1.", dbtest.Success, testID)

			if len(tasks3) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single task.", dbtest.Success, testID)

			tasks4, err := core.QueryWorkspaceTasks(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve tasks for page 2.", dbtest.Success, testID)

			if len(tasks4) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single task : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single task.", dbtest.Success, testID)

			if tasks3[0].ID == tasks4[0].ID {
				t.Logf("\t\tTest %d:\ttask1: %v", testID, tasks3[0].ID)
				t.Logf("\t\tTest %d:\ttask2: %v", testID, tasks4[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different tasks : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different tasks.", dbtest.Success, testID)
		}
	}
}
