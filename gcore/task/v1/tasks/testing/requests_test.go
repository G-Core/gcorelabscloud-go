package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/tasks/%d/%d/active", projectID, regionID)
}

func prepareGetTestURLParams(id string) string {
	return fmt.Sprintf("/v1/tasks/%s", id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(id)
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

	client := fake.ServiceTokenClient("tasks", "v1")
	count := 0

	err := tasks.ListActive(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := tasks.ExtractTasks(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Task1, ct)
		require.Equal(t, ExpectedTasks, actual)
		require.Nil(t, ct.UpdatedOn)
		return true, nil
	})

	require.NoError(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Task1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("tasks", "v1")
	ct, err := tasks.Get(client, Task1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Task1, *ct)
	require.Equal(t, createdTime, ct.CreatedOn)
	require.Nil(t, ct.UpdatedOn)
}

func TestWaitTask(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Task1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, FinishedTaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("tasks", "v1")

	err := tasks.WaitTaskAndProcessResult(client, tasks.TaskID(Task1.ID), true, 600, func(task tasks.TaskID) error {
		require.Equal(t, Task1.ID, string(task))
		return nil
	})

	require.NoError(t, err)
}
