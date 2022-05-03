package time_entry

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

func TestClient(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testtimeentry")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with time entry records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single time entry.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			ntc := NewTimeEntry{
				Wid:         "7da3ca14-6366-47cf-b953-f706226567d8",
				Pid:         "45cf87a3-5915-4079-a9af-6c559239ddbf",
				Tid:         "346efd40-6d6e-46d5-b60b-5db9fc171779",
				Start:       time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
				Duration:    60 * 60 * 8,
				CreatedWith: "API",
			}

			timeEntryCreated, err := core.Create(ctx, ntc, "5cf37266-3473-4006-984f-9325122678b7", now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create timeEntryCreated : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create timeEntryCreated.", dbtest.Success, testID)

			savedcreated, err := core.QueryByID(ctx, timeEntryCreated.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve timeEntryCreated by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve timeEntryCreated by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(timeEntryCreated, savedcreated); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same timeEntryCreated. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same timeEntryCreated.", dbtest.Success, testID)

			nts := StartTimeEntry{
				Wid:         "7da3ca14-6366-47cf-b953-f706226567d8",
				Pid:         "45cf87a3-5915-4079-a9af-6c559239ddbf",
				Tid:         "346efd40-6d6e-46d5-b60b-5db9fc171779",
				CreatedWith: "API",
			}
			timeEntryStart, err := core.Start(ctx, nts, "5cf37266-3473-4006-984f-9325122678b7", now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create timeEntryStart : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create timeEntryStart.", dbtest.Success, testID)

			savedStart, err := core.QueryByID(ctx, timeEntryStart.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve timeEntryStart by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve timeEntryStart by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(timeEntryStart, savedStart); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same timeEntryStart. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same timeEntryStart.", dbtest.Success, testID)

			timeEntryStop, err := core.Stop(ctx, timeEntryStart.ID, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to stop timeEntryStop : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to stop timeEntryStop.", dbtest.Success, testID)
			savedStop, err := core.QueryByID(ctx, timeEntryStop.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve timeEntryStop by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve timeEntryStop by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(timeEntryStop, savedStop); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same timeEntryStop. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same timeEntryStop.", dbtest.Success, testID)

			upd := UpdateTimeEntry{
				Description: dbtest.StringPointer("just Updated"),
				Billable:    dbtest.BoolPointer(true),
				Start:       dbtest.TimePointer(time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)),
				//Tags:        []string{"tags"},
			}

			if err := core.Update(ctx, timeEntryCreated.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update timeEntry.", dbtest.Success, testID)

			if savedcreated.Start == *upd.Start {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, savedcreated.Start)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Start)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, savedcreated.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete timeEntry.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, savedcreated.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve timeEntry.", dbtest.Success, testID)

			//ut := UpdateTimeEntryTags{
			//	Tags:    []string{"tags"},
			//	TagMode: "remove",
			//}
			//if err := core.UpdateTags(ctx, savedStart.ID, ut, now); err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to update timeEntry tags : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to update timeEntry tags.", dbtest.Success, testID)

			//tu, err := core.QueryByID(ctx, savedStart.ID)
			//if err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve timeEntry : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve timeEntry.", dbtest.Success, testID)
			//
			//if diff := cmp.Diff(tu.Tags, upd.Tags); diff == "" {
			//	t.Fatalf("\t%s\tTest %d:\tShould get back the same timeEntryStop. Diff:\n%s", dbtest.Failed, testID, diff)
			//}
			//t.Logf("\t%s\tTest %d:\tShould get back the same timeEntryStop.", dbtest.Success, testID)

			_, err = core.QueryDash(ctx, "5cf37266-3473-4006-984f-9325122678b7")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve timeEntry.", dbtest.Success, testID)
			//TODO: check dash
		}
	}
}

func TestPagingClient(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testpaging")
	t.Cleanup(teardown)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbschema.Seed(ctx, db)

	core := NewCore(log, db)

	t.Log("Given the need to page through timeEntry records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 timeEntries.", testID)
		{
			ctx := context.Background()

			timeEntryRunning1, err := core.QueryRunning(ctx, "5cf37266-3473-4006-984f-9325122678b7", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve running timeEntries for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve running timeEntries for page 1.", dbtest.Success, testID)

			if len(timeEntryRunning1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single running timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single running timeEntry.", dbtest.Success, testID)

			timeEntryRunning2, err := core.QueryRunning(ctx, "5cf37266-3473-4006-984f-9325122678b7", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve running timeEntries for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve running timeEntries for page 2.", dbtest.Success, testID)

			if len(timeEntryRunning2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single running timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single running timeEntry.", dbtest.Success, testID)

			if timeEntryRunning1[0].ID == timeEntryRunning2[0].ID {
				t.Logf("\t\tTest %d:\ttimeEntry1: %v", testID, timeEntryRunning1[0].ID)
				t.Logf("\t\tTest %d:\ttimeEntry2: %v", testID, timeEntryRunning2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different running timeEntries : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different timeEntries.", dbtest.Success, testID)

			//==================================================================================================================

			timeEntryRange1, err := core.QueryRange(ctx, "5cf37266-3473-4006-984f-9325122678b7", 1, 1, time.Date(2006, time.October, 1, 0, 0, 0, 0, time.UTC), time.Now())
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace timeEntryRange for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace timeEntryRange for page 1.", dbtest.Success, testID)

			if len(timeEntryRange1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single timeEntry.", dbtest.Success, testID)

			timeEntryRange2, err := core.QueryRange(ctx, "5cf37266-3473-4006-984f-9325122678b7", 2, 1, time.Date(2006, time.October, 1, 0, 0, 0, 0, time.UTC), time.Now())
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace timeEntryRange for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace timeEntryRange for page 2.", dbtest.Success, testID)

			if len(timeEntryRange2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single timeEntry : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single timeEntry.", dbtest.Success, testID)

			if timeEntryRange1[0].ID == timeEntryRange2[0].ID {
				t.Logf("\t\tTest %d:\ttimeEntry1: %v", testID, timeEntryRange1[0].ID)
				t.Logf("\t\tTest %d:\ttimeEntry2: %v", testID, timeEntryRange2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different timeEntryRange : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different timeEntryRange.", dbtest.Success, testID)

		}
	}
}
