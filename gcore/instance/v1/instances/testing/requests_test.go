package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"

	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/instances/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, id)
}

func prepareGetActionTestURLParams(id string, action string) string {
	return fmt.Sprintf("/v1/instances/%d/%d/%s/%s", fake.ProjectID, fake.RegionID, id, action)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareListInterfacesTestURL(id string) string {
	return prepareGetActionTestURLParams(id, "interfaces")
}

func prepareListSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams(id, "securitygroups")
}

func prepareAssignSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams(id, "addsecuritygroup")
}

func prepareUnAssignSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams(id, "delsecuritygroup")
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

	client := fake.ServiceTokenClient("instances", "v1")
	count := 0

	opts := instances.ListOpts{}

	err := instances.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
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

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Instance1.ID)

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

	client := fake.ServiceTokenClient("instances", "v1")

	ct, err := instances.Get(client, Instance1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Instance1, *ct)

}

func TestListAllInterfaces(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInterfacesTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, InterfacesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	interfaces, err := instances.ListInterfacesAll(client, instanceID)

	require.NoError(t, err)
	require.Len(t, interfaces, 1)
	require.Equal(t, PortID, interfaces[0].PortID)
	require.Equal(t, ExpectedInstanceInterfacesSlice, interfaces)
}

func TestListAllSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, SecurityGroupsListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	sgs, err := instances.ListSecurityGroupsAll(client, instanceID)

	require.NoError(t, err)
	require.Len(t, sgs, 1)
	require.Equal(t, SecurityGroup1, sgs[0])
	require.Equal(t, ExpectedSecurityGroupsSlice, sgs)
}

func TestUnAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUnAssignSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, UnAssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := instances.UnAssignSecurityGroup(client, instanceID, opts).ExtractErr()

	require.NoError(t, err)
}

func TestAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAssignSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, AssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := instances.AssignSecurityGroup(client, instanceID, opts).ExtractErr()

	require.NoError(t, err)
}
