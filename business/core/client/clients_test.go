package client

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
	log, db, teardown := dbtest.NewUnit(t, c, "testclient")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with client records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single client.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nc := NewClient{
				Name: "Ahmed Shaef",
				WID:  "7da3ca14-6366-47cf-b953-f706226567d8",
			}

			client, err := core.Create(ctx, nc, "5cf37266-3473-4006-984f-9325122678b7", now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create client.", dbtest.Success, testID)

			saved, err := core.QueryByID(ctx, client.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve client by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve client by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(client, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same client. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same client.", dbtest.Success, testID)

			upd := UpdateClient{
				Name: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, client.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update client.", dbtest.Success, testID)

			if saved.Name == *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			if err := core.Delete(ctx, client.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete client.", dbtest.Success, testID)

			_, err = core.QueryByID(ctx, client.ID)
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve client.", dbtest.Success, testID)
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

	t.Log("Given the need to page through client records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 clients.", testID)
		{
			ctx := context.Background()

			clients1, err := core.Query(ctx, "5cf37266-3473-4006-984f-9325122678b7", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve clients for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve clients for page 1.", dbtest.Success, testID)

			if len(clients1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single client.", dbtest.Success, testID)

			clients2, err := core.Query(ctx, "5cf37266-3473-4006-984f-9325122678b7", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve clients for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve clients for page 2.", dbtest.Success, testID)

			if len(clients2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single client.", dbtest.Success, testID)

			if clients1[0].ID == clients2[0].ID {
				t.Logf("\t\tTest %d:\tclient1: %v", testID, clients1[0].ID)
				t.Logf("\t\tTest %d:\tclient2: %v", testID, clients2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different clients : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different clients.", dbtest.Success, testID)

			//==================================================================================================================

			clients3, err := core.QueryWorkspaceClients(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace clients for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace clients for page 1.", dbtest.Success, testID)

			if len(clients3) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single client.", dbtest.Success, testID)

			clients4, err := core.QueryWorkspaceClients(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace clients for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace clients for page 2.", dbtest.Success, testID)

			if len(clients4) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single client : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single client.", dbtest.Success, testID)

			if clients3[0].ID == clients4[0].ID {
				t.Logf("\t\tTest %d:\tclient1: %v", testID, clients3[0].ID)
				t.Logf("\t\tTest %d:\tclient2: %v", testID, clients4[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different clients : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different clients.", dbtest.Success, testID)

		}
	}
}
