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
	"github.com/AhmedShaef/wakt/business/core/team"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TeamTests holds methods for each team subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type TeamTests struct {
	app        http.Handler
	adminToken string
}

// TestTeams runs a series of tests to exercise Team behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestTeams(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttestteam")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := TeamTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		adminToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postTeam400", tests.postTeam400)
	t.Run("postTeam403", tests.postTeam403)
	t.Run("postTeam401", tests.postTeam401)
	t.Run("deleteTeamNotFound", tests.deleteTeamNotFound)
	t.Run("putTeam404", tests.putTeam404)
	t.Run("crudTeams", tests.crudTeam)
}

// postTeam400 validates a team can't be created with the endpoint
// unless a valid team document is submitted.
func (pt *TeamTests) postTeam400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/team", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new team can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete team value.", testID)
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

			exp := v1Web.ErrorResponse{
				Error: "ID is not in its proper form",
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

// postTeam403 validates a team can't be created with the endpoint
// unless the user is authenticated
func (pt *TeamTests) postTeam403(t *testing.T) {
	np := team.NewTeam{
		Pid:  "45cf87a3-5915-4079-a9af-6c559239ddbf",
		UID:  "5cf37266-3473-4006-984f-9325122678b7",
		Puis: "c7142720-91d3-4d1e-841d-680042b6500c",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/team", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new team can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete team value.", testID)
		{
			if w.Code != http.StatusForbidden {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 403 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// postTeam401 validates a team can't be created with the endpoint
// unless the user is authenticated
func (pt *TeamTests) postTeam401(t *testing.T) {
	np := team.NewTeam{
		Pid: "45cf87a3-5915-4079-a9af-6c559239ddbf",
		UID: "5cf37266-3473-4006-984f-9325122678b7",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/team", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new team can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete team value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// deleteTeamNotFound validates deleting a team that does not exist is not a failure.
func (pt *TeamTests) deleteTeamNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/team/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a team that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new team %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putTeam404 validates updating a team that does not exist.
func (pt *TeamTests) putTeam404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := team.UpdateTeam{
		Manager: dbtest.BoolPointer(true),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/team/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a team that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new team %s.", testID, id)
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

// crudTeam performs a complete test of CRUD against the api.
func (pt *TeamTests) crudTeam(t *testing.T) {
	p := pt.postTeam201(t)
	defer pt.deleteTeam204(t, p[0].ID)

	pt.putTeam204(t, p[0].ID)
}

// postTeam201 validates a team can be created with the endpoint.
func (pt *TeamTests) postTeam201(t *testing.T) []team.Team {
	np := team.NewTeam{
		Pid:     "45cf87a3-5915-4079-a9af-6c559239ddbf",
		UID:     "5cf37266-3473-4006-984f-9325122678b7",
		Puis:    "efcc74aa-86d2-4e11-80f9-3ca912af8269",
		Manager: true,
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/team", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got []team.Team

	t.Log("Given the need to create a new team with the teams endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared team value.", testID)
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
			exp[0].WID = "7da3ca14-6366-47cf-b953-f706226567d8"
			exp[0].UID = "5cf37266-3473-4006-984f-9325122678b7"
			exp[0].Pid = "45cf87a3-5915-4079-a9af-6c559239ddbf"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}

	return got
}

// deleteTeam200 validates deleting a team that does exist.
func (pt *TeamTests) deleteTeam204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/team/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a team that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new team %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putTeam204 validates updating a team that does exist.
func (pt *TeamTests) putTeam204(t *testing.T, id string) {
	body := `{"rate": 10}`
	r := httptest.NewRequest(http.MethodPut, "/v1/team/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.adminToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a team with the teams endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified team value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}
