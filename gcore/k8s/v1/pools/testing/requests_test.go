package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int, clusterID string) string {
	return fmt.Sprintf("/v1/k8s/clusters/%d/%d/%s/pools", projectID, regionID, clusterID)
}

func prepareGetTestURLParams(projectID int, regionID int, clusterID, id string) string {
	return fmt.Sprintf("/v1/k8s/clusters/%d/%d/%s/pools/%s", projectID, regionID, clusterID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID, Pool1.ClusterID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, Pool1.ClusterID, id)
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

	client := fake.ServiceTokenClient("k8s/clusters", "v1")
	count := 0

	err := pools.List(client, Pool1.ClusterID, pools.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := pools.ExtractClusterPools(page)
		require.NoError(t, err)
		ng1 := actual[0]
		require.Equal(t, PoolList1, ng1)
		require.Equal(t, ExpectedClusterListPoolSlice, actual)
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

	client := fake.ServiceTokenClient("k8s/clusters", "v1")

	groups, err := pools.ListAll(client, Pool1.ClusterID, nil)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	ng1 := groups[0]
	require.Equal(t, PoolList1, ng1)
	require.Equal(t, ExpectedClusterListPoolSlice, groups)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Pool1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse1)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v1")

	ct, err := pools.Get(client, Pool1.ClusterID, Pool1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, Pool1, *ct)
	th.CheckDeepEquals(t, &Pool1, ct)

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

	options := pools.CreateOpts{
		Name:             Pool1.Name,
		FlavorID:         Pool1.FlavorID,
		NodeCount:        1,
		DockerVolumeSize: 5,
	}
	client := fake.ServiceTokenClient("k8s/clusters", "v1")
	tasks, err := pools.Create(client, Pool1.ClusterID, options).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	listenURL := prepareGetTestURL(Pool1.UUID)
	th.Mux.HandleFunc(listenURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v1")
	tasks, err := pools.Delete(client, Pool1.ClusterID, Pool1.UUID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Pool1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v1")

	options := pools.UpdateOpts{
		Name:         "",
		MinNodeCount: 3,
		MaxNodeCount: 4,
	}

	tasks, err := pools.Update(client, Pool1.ClusterID, Pool1.UUID, options).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}
