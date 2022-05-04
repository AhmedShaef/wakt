package user

import (
	"context"
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

func TestUser(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testuser")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with user records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single user.", testID)
		{
			ctx := context.Background()
			now := time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC)

			nu := NewUser{
				Email:    "example@example.com",
				Password: "f51gr6g5+erg4+51+erg",
			}

			user, err := core.Create(ctx, nu, now, "", "")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", dbtest.Success, testID)

			_, err = core.Authenticate(ctx, now, nu.Email, nu.Password)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to authenticate user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to authenticate user.", dbtest.Success, testID)

			//claim := auth.Claims{
			//	RegisteredClaims: jwt.RegisteredClaims{
			//		Subject: user.ID,
			//		Issuer:  "wakt project",
			//	},
			//	//Roles: []string{auth.RoleAdmin},
			//}
			//if diff := cmp.Diff(claims, claim); diff != "" {
			//	t.Fatalf("\t%s\tTest %d:\tShould get back the Auth. Diff:\n%s", dbtest.Failed, testID, diff)
			//}
			t.Logf("\t%s\tTest %d:\tShould get back the Auth.", dbtest.Success, testID)
			saved, err := core.QueryByID(ctx, user.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(user, saved); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

			upd := UpdateUser{
				FullName: dbtest.StringPointer("Shehab Shaef"),
			}

			if err := core.Update(ctx, user.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", dbtest.Success, testID)

			if saved.FullName == *upd.FullName {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.FullName)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.FullName)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", dbtest.Success, testID)
			}

			chp := ChangePassword{
				OldPassword: "f51gr6g5+erg4+51+erg",
				Password:    dbtest.StringPointer("f51gr6g5+erg4+51+erg2"),
			}
			if err := core.ChangePassword(ctx, user.ID, chp, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to change password : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to change password.", dbtest.Success, testID)

			saved2, err := core.QueryByID(ctx, user.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(user, saved); diff == "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

			if diff := cmp.Diff(saved.PasswordHash, saved2.PasswordHash); diff == "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same password. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to see updates to password.", dbtest.Success, testID)

			t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", dbtest.Success, testID)
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

	t.Log("Given the need to page through user records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen paging through 2 users.", testID)
		{
			ctx := context.Background()

			user1, err := core.QueryWorkspaceUsers(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace users for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace users for page 1.", dbtest.Success, testID)

			if len(user1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", dbtest.Success, testID)

			user2, err := core.QueryWorkspaceUsers(ctx, "7da3ca14-6366-47cf-b953-f706226567d8", 2, 2)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve workspace users for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve workspace users for page 2.", dbtest.Success, testID)

			if len(user2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", dbtest.Success, testID)

			if user1[0].ID == user2[0].ID {
				t.Logf("\t\tTest %d:\tuser1: %v", testID, user1[0].ID)
				t.Logf("\t\tTest %d:\tuser2: %v", testID, user2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different users : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different users.", dbtest.Success, testID)

			//===================================================================================================================

			user3, err := core.QueryProjectUsers(ctx, "45cf87a3-5915-4079-a9af-6c559239ddbf", 1, 1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve project users for page 1 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve project users for page 1.", dbtest.Success, testID)

			if len(user3) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", dbtest.Success, testID)

			user4, err := core.QueryProjectUsers(ctx, "45cf87a3-5915-4079-a9af-6c559239ddbf", 2, 2)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve project users for page 2 : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve project users for page 2.", dbtest.Success, testID)

			if len(user4) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single user.", dbtest.Success, testID)

			if user3[0].ID == user4[0].ID {
				t.Logf("\t\tTest %d:\tuser3: %v", testID, user3[0].ID)
				t.Logf("\t\tTest %d:\tuser4: %v", testID, user4[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different users : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different users.", dbtest.Success, testID)

		}
	}
}
