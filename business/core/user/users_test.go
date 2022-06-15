package user

import (
	"context"
	"fmt"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	"github.com/AhmedShaef/wakt/foundation/docker"
	"github.com/golang-jwt/jwt/v4"
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

			authenticate, err := core.Authenticate(ctx, nu.Email, nu.Password)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to authenticate user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to authenticate user.", dbtest.Success, testID)

			claim := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   user.ID,
					Issuer:    "wakt project",
					ExpiresAt: authenticate.ExpiresAt,
					IssuedAt:  authenticate.IssuedAt,
				},
			}
			if diff := cmp.Diff(authenticate, claim); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the Auth. Diff:\n%s", dbtest.Failed, testID, diff)
			}
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

			saved2, err := core.QueryByEmail(ctx, "example@example.com")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(user, saved); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

			if diff := cmp.Diff(saved.PasswordHash, saved2.PasswordHash); diff == "" {
				t.Errorf("\t%s\tTest %d:\tShould get back the same password. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", dbtest.Success, testID)

			ui := UpdateImage{
				ImageName: "image-785278.png",
			}
			if err := core.UpdateImage(ctx, user.ID, ui, now); err != nil {
				t.Errorf("\t%s\tTest %d:\tShould be able to update the image : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update the image.", dbtest.Success, testID)

		}
	}
}
