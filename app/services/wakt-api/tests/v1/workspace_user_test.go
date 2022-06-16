package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// WorkspaceUserTests holds methods for each workspace_user subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type WorkspaceUserTests struct {
	app       http.Handler
	userToken string
}

// TestWorkspaceUsers runs a series of tests to exercise WorkspaceUser behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestWorkspaceUsers(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttestworkspace_user")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := WorkspaceUserTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postWorkspaceUser400", tests.postWorkspaceUser400)
	t.Run("postWorkspaceUser401", tests.postWorkspaceUser401)
	t.Run("deleteWorkspaceUserNotFound", tests.deleteWorkspaceUserNotFound)
	t.Run("putWorkspaceUser404", tests.putWorkspaceUser404)
	t.Run("crudWorkspaceUsers", tests.crudWorkspaceUser)
}

// postWorkspaceUser400 validates a workspace_user can't be created with the endpoint
// unless a valid workspace_user document is submitted.
func (pt *WorkspaceUserTests) postWorkspaceUser400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/workspace_user", strings.NewReader(`{"inviter_id": "32c1494f-1c1f-4981-857f-b0526cb654ec"}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new workspace_user can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete workspace_user value.", testID)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", dbtest.Success, testID)

			// Inspect the response.
			var got v1Web.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type : %v", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type.", dbtest.Success, testID)

			fields := validate.FieldErrors{
				{Field: "emails", Error: "emails is a required field"},
			}
			exp := v1Web.ErrorResponse{
				Error:  "data validation error",
				Fields: fields.Fields(),
			}

			// We can't rely on the order of the field errors so they have to be
			// sorted. Tell the cmp package how to sort them.
			sorter := cmpopts.SortSlices(func(a, b validate.FieldError) bool {
				return a.Field < b.Field
			})

			if diff := cmp.Diff(got, exp, sorter); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// postWorkspaceUser401 validates a workspace_user can't be created with the endpoint
// unless the user is authenticated
func (pt *WorkspaceUserTests) postWorkspaceUser401(t *testing.T) {
	np := workspace_user.InviteUsers{
		Emails:    []string{"example@example.com"},
		InviterID: "5cf37266-3473-4006-984f-9325122678b7",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/workspace_user", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new workspace_user can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete workspace_user value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// deleteWorkspaceUserNotFound validates deleting a workspace_user that does not exist is not a failure.
func (pt *WorkspaceUserTests) deleteWorkspaceUserNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/workspace_user/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a workspace_user that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new workspace_user %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putWorkspaceUser404 validates updating a workspace_user that does not exist.
func (pt *WorkspaceUserTests) putWorkspaceUser404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := workspace_user.UpdateWorkspaceUser{
		Active: dbtest.BoolPointer(true),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/workspace_user/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a workspace_user that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new workspace_user %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", dbtest.Success, testID)

			got := w.Body.String()
			exp := "not found"
			if !strings.Contains(got, exp) {
				t.Logf("\t\tTest %d:\tGot : %v", testID, got)
				t.Logf("\t\tTest %d:\tExp: %v", testID, exp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// crudWorkspaceUser performs a complete test of CRUD against the api.
func (pt *WorkspaceUserTests) crudWorkspaceUser(t *testing.T) {
	_ = pt.postWorkspaceUser201(t)
	defer pt.deleteWorkspaceUser204(t, "32c1494f-1c1f-4981-857f-b0526cb654ec")

	pt.putWorkspaceUser204(t, "32c1494f-1c1f-4981-857f-b0526cb654ec")
}

// postWorkspaceUser201 validates a workspace_user can be created with the endpoint.
func (pt *WorkspaceUserTests) postWorkspaceUser201(t *testing.T) []workspace_user.WorkspaceUser {
	np := workspace_user.InviteUsers{
		Emails:    []string{"example1@example.com", "example2@example.com"},
		InviterID: "32c1494f-1c1f-4981-857f-b0526cb654ec",
	}
	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/workspace_user", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got []workspace_user.WorkspaceUser

	t.Log("Given the need to create a new workspace_user with the workspace_users endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared workspace_user value.", testID)
		{
			if w.Code != http.StatusCreated {
				t.Errorf("\t%s\tTest %d:\tShould receive a status code of 201 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 201 for the response.", dbtest.Success, testID)

			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Errorf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like ID and Dates so we copy p.
			exp := got
			exp[0].Active = false

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Errorf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}

	return got
}

// deleteWorkspaceUser200 validates deleting a workspace_user that does exist.
func (pt *WorkspaceUserTests) deleteWorkspaceUser204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/workspace_user/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a workspace_user that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new workspace_user %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putWorkspaceUser204 validates updating a workspace_user that does exist.
func (pt *WorkspaceUserTests) putWorkspaceUser204(t *testing.T, id string) {
	body := `{"active": true}`
	r := httptest.NewRequest(http.MethodPut, "/v1/workspace_user/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a workspace_user with the workspace_users endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified workspace_user value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}
