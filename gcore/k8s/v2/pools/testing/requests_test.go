package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int, clusterName string) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d/%s/pools", projectID, regionID, clusterName)
}

func prepareGetTestURLParams(projectID int, regionID int, clusterName, poolName string) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d/%s/pools/%s", projectID, regionID, clusterName, poolName)
}

func prepareActionTestURLParams(projectID int, regionID int, clusterName, poolName, action string) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d/%s/pools/%s/%s", projectID, regionID, clusterName, poolName, action)
}

func prepareListTestURL(clusterName string) string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID, clusterName)
}

func prepareGetTestURL(clusterName, poolName string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, clusterName, poolName)
}

func prepareResizeTestURL(clusterName, poolName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, poolName, "resize")
}

func prepareListInstancesTestURL(clusterName, poolName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, poolName, "instances")
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(Cluster1Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	count := 0

	err := pools.List(client, Cluster1Name).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := pools.ExtractClusterPools(page)
		require.NoError(t, err)
		ng1 := actual[0]
		require.Equal(t, Pool1, ng1)
		require.Equal(t, ExpectedClusterPoolListSlice, actual)
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

	th.Mux.HandleFunc(prepareListTestURL(Cluster1Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")

	groups, err := pools.ListAll(client, Cluster1Name)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	ng1 := groups[0]
	require.Equal(t, Pool1, ng1)
	require.Equal(t, ExpectedClusterPoolListSlice, groups)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(Cluster1Name), func(w http.ResponseWriter, r *http.Request) {
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
		Name:              "pool-1",
		FlavorID:          "g0-standard-2-4",
		MinNodeCount:      1,
		MaxNodeCount:      2,
		BootVolumeSize:    50,
		BootVolumeType:    volumes.SsdHiIops,
		ServerGroupPolicy: servergroups.AffinityPolicy,
	}
	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := pools.Create(client, Cluster1Name, options).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")

	ct, err := pools.Get(client, Cluster1Name, Pool1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Pool1, *ct)
	th.CheckDeepEquals(t, &Pool1, ct)
}

func TestUpdate(t *testing.T) {
	t.Run("EnableAutoHealing", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestJSONRequest(t, r, UpdateRequest1)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := fmt.Fprint(w, UpdateResponse1)
			if err != nil {
				log.Error(err)
			}
		})

		client := fake.ServiceTokenClient("k8s/clusters", "v2")

		enabled := true
		options := pools.UpdateOpts{
			AutoHealingEnabled: &enabled,
		}

		ct, err := pools.Update(client, Cluster1Name, Pool1.Name, options).Extract()

		require.NoError(t, err)
		require.Equal(t, true, ct.AutoHealingEnabled)
	})
	t.Run("DisableAutoHealing", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestJSONRequest(t, r, UpdateRequest2)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := fmt.Fprint(w, UpdateResponse2)
			if err != nil {
				log.Error(err)
			}
		})

		client := fake.ServiceTokenClient("k8s/clusters", "v2")

		disabled := false
		options := pools.UpdateOpts{
			AutoHealingEnabled: &disabled,
		}

		ct, err := pools.Update(client, Cluster1Name, Pool1.Name, options).Extract()

		require.NoError(t, err)
		require.Equal(t, false, ct.AutoHealingEnabled)
	})
	t.Run("NodeCount", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestJSONRequest(t, r, UpdateRequest3)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := fmt.Fprint(w, UpdateResponse3)
			if err != nil {
				log.Error(err)
			}
		})

		client := fake.ServiceTokenClient("k8s/clusters", "v2")

		options := pools.UpdateOpts{
			MinNodeCount: 2,
			MaxNodeCount: 3,
		}

		ct, err := pools.Update(client, Cluster1Name, Pool1.Name, options).Extract()

		require.NoError(t, err)
		require.Equal(t, 2, ct.MinNodeCount)
		require.Equal(t, 3, ct.MaxNodeCount)
	})
	t.Run("LabelsTaints", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestJSONRequest(t, r, UpdateRequest4)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := fmt.Fprint(w, UpdateResponse4)
			if err != nil {
				log.Error(err)
			}
		})

		client := fake.ServiceTokenClient("k8s/clusters", "v2")

		options := pools.UpdateOpts{
			Labels: &map[string]string{"foo": "bar"},
			Taints: &map[string]string{"qux": "wat:NoSchedule"},
		}

		ct, err := pools.Update(client, Cluster1Name, Pool1.Name, options).Extract()

		require.NoError(t, err)
		require.Equal(t, map[string]string{"foo": "bar"}, ct.Labels)
		require.Equal(t, map[string]string{"qux": "wat:NoSchedule"}, ct.Taints)
	})
	t.Run("ClearLabelsTaints", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		th.Mux.HandleFunc(prepareGetTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestJSONRequest(t, r, UpdateRequest5)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := fmt.Fprint(w, UpdateResponse5)
			if err != nil {
				log.Error(err)
			}
		})

		client := fake.ServiceTokenClient("k8s/clusters", "v2")

		options := pools.UpdateOpts{
			Labels: &map[string]string{},
			Taints: &map[string]string{},
		}

		ct, err := pools.Update(client, Cluster1Name, Pool1.Name, options).Extract()

		require.NoError(t, err)
		require.Equal(t, map[string]string{}, ct.Labels)
		require.Equal(t, map[string]string{}, ct.Taints)
	})
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	listenURL := prepareGetTestURL(Cluster1Name, Pool1.Name)
	th.Mux.HandleFunc(listenURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := pools.Delete(client, Cluster1Name, Pool1.Name).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestResize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResizeTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResizeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, ResizeResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := pools.ResizeOpts{
		NodeCount: 2,
	}

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := pools.Resize(client, Cluster1Name, Pool1.Name, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestInstances(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstancesTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListInstancesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	count := 0

	err := pools.ListInstances(client, Cluster1Name, Pool1.Name).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := instances.ExtractInstances(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Instance1, ct)
		require.Equal(t, ExpectedInstancesSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestInstancesAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstancesTestURL(Cluster1Name, Pool1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListInstancesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	actual, err := pools.ListInstancesAll(client, Cluster1Name, Pool1.Name)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Instance1, ct)
	require.Equal(t, ExpectedInstancesSlice, actual)
}
