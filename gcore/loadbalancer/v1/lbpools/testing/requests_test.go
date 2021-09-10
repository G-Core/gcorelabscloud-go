package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/lbpools/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareActionTestURL(projectID int, regionID int, id string, action string) string {
	return fmt.Sprintf("/v1/lbpools/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareActionDetailTestURL(projectID int, regionID int, id string, action string, actionID string) string {
	return fmt.Sprintf("/v1/lbpools/%d/%d/%s/%s/%s", projectID, regionID, id, action, actionID)
}

func prepareCreateMemberURL(id string) string {
	return prepareActionTestURL(fake.ProjectID, fake.RegionID, id, "member")
}

func prepareHealthMonitorURL(id string) string {
	return prepareActionTestURL(fake.ProjectID, fake.RegionID, id, "healthmonitor")
}

func prepareDeleteMemberURL(id string, memberID string) string {
	return prepareActionDetailTestURL(fake.ProjectID, fake.RegionID, id, "member", memberID)
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

	client := fake.ServiceTokenClient("lbpools", "v1")
	count := 0

	opts := lbpools.ListOpts{LoadBalancerID: &LBPool1.ID}

	err := lbpools.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		pools, err := lbpools.ExtractPools(page)
		require.NoError(t, err)
		pool := pools[0]
		require.Equal(t, LBPool1, pool)
		require.Equal(t, ExpectedLBPoolsSlice, pools)
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

	client := fake.ServiceTokenClient("lbpools", "v1")

	opts := lbpools.ListOpts{LoadBalancerID: &LBPool1.ID}

	pools, err := lbpools.ListAll(client, opts)
	require.NoError(t, err)
	pool := pools[0]
	require.Equal(t, LBPool1, pool)
	require.Equal(t, ExpectedLBPoolsSlice, pools)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(LBPool1.ID)

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

	client := fake.ServiceTokenClient("lbpools", "v1")

	ct, err := lbpools.Get(client, LBPool1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, LBPool1, *ct)

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

	options := lbpools.CreateOpts{
		Name:            LBPool1.Name,
		Protocol:        LBPool1.Protocol,
		LBPoolAlgorithm: LBPool1.LoadBalancerAlgorithm,
		Members: []lbpools.CreatePoolMemberOpts{
			{
				Address:      *Member1.Address,
				ProtocolPort: Member1.ProtocolPort,
				Weight:       Member1.Weight,
				SubnetID:     Member1.SubnetID,
				InstanceID:   Member1.InstanceID,
			},
			{
				Address:      *Member2.Address,
				ProtocolPort: Member2.ProtocolPort,
				Weight:       Member2.Weight,
				SubnetID:     Member2.SubnetID,
				InstanceID:   Member2.InstanceID,
			},
		},
		LoadBalancerID:     LoadBalancerID,
		ListenerID:         ListenerID,
		HealthMonitor:      nil,
		SessionPersistence: nil,
	}

	client := fake.ServiceTokenClient("lbpools", "v1")
	tasks, err := lbpools.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(LBPool1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("lbpools", "v1")
	tasks, err := lbpools.Delete(client, LBPool1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(LBPool1.ID)

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

	client := fake.ServiceTokenClient("lbpools", "v1")

	opts := lbpools.UpdateOpts{
		Name: LBPool1.Name,
	}

	tasks, err := lbpools.Update(client, LBPool1.ID, opts).Extract()

	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestDeleteMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	memberID := LBPool1.Members[0].ID

	th.Mux.HandleFunc(prepareDeleteMemberURL(LBPool1.ID, memberID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteMemberResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("lbpools", "v1")
	tasks, err := lbpools.DeleteMember(client, LBPool1.ID, memberID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestCreateMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateMemberURL(LBPool1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePoolMemberRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := lbpools.CreatePoolMemberOpts{
		Address:      *Member1.Address,
		ProtocolPort: Member1.ProtocolPort,
		Weight:       Member1.Weight,
		SubnetID:     Member1.SubnetID,
		InstanceID:   Member1.InstanceID,
	}

	client := fake.ServiceTokenClient("lbpools", "v1")
	tasks, err := lbpools.CreateMember(client, LBPool1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDeleteHealthMonitor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareHealthMonitorURL(LBPool1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("lbpools", "v1")
	err := lbpools.DeleteHealthMonitor(client, LBPool1.ID).ExtractErr()
	require.NoError(t, err)

}

func TestCreateHealthMonitor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareHealthMonitorURL(LBPool1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateHealthMonitorRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	method := types.HTTPMethodGET
	opts := lbpools.CreateHealthMonitorOpts{
		Type:           types.HealthMonitorTypeHTTP,
		Delay:          5,
		MaxRetries:     1,
		Timeout:        30,
		MaxRetriesDown: 3,
		HTTPMethod:     &method,
		URLPath:        "/",
		ExpectedCodes:  "200,301,302",
	}

	client := fake.ServiceTokenClient("lbpools", "v1")
	tasks, err := lbpools.CreateHealthMonitor(client, LBPool1.ID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}
