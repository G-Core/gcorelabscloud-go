package testing

import (
	"fmt"
	metadataV1Testing "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata/testing"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, id)
}

func prepareCustomSecurityGroupTestURL(id string) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/securitygroup", fake.ProjectID, fake.RegionID, id)
}

func prepareResizeLoadBalancerURL(id string) string {
	return fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/resize", fake.ProjectID, fake.RegionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}
func prepareGetActionTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/floatingips/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}
func prepareMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/floatingips/%d/%d/%s/%s", "v1", fake.ProjectID, fake.RegionID, id, "metadata_item")
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

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	count := 0

	err := loadbalancers.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := loadbalancers.ExtractLoadBalancers(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, LoadBalancer1, ct)
		require.Equal(t, ExpectedLoadBalancerSlice, actual)
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

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	lbs, err := loadbalancers.ListAll(client, nil)
	require.NoError(t, err)
	lb := lbs[0]
	require.Equal(t, LoadBalancer1, lb)
	require.Equal(t, ExpectedLoadBalancerSlice, lbs)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(LoadBalancer1.ID)

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

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	ct, err := loadbalancers.Get(client, LoadBalancer1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, LoadBalancer1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)

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

	urlPath := "/"
	maxRetriesDown := 3
	memberWeight1 := 2
	memberWeight2 := 4
	instanceID1 := "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
	instanceID2 := "169942e0-9b53-42df-95ef-1a8b6525c2bd"

	options := loadbalancers.CreateOpts{
		Name: LoadBalancer1.Name,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:             "listener_name",
			ProtocolPort:     80,
			Protocol:         types.ProtocolTypeHTTP,
			Certificate:      "",
			CertificateChain: "",
			PrivateKey:       "",
			Pools: []loadbalancers.CreatePoolOpts{{
				Name:     "pool_name",
				Protocol: types.ProtocolTypeHTTP,
				Members: []loadbalancers.CreatePoolMemberOpts{{
					InstanceID:   instanceID1,
					Address:      net.ParseIP("192.168.1.101"),
					ProtocolPort: 8000,
					Weight:       memberWeight1,
					SubnetID:     "",
				}, {
					Address:      net.ParseIP("192.168.1.102"),
					ProtocolPort: 8000,
					Weight:       memberWeight2,
					SubnetID:     "",
					InstanceID:   instanceID2,
				},
				},
				HealthMonitor: &loadbalancers.CreateHealthMonitorOpts{
					Type:           types.HealthMonitorTypeHTTP,
					Delay:          10,
					MaxRetries:     3,
					Timeout:        5,
					MaxRetriesDown: maxRetriesDown,
					HTTPMethod:     types.HTTPMethodPointer(types.HTTPMethodGET),
					URLPath:        urlPath,
				},
				LoadBalancerAlgorithm: types.LoadBalancerAlgorithmRoundRobin,
				SessionPersistence:    nil,
			}},
		}},
		VipPortID:   "169942e0-9b53-42df-95ef-1a8b6525c2bd",
		VIPIPFamily: types.DualStackIPFamilyType,
	}

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	tasks, err := loadbalancers.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(LoadBalancer1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")
	tasks, err := loadbalancers.Delete(client, LoadBalancer1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(LoadBalancer1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	opts := loadbalancers.UpdateOpts{
		Name: LoadBalancer1.Name,
	}

	ct, err := loadbalancers.Update(client, LoadBalancer1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, LoadBalancer1, *ct)
	require.Equal(t, LoadBalancer1.Name, ct.Name)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)

}

func TestGetCustomSecurityGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCustomSecurityGroupTestURL(LoadBalancer1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListCustomSecurityGroupResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	lbSg, err := loadbalancers.ListCustomSecurityGroup(client, LoadBalancer1.ID).Extract()
	require.NoError(t, err)
	lb := lbSg[0]
	require.Equal(t, LbSecurityGroup1, lb)
	require.Equal(t, ExpectedLbSecurityGroupSlice, lbSg)
}

func TestCreateCustomSecurityGroup(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCustomSecurityGroupTestURL(LoadBalancer1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	err := loadbalancers.CreateCustomSecurityGroup(client, LoadBalancer1.ID).ExtractErr()
	require.NoError(t, err)
}

func TestResizeLoadBalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResizeLoadBalancerURL(LoadBalancer1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResizeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ResizeResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("loadbalancers", "v1")

	opts := loadbalancers.ResizeOpts{
		Flavor: Flavor,
	}

	tasks, err := loadbalancers.Resize(client, LoadBalancer1.ID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestMetadataListAll(t *testing.T) {
	metadataV1Testing.BuildTestMetadataListAll("loadbalancers", LoadBalancer1.ID)(t)
}

func TestMetadataGet(t *testing.T) {
	metadataV1Testing.BuildTestMetadataGet("loadbalancers", LoadBalancer1.ID)(t)
}

func TestMetadataCreate(t *testing.T) {
	metadataV1Testing.BuildTestMetadataCreate("loadbalancers", LoadBalancer1.ID)(t)
}

func TestMetadataUpdate(t *testing.T) {
	metadataV1Testing.BuildTestMetadataUpdate("loadbalancers", LoadBalancer1.ID)(t)

}

func TestMetadataDelete(t *testing.T) {
	metadataV1Testing.BuildTestMetadataDelete("loadbalancers", LoadBalancer1.ID)(t)
}
