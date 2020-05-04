package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/pools"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/clusters"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/k8s/%d/%d/clusters", projectID, regionID)
}

func prepareVersionURL() string {
	return fmt.Sprintf("/v1/k8s/versions")
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/k8s/%d/%d/clusters/%s", projectID, regionID, id)
}

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/k8s/%d/%d/clusters/%s/%s", projectID, regionID, id, action)
}

func prepareClusterTestURLParams(projectID int, regionID int, id, name string) string {
	return fmt.Sprintf("/v1/k8s/%d/%d/clusters/%s/%s", projectID, regionID, id, name)
}
func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareGetConfigTestURL(id string) string {
	return prepareClusterTestURLParams(fake.ProjectID, fake.RegionID, id, "config")
}

func prepareResizeTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "resize")
}

func prepareUpgradeTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "upgrade")
}

func prepareClusterCertificateURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "certificates")
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

	client := fake.ServiceTokenClient("k8s", "v1")
	count := 0

	err := clusters.List(client, clusters.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, ClusterList1, ct)
		require.Equal(t, ExpectedClusterSlice, actual)
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

	client := fake.ServiceTokenClient("k8s", "v1")
	actual, err := clusters.ListAll(client, clusters.ListOpts{})
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, ClusterList1, ct)
	require.Equal(t, ExpectedClusterSlice, actual)
	require.NotNil(t, ct.HealthStatus)
}

func TestVersionAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareVersionURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, VersionResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")
	actual, err := clusters.VersionsAll(client)
	require.NoError(t, err)
	require.Equal(t, []string{"1.14", "1.17"}, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")

	ct, err := clusters.Get(client, Cluster1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, Cluster1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)
	require.NotNil(t, ct.HealthStatus)

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

	options := clusters.CreateOpts{
		Name:                      Cluster1.Name,
		FixedNetwork:              fixedNetwork,
		FixedSubnet:               fixedSubnet,
		MasterCount:               1,
		KeyPair:                   "keypair",
		AutoHealingEnabled:        false,
		MasterLBFloatingIPEnabled: false,
		Version:                   version,
		Pools: []pools.CreateOpts{{
			Name:             Cluster1.Pools[0].Name,
			FlavorID:         Cluster1.Pools[0].FlavorID,
			NodeCount:        Cluster1.Pools[0].NodeCount,
			DockerVolumeSize: 10,
			MinNodeCount:     Cluster1.Pools[0].MinNodeCount,
			MaxNodeCount:     Cluster1.Pools[0].MinNodeCount + 1,
		}},
	}

	client := fake.ServiceTokenClient("k8s", "v1")
	tasks, err := clusters.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")
	tasks, err := clusters.Delete(client, Cluster1.UUID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestResize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResizeTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
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

	options := clusters.ResizeOpts{
		NodeCount:     2,
		NodesToRemove: nil,
		Pool:          pool,
	}

	client := fake.ServiceTokenClient("k8s", "v1")
	tasks, err := clusters.Resize(client, Cluster1.UUID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestGetConfig(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetConfigTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, string(ConfigResponse))
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")

	cfg, err := clusters.GetConfig(client, Cluster1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, Config1, *cfg)

}

func TestUpgrade(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUpgradeTestURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpgradeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, UpgradeResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := clusters.UpgradeOpts{
		Version: version,
		Pool:    pool,
	}

	client := fake.ServiceTokenClient("k8s", "v1")
	tasks, err := clusters.Upgrade(client, Cluster1.UUID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestGetCACertificate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareClusterCertificateURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ClusterCAResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")

	ct, err := clusters.Certificate(client, Cluster1.UUID).Extract()

	require.NoError(t, err)
	require.Equal(t, ClusterCertificate, *ct)

}

func TestSignCertificate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareClusterCertificateURL(Cluster1.UUID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ClusterCsrRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ClusterSignResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s", "v1")

	opts := clusters.ClusterSignCertificateOpts{
		CSR: "string",
	}

	ct, err := clusters.SignCertificate(client, Cluster1.UUID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, ClusterSignedCertificate, *ct)

}
