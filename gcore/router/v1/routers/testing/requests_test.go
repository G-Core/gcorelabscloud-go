package testing

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"net"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/routers/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/routers/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
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

	client := fake.ServiceTokenClient("routers", "v1")
	count := 0

	opts := routers.ListOpts{
		ID: Router1.ID,
	}

	err := routers.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := routers.ExtractRouters(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Router1, ct)
		require.Equal(t, ExpectedRouterSlice, actual)
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

	client := fake.ServiceTokenClient("routers", "v1")

	opts := routers.ListOpts{
		ID: Router1.ID,
	}

	results, err := routers.ListAll(client, opts)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, Router1, ct)
	require.Equal(t, ExpectedRouterSlice, results)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Router1.ID)

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

	client := fake.ServiceTokenClient("routers", "v1")
	ct, err := routers.Get(client, Router1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Router1, *ct)
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

	var gccidr gcorecloud.CIDR
	_, netIPNet, _ := net.ParseCIDR("10.0.3.0/24")
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask

	snat := true

	options := routers.CreateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		},
		Interfaces: []routers.Interface{
			{
				Type:     types.SubnetInterfaceType,
				SubnetID: "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
			},
		},
		Routes: []subnets.HostRoute{
			{
				Destination: gccidr,
				NextHop:     net.ParseIP("10.0.0.13"),
			},
		},
	}

	client := fake.ServiceTokenClient("routers", "v1")
	tasks, err := routers.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Router1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("routers", "v1")
	tasks, err := routers.Delete(client, Router1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Router2.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponseUpdate)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("routers", "v1")

	var gccidr gcorecloud.CIDR
	_, netIPNet, _ := net.ParseCIDR("10.0.4.0/24")
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask

	snat := false

	opts := routers.UpdateOpts{
		Name: Router2.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f2",
		},
		Routes: []subnets.HostRoute{
			{
				Destination: gccidr,
				NextHop:     net.ParseIP("10.0.0.14"),
			},
		},
	}

	ct, err := routers.Update(client, Router2.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Router2, *ct)
	require.Equal(t, Router2.Name, ct.Name)
	require.Equal(t, Router2.ExternalGatewayInfo.EnableSNat, ct.ExternalGatewayInfo.EnableSNat)
	require.Equal(t, Router2.ExternalGatewayInfo.NetworkID, ct.ExternalGatewayInfo.NetworkID)
}
