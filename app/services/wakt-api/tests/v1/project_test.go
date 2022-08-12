package tests

import (
	"bytes"
	"encoding/json"
	"github.com/AhmedShaef/wakt/business/core/project"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/AhmedShaef/wakt/app/services/wakt-api/handlers"
	"github.com/AhmedShaef/wakt/business/data/dbtest"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// ProjectTests holds methods for each project subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ProjectTests struct {
	app       http.Handler
	userToken string
}

// TestProjects runs a series of tests to exercise Project behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestProjects(t *testing.T) {
	t.Parallel()

	test := dbtest.NewIntegration(t, c, "inttestproject")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := ProjectTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postProject400", tests.postProject400)
	t.Run("postProject401", tests.postProject401)
	t.Run("getProject404", tests.getProject404)
	t.Run("getProject400", tests.getProject400)
	t.Run("deleteProjectNotFound", tests.deleteProjectNotFound)
	t.Run("putProject404", tests.putProject404)
	t.Run("crudProjects", tests.crudProject)
}

// postProject400 validates a project can't be created with the endpoint
// unless a valid project document is submitted.
func (pt *ProjectTests) postProject400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/project", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new project can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete project value.", testID)
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

// postProject401 validates a project can't be created with the endpoint
// unless the user is authenticated
func (pt *ProjectTests) postProject401(t *testing.T) {
	np := project.NewProject{
		Name: "Comic Books",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/project", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new project can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete project value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

// getProject400 validates a project request for a malformed id.
func (pt *ProjectTests) getProject400(t *testing.T) {
	id := "12345"

	r := httptest.NewRequest(http.MethodGet, "/v1/project/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a project with a malformed id.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new project %s.", testID, id)
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

// getProject404 validates a project request for a project that does not exist with the endpoint.
func (pt *ProjectTests) getProject404(t *testing.T) {
	id := "7da3ca14-6366-47cf-b953-f706226567d8"

	r := httptest.NewRequest(http.MethodGet, "/v1/project/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a project with an unknown id.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new project %s.", testID, id)
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

// deleteProjectNotFound validates deleting a project that does not exist is not a failure.
func (pt *ProjectTests) deleteProjectNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/project/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a project that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new project %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putProject404 validates updating a project that does not exist.
func (pt *ProjectTests) putProject404(t *testing.T) {
	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	up := project.UpdateProject{
		Name: dbtest.StringPointer("Nonexistent"),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/project/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a project that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new project %s.", testID, id)
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

// crudProject performs a complete test of CRUD against the api.
func (pt *ProjectTests) crudProject(t *testing.T) {
	p := pt.postProject201(t)
	defer pt.deleteProject204(t, p.ID)

	pt.getProject200(t, p.ID)
	pt.getProjectTask200(t, p.ID)
	pt.putProject204(t, p.ID)
}

// postProject201 validates a project can be created with the endpoint.
func (pt *ProjectTests) postProject201(t *testing.T) project.Project {
	np := project.NewProject{
		Name: "Comic Books",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/project", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	// This needs to be returned for other dbtest.
	var got project.Project

	t.Log("Given the need to create a new project with the projects endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared project value.", testID)
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

// deleteProject200 validates deleting a project that does exist.
func (pt *ProjectTests) deleteProject204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/project/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a project that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new project %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// getProject200 validates a project request for an existing id.
func (pt *ProjectTests) getProject200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/project/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a project that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new project %s.", testID, id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got project.Project
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp.ID = id
			exp.Name = "Comic Books"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// getProjectTask200 validates a project request for an existing id.
func (pt *ProjectTests) getProjectTask200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/project/"+id+"/task/1/20", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a project project that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new project %s.", testID, id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got []project.Project
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}
		}
	}
}

// putProject204 validates updating a project that does exist.
func (pt *ProjectTests) putProject204(t *testing.T, id string) {
	body := `{"name": "Graphic Novels"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/project/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)
	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update a project with the projects endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified project value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/v1/project/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+pt.userToken)
			pt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", dbtest.Success, testID)

			var ru project.Project
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			if ru.Name != "Graphic Novels" {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", dbtest.Failed, testID, ru.Name, "Graphic Novels")
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", dbtest.Success, testID)
		}
	}
}
