package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/region/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/region/v1/regions"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURL() string {
	return "/v1/regions"
}

func prepareGetTestURL(id int) string {
	return fmt.Sprintf("/v1/regions/%d", id)
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

	client := fake.ServiceTokenClient("regions", "v1")
	count := 0

	err := regions.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := regions.ExtractRegions(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Region1, ct)
		require.Equal(t, ExpectedRegionSlice, actual)
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

	client := fake.ServiceTokenClient("regions", "v1")

	results, err := regions.ListAll(client)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, Region1, ct)
	require.Equal(t, ExpectedRegionSlice, results)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Region1.ID)

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

	client := fake.ServiceTokenClient("regions", "v1")

	ct, err := regions.Get(client, Region1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Region1, *ct)
	require.Equal(t, createdTime, ct.CreatedOn)

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

	options := regions.CreateOpts{
		DisplayName:       Region1.DisplayName,
		KeystoneName:      Region1.DisplayName,
		State:             types.RegionStateActive,
		EndpointType:      Region1.EndpointType,
		ExternalNetworkID: Region1.ExternalNetworkID,
		SpiceProxyURL:     nil,
		KeystoneID:        Region1.KeystoneID,
	}

	err := gcorecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	client := fake.ServiceTokenClient("regions", "v1")
	region, err := regions.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Region1, *region)
	require.Equal(t, createdTime, region.CreatedOn)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Region1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
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

	client := fake.ServiceTokenClient("regions", "v1")

	options := regions.UpdateOpts{
		DisplayName:       Region1.DisplayName,
		State:             types.RegionStateDeleted,
		EndpointType:      "",
		ExternalNetworkID: "",
		SpiceProxyURL:     nil,
	}

	err := gcorecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	region, err := regions.Update(client, Region1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Region1, *region)
	require.Equal(t, createdTime, region.CreatedOn)

}
