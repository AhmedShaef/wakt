package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers"
	"github.com/AhmedShaef/wakt/business/core/time_entry"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TimeEntryTests holds methods for each time_entry subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type TimeEntryTests struct {
	app       http.Handler
	userToken string
}

// TestTimeEntry runs a series of tests to exercise TimeEntry behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestTimeEntry(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttesttimeentry")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := TimeEntryTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postTimeEntry400", tests.postTimeEntry400)
	t.Run("startTimeEntry400", tests.startTimeEntry400)
	t.Run("postTimeEntry401", tests.postTimeEntry401)
	t.Run("startTimeEntry401", tests.startTimeEntry401)
	t.Run("getTimeEntry404", tests.getTimeEntry404)
	t.Run("getTimeEntry400", tests.getTimeEntry400)
	t.Run("deleteTimeEntryNotFound", tests.deleteTimeEntryNotFound)
	t.Run("putTimeEntry404", tests.putTimeEntry404)
	t.Run("stopTimeEntry404", tests.stopTimeEntry404)
	t.Run("crudTimeEntry", tests.crudTimeEntry)
}

// postTimeEntry400 validates a time_entry can't be created with the endpoint
// unless a valid time_entry document is submitted.
func (pt *TimeEntryTests) postTimeEntry400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new time_entry can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete time_entry value.", testID)
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
				{Field: "created_with", Error: "created_with is a required field"},
				{Field: "duration", Error: "duration is a required field"},
				{Field: "start", Error: "start is a required field"},
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

// startTimeEntry400 validates a time_entry can't be created with the endpoint
// unless a valid time_entry document is submitted.
func (pt *TimeEntryTests) startTimeEntry400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry/start", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new time_entry can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete time_entry value.", testID)
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
				{Field: "created_with", Error: "created_with is a required field"},
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

// postTimeEntry401 validates a time_entry can't be created with the endpoint
// unless the user is authenticated
func (pt *TimeEntryTests) postTimeEntry401(t *testing.T) {
	np := time_entry.NewTimeEntry{
		Start:       time.Now(),
		Duration:    time.Duration(60),
		CreatedWith: "API",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new time_entry can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete time_entry value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// startTimeEntry401 validates a time_entry can't be created with the endpoint
// unless the user is authenticated
func (pt *TimeEntryTests) startTimeEntry401(t *testing.T) {
	np := time_entry.StartTimeEntry{
		CreatedWith: "API",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new time_entry can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete time_entry value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// getTimeEntry400 validates a time_entry request for a malformed id.
func (pt *TimeEntryTests) getTimeEntry400(t *testing.T) {
	id := "12345"

	r := httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry with a malformed id.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", dbtest.Success, testID)

			got := w.Body.String()
			exp := `{"error":"ID is not in its proper form"}`
			if got != exp {
				t.Logf("\t\tTest %d:\tGot : %v", testID, got)
				t.Logf("\t\tTest %d:\tExp: %v", testID, exp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// getTimeEntry404 validates a time_entry request for a time_entry that does not exist with the endpoint.
func (pt *TimeEntryTests) getTimeEntry404(t *testing.T) {
	id := "45cf87a3-5915-4079-a9af-6c559239ddbf"

	r := httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry with an unknown id.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
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

// deleteTimeEntryNotFound validates deleting a time_entry that does not exist is not a failure.
func (pt *TimeEntryTests) deleteTimeEntryNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/time_entry/delete/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a time_entry that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putTimeEntry404 validates updating a time_entry that does not exist.
func (pt *TimeEntryTests) putTimeEntry404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := time_entry.UpdateTimeEntry{
		Start: dbtest.TimePointer(time.Now()),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/update/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a time_entry that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
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

// putTags404 validates updating a time_entry that does not exist.
func (pt *TimeEntryTests) putTags404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := time_entry.UpdateTimeEntryTags{
		Tags:    []string{"tags", "project"},
		TagMode: "add",
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a time_entry that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
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

// stopTimeEntry404 validates updating a time_entry that does not exist.
func (pt *TimeEntryTests) stopTimeEntry404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/"+id+"/stop", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a time_entry that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
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

// crudTimeEntry performs a complete test of CRUD against the api.
func (pt *TimeEntryTests) crudTimeEntry(t *testing.T) {
	p := pt.postTimeEntry201(t)
	S := pt.startTimeEntry201(t)
	defer pt.deleteTimeEntry204(t, p.ID)

	pt.getTimeEntry200(t, p.ID)
	pt.getRunTimeEntry200(t, S.ID)
	pt.getRangeTimeEntry200(t, p.ID)
	pt.getDash200(t)
	pt.putTimeEntry204(t, p.ID)
	pt.putTags204(t, p.ID)
	pt.stopTimeEntry204(t, S.ID)
}

// postTimeEntry201 validates a time_entry can be created with the endpoint.
func (pt *TimeEntryTests) postTimeEntry201(t *testing.T) time_entry.TimeEntry {
	np := time_entry.NewTimeEntry{
		Start:       time.Now(),
		Duration:    time.Duration(60),
		CreatedWith: "API",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got time_entry.TimeEntry

	t.Log("Given the need to create a new time_entry with the time_entry endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared time_entry value.", testID)
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
			exp.Wid = "7da3ca14-6366-47cf-b953-f706226567d8"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}

	return got
}

// startTimeEntry201 validates a time_entry can be created with the endpoint.
func (pt *TimeEntryTests) startTimeEntry201(t *testing.T) time_entry.TimeEntry {
	np := time_entry.StartTimeEntry{
		CreatedWith: "API",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/time_entry/start", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got time_entry.TimeEntry

	t.Log("Given the need to create a new time_entry with the time_entry endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared time_entry value.", testID)
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
			exp.Wid = "7da3ca14-6366-47cf-b953-f706226567d8"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}

	return got
}

// deleteTimeEntry200 validates deleting a time_entry that does exist.
func (pt *TimeEntryTests) deleteTimeEntry204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/time_entry/delete/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a time_entry that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new time_entry %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// getTimeEntry200 validates a time_entry request for an existing id.
func (pt *TimeEntryTests) getTimeEntry200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new time_entry %s.", testID, id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp.ID = id

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// getRunTimeEntry200 validates a time_entry request for an existing id.
func (pt *TimeEntryTests) getRunTimeEntry200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/time_entry/running/1/20", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new time_entrys.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got []time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp[0].ID = id

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// getRangeTimeEntry200 validates a time_entry request for an existing id.
func (pt *TimeEntryTests) getRangeTimeEntry200(t *testing.T, id string) {

	start := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, time.January, 2, 0, 0, 0, 0, time.UTC)

	ranges := "?start_date=" + start.Format(time.RFC3339) + "&end_date=" + end.Format(time.RFC3339)

	r := httptest.NewRequest(http.MethodGet, "/v1/time_entry/1/20"+ranges, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new time_entrys.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got []time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp[0].ID = id

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// getDash200 validates a time_entry request for an existing id.
func (pt *TimeEntryTests) getDash200(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/dashboard", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a time_entry that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new time_entrys.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got []time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got

			if diff := cmp.Diff(len(got), len(exp)); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// putTimeEntry204 validates updating a time_entry that does exist.
func (pt *TimeEntryTests) putTimeEntry204(t *testing.T, id string) {
	body := `{"created_with": "cURL"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/update/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a time_entry with the time_entry endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified time_entry value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+pt.userToken)
			pt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", dbtest.Success, testID)

			var ru time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			if ru.CreatedWith != "cURL" {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", dbtest.Failed, testID, ru.CreatedWith, "cURL")
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", dbtest.Success, testID)
		}
	}
}

// putTags204 validates updating a time_entry that does exist.
func (pt *TimeEntryTests) putTags204(t *testing.T, id string) {
	body := `{"tags": ["cURL", "golang"], "tag_mode": "add"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/tags/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a time_entry with the time_entry endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified time_entry value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+pt.userToken)
			pt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", dbtest.Success, testID)

			var ru time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			if ru.Tags[0] != "golang" {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", dbtest.Failed, testID, ru.Tags[0], "golang")
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", dbtest.Success, testID)
		}
	}
}

// stopTimeEntry204 validates updating a time_entry that does exist.
func (pt *TimeEntryTests) stopTimeEntry204(t *testing.T, id string) {
	body := `{"name": "Graphic Novels"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/time_entry/"+id+"/stop", strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a time_entry with the time_entry endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified time_entry value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/v1/time_entry/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+pt.userToken)
			pt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", dbtest.Success, testID)

			var ru time_entry.TimeEntry
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			if ru.ID != id {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", dbtest.Failed, testID, ru.ID, id)
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", dbtest.Success, testID)
		}
	}
}
