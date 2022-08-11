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
	"github.com/AhmedShaef/wakt/business/core/group"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// GroupTests holds methods for each group subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type GroupTests struct {
	app       http.Handler
	userToken string
}

// TestGroups runs a series of tests to exercise Group behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestGroups(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttestgroup")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := GroupTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postGroup400", tests.postGroup400)
	t.Run("postGroup401", tests.postGroup401)
	t.Run("deleteGroupNotFound", tests.deleteGroupNotFound)
	t.Run("putGroup404", tests.putGroup404)
	t.Run("crudGroups", tests.crudGroup)
}

// postGroup400 validates a group can't be created with the endpoint
// unless a valid group document is submitted.
func (pt *GroupTests) postGroup400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/group", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new group can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete group value.", testID)
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
				{Field: "name", Error: "name is a required field"},
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

// postGroup401 validates a group can't be created with the endpoint
// unless the user is authenticated
func (pt *GroupTests) postGroup401(t *testing.T) {
	np := group.NewGroup{
		Name: "Comic Books",
		WID:  "7da3ca14-6366-47cf-b953-f706226567d8",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/group", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new group can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete group value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// deleteGroupNotFound validates deleting a group that does not exist is not a failure.
func (pt *GroupTests) deleteGroupNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/group/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a group that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new group %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putGroup404 validates updating a group that does not exist.
func (pt *GroupTests) putGroup404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := group.UpdateGroup{
		Name: dbtest.StringPointer("Nonexistent"),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/group/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a group that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new group %s.", testID, id)
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

// crudGroup performs a complete test of CRUD against the api.
func (pt *GroupTests) crudGroup(t *testing.T) {
	p := pt.postGroup201(t)
	defer pt.deleteGroup204(t, p.ID)

	pt.putGroup204(t, p.ID)
}

// postGroup201 validates a group can be created with the endpoint.
func (pt *GroupTests) postGroup201(t *testing.T) group.Group {
	np := group.NewGroup{
		Name: "Comic Books",
		WID:  "7da3ca14-6366-47cf-b953-f706226567d8",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/group", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got group.Group

	t.Log("Given the need to create a new group with the groups endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared group value.", testID)
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 201 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 201 for the response.", dbtest.Success, testID)

			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like ID and Dates so we copy p.
			exp := got
			exp.Name = "Comic Books"
			exp.WID = "7da3ca14-6366-47cf-b953-f706226567d8"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}

	return got
}

// deleteGroup200 validates deleting a group that does exist.
func (pt *GroupTests) deleteGroup204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/group/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a group that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new group %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putGroup204 validates updating a group that does exist.
func (pt *GroupTests) putGroup204(t *testing.T, id string) {
	body := `{"name": "Graphic Novels"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/group/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a group with the groups endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified group value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}
