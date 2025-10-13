package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/dbaas/postgres/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		q := r.URL.Query()
		resp := ListResponse
		if q.Get("limit") != "" || q.Get("offset") != "" {
			resp = ListResponsePagination
		}
		if _, err := fmt.Fprint(w, resp); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	pages, err := clusters.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	actual, err := clusters.ExtractClusters(pages)
	th.AssertNoErr(t, err)

	expected := []clusters.PostgresSQLClusterShort{FirstClusterShort, SecondClusterShort}

	th.CheckDeepEquals(t, expected, actual)
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		q := r.URL.Query()
		resp := ListResponse
		if q.Get("limit") != "" || q.Get("offset") != "" {
			resp = ListResponsePagination
		}
		if _, err := fmt.Fprint(w, resp); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	all, err := clusters.ListAll(client, nil)
	th.AssertNoErr(t, err)

	expected := []clusters.PostgresSQLClusterShort{FirstClusterShort, SecondClusterShort}
	th.CheckDeepEquals(t, expected, all)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareClusterTestURL(FirstClusterShort.ClusterName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	actual, err := clusters.Get(client, FirstClusterShort.ClusterName).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &FirstClusterDetail, actual)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateResponse)
	})

	createOpts := clusters.CreateOpts{
		ClusterName: "test-cluster-1",
		Databases: []clusters.DatabaseOpts{
			{
				Name:  "testdb",
				Owner: "testuser",
			},
		},
		Flavor: clusters.FlavorOpts{
			CPU:       2,
			MemoryGiB: 4,
		},
		HighAvailability: &clusters.HighAvailabilityOpts{
			ReplicationMode: clusters.HighAvailabilityReplicationModeAsync,
		},
		Network: clusters.NetworkOpts{
			ACL:         []string{"0.0.0.0/0"},
			NetworkType: "public",
		},
		PGServerConfiguration: clusters.PGServerConfigurationOpts{
			PGConf:  "standard",
			Version: "15.0",
			Pooler: &clusters.PoolerOpts{
				Mode: clusters.PoolerModeSession,
				Type: "pgbouncer",
			},
		},
		Storage: clusters.PGStorageConfigurationOpts{
			SizeGiB: 50,
			Type:    "standard",
		},
		Users: []clusters.PgUserOpts{
			{
				Name:           "testuser",
				RoleAttributes: []clusters.RoleAttribute{clusters.RoleAttributeLogin, clusters.RoleAttributeCreateDB},
			},
		},
	}

	res := clusters.Create(client, createOpts)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedTaskResults, task)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareClusterTestURL(FirstClusterShort.ClusterName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	updateOpts := clusters.UpdateOpts{
		Flavor: &clusters.FlavorOpts{
			CPU:       4,
			MemoryGiB: 8,
		},
		Storage: &clusters.PGStorageConfigurationUpdateOpts{
			SizeGiB: 100,
		},
	}

	res := clusters.Update(client, FirstClusterShort.ClusterName, updateOpts)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedTaskResults, task)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	th.Mux.HandleFunc(prepareClusterTestURL(FirstClusterShort.ClusterName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, DeleteResponse)
	})

	res := clusters.Delete(client, FirstClusterShort.ClusterName, nil)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	var expectedTasks = tasks.TaskResults{
		Tasks: []tasks.TaskID{
			"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}
	th.CheckDeepEquals(t, &expectedTasks, task)
}

func TestCreateValidation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	// Test with missing required fields
	createOpts := clusters.CreateOpts{
		ClusterName: "test-cluster",
		// Missing required fields
	}

	res := clusters.Create(client, createOpts)
	if res.Err == nil {
		t.Fatal("Expected error for invalid CreateOpts, got nil")
	}
}

func TestUpdateValidation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("dbaas/postgres/clusters", "v1")

	// Test with empty UpdateOpts (should return error)
	updateOpts := clusters.UpdateOpts{}

	res := clusters.Update(client, FirstClusterShort.ClusterName, updateOpts)
	if res.Err == nil {
		t.Fatal("Expected error for empty UpdateOpts, got nil")
	}
}
