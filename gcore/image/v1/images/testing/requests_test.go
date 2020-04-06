package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/image/v1/images"

	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/images/%d/%d", version, projectID, regionID)
}

func prepareGetTestURLParams(version string, projectID int, regionID int, id string) string {
	return fmt.Sprintf("/%s/images/%d/%d/%s", version, projectID, regionID, id)
}

func prepareCreateTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/downloadimage/%d/%d", version, projectID, regionID)
}

func prepareListTestURL() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareDeleteTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareCreateTestURL() string {
	return prepareCreateTestURLParams("v1", fake.ProjectID, fake.RegionID)
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

	client := fake.ServiceTokenClient("images", "v1")
	count := 0

	opts := images.ListOpts{}

	err := images.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := images.ExtractImages(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Image1, ct)
		require.Equal(t, ExpectedImagesSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

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

	client := fake.ServiceTokenClient("images", "v1")

	opts := images.ListOpts{}

	actual, err := images.ListAll(client, opts)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Image1, ct)
	require.Equal(t, ExpectedImagesSlice, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Image1.ID)

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

	client := fake.ServiceTokenClient("images", "v1")

	ct, err := images.Get(client, Image1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Image1, *ct)

}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateTestURL(), func(w http.ResponseWriter, r *http.Request) {
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

	options := images.CreateOpts{
		URL:       "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img",
		Name:      "image_name",
		CowFormat: false,
	}

	err := options.Validate()
	require.NoError(t, err)

	client := fake.ServiceTokenClient("downloadimage", "v1")
	tasks, err := images.Create(client, options).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeleteTestURL(Image1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("images", "v1")
	tasks, err := images.Delete(client, Image1.ID).ExtractTasks()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}
