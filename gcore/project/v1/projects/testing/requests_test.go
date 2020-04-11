package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/project/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/project/v1/projects"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURL() string {
	return "/v1/projects"
}

func prepareGetTestURL(id int) string {
	return fmt.Sprintf("/v1/projects/%d", id)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("projects", "v1")
	count := 0

	err := projects.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := projects.ExtractProjects(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Project1, ct)
		require.Equal(t, ExpectedProjectSlice, actual)
		return true, nil
	})

	require.NoError(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("projects", "v1")

	results, err := projects.ListAll(client)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, Project1, ct)
	require.Equal(t, ExpectedProjectSlice, results)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Project1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("projects", "v1")

	ct, err := projects.Get(client, Project1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Project1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)

}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := projects.CreateOpts{
		ClientID:    1,
		State:       types.ProjectStateActive,
		Name:        "default",
		Description: "",
	}

	err := gcorecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	client := fake.ServiceTokenClient("projects", "v1")
	project, err := projects.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Project1, *project)
	require.Equal(t, createdTime, project.CreatedAt)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Project1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UpdateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("projects", "v1")

	options := projects.UpdateOpts{
		Name:        Project1.Name,
		Description: "description",
	}

	err := gcorecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	project, err := projects.Update(client, Project1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Project1, *project)
	require.Equal(t, createdTime, project.CreatedAt)

}
