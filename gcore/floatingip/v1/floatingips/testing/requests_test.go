package testing

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/floatingip/v1/floatingips"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/floatingips/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, id)
}

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareAssignTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "assign")
}

func prepareUnAssignTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "unassign")
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

	client := fake.ServiceTokenClient("floatingips", "v1")
	count := 0

	err := floatingips.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := floatingips.ExtractFloatingIPs(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, floatingIPDetails, ct)
		require.Equal(t, ExpectedFloatingIPSlice, actual)
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

	client := fake.ServiceTokenClient("floatingips", "v1")

	groups, err := floatingips.ListAll(client)
	require.NoError(t, err)
	ct := groups[0]
	require.Equal(t, floatingIPDetails, ct)
	require.Equal(t, ExpectedFloatingIPSlice, groups)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(floatingIPDetails.ID)

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

	client := fake.ServiceTokenClient("floatingips", "v1")

	ct, err := floatingips.Get(client, floatingIP.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, floatingIP, *ct)
	require.Equal(t, floatingIPCreatedAt, ct.CreatedAt)
	require.Equal(t, floatingIPUpdatedAt, *ct.UpdatedAt)

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

	options := floatingips.CreateOpts{
		PortID:         "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		FixedIPAddress: net.ParseIP("192.168.10.15"),
	}

	client := fake.ServiceTokenClient("floatingips", "v1")
	tasks, err := floatingips.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(floatingIP.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("floatingips", "v1")
	tasks, err := floatingips.Delete(client, floatingIP.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestAssign(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareAssignTestURL(floatingIP.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AssignRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AssignResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := floatingips.CreateOpts{
		PortID:         "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		FixedIPAddress: net.ParseIP("192.168.10.15"),
	}

	client := fake.ServiceTokenClient("floatingips", "v1")
	ip, err := floatingips.Assign(client, floatingIP.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, floatingIP, *ip)
	require.Equal(t, floatingIPCreatedAt, ip.CreatedAt)
	require.Equal(t, floatingIPUpdatedAt, *ip.UpdatedAt)

}

func TestUnAssign(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareUnAssignTestURL(floatingIP.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UnassignResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("floatingips", "v1")
	ip, err := floatingips.UnAssign(client, floatingIP.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, floatingIP, *ip)
	require.Equal(t, floatingIPCreatedAt, ip.CreatedAt)
	require.Equal(t, floatingIPUpdatedAt, *ip.UpdatedAt)

}
