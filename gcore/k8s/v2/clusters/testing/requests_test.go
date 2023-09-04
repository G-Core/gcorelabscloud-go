package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, clusterName string) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d/%s", projectID, regionID, clusterName)
}

func prepareActionTestURLParams(projectID int, regionID int, clusterName, action string) string {
	return fmt.Sprintf("/v2/k8s/clusters/%d/%d/%s/%s", projectID, regionID, clusterName, action)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(clusterName string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, clusterName)
}

func prepareGetCertificateTestURL(clusterName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, "certificates")
}

func prepareGetConfigTestURL(clusterName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, "config")
}

func prepareListInstancesTestURL(clusterName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, "instances")
}

func prepareUpgradeTestURL(clusterName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, clusterName, "upgrade")
}

func prepareVersionTestURL() string {
	return "/v2/k8s/versions"
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

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	count := 0

	err := clusters.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Cluster1, ct)
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

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	actual, err := clusters.ListAll(client)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Cluster1, ct)
	require.Equal(t, ExpectedClusterSlice, actual)
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
		Name:         "cluster-1",
		KeyPair:      "keypair",
		FixedNetwork: fixedNetwork,
		FixedSubnet:  fixedSubnet,
		Version:      "v1.26.7",
		Pools: []pools.CreateOpts{
			{
				Name:               "pool-1",
				FlavorID:           "g0-standard-2-4",
				MinNodeCount:       1,
				MaxNodeCount:       2,
				BootVolumeType:     volumes.SsdHiIops,
				BootVolumeSize:     50,
				AutoHealingEnabled: true,
			},
		},
	}

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := clusters.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
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

	ct, err := clusters.Get(client, Cluster1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Cluster1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := clusters.Delete(client, Cluster1.Name).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestGetCertificate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetCertificateTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, string(CertificatesResponse))
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")

	cert, err := clusters.GetCertificate(client, Cluster1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Certificate1, *cert)
}

func TestGetConfig(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetConfigTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, string(Config1Response))
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")

	cfg, err := clusters.GetConfig(client, Cluster1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Config1, *cfg)
}

func TestInstances(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstancesTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
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

	err := clusters.ListInstances(client, Cluster1.Name).EachPage(func(page pagination.Page) (bool, error) {
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

	th.Mux.HandleFunc(prepareListInstancesTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
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
	actual, err := clusters.ListInstancesAll(client, Cluster1.Name)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Instance1, ct)
	require.Equal(t, ExpectedInstancesSlice, actual)
}

func TestUpgrade(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUpgradeTestURL(Cluster1.Name), func(w http.ResponseWriter, r *http.Request) {
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
		Version: "v1.27.4",
	}

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	tasks, err := clusters.Upgrade(client, Cluster1.Name, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareVersionTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, VersionResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	count := 0

	err := clusters.Versions(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractVersions(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Version1, ct)
		require.Equal(t, ExpectedVersionSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestVersionsAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareVersionTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, VersionResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("k8s/clusters", "v2")
	actual, err := clusters.VersionsAll(client)
	require.NoError(t, err)
	require.Equal(t, ExpectedVersionSlice, actual)
}
