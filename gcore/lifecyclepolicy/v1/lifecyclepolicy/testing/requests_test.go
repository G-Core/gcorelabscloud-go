package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func getURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v", fake.ProjectID, fake.RegionID, id)
}

func listURL() string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d", fake.ProjectID, fake.RegionID)
}

func deleteURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v", fake.ProjectID, fake.RegionID, id)
}

func createURL() string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d", fake.ProjectID, fake.RegionID)
}

func updateURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v", fake.ProjectID, fake.RegionID, id)
}

func addVolumesURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v/add_volumes_to_policy", fake.ProjectID, fake.RegionID, id)
}

func removeVolumesURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v/remove_volumes_from_policy", fake.ProjectID, fake.RegionID, id)
}

func addSchedulesURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v/add_schedules", fake.ProjectID, fake.RegionID, id)
}

func removeSchedulesURL(id int) string {
	return fmt.Sprintf("/v1/lifecycle_policy/%d/%d/%v/remove_schedules", fake.ProjectID, fake.RegionID, id)
}

func TestGet(t *testing.T) {
	for i, policy := range policies {
		th.SetupHTTP()
		client := fake.ServiceTokenClient("lifecycle_policy", "v1")

		th.Mux.HandleFunc(getURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodGet)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			response := responsesWithoutVolumes[i]
			query := r.URL.Query()
			if query.Get("need_volumes") == "true" {
				response = responsesWithVolumes[i]
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, response)
			if err != nil {
				log.Error(err)
			}
		})

		result := lifecyclepolicy.Get(client, policy.ID, lifecyclepolicy.GetOpts{NeedVolumes: true})
		ct, err := result.Extract()
		require.NoError(t, err)
		require.Equal(t, policy, *ct)

		policy.Volumes = nil
		result = lifecyclepolicy.Get(client, policy.ID, lifecyclepolicy.GetOpts{NeedVolumes: false})
		ct, err = result.Extract()
		require.NoError(t, err)
		require.Equal(t, policy, *ct)

		th.TeardownHTTP()
	}
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("lifecycle_policy", "v1")

	th.Mux.HandleFunc(listURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, listPoliciesResponseWithVolumes)
		if err != nil {
			log.Error(err)
		}
	})

	result := lifecyclepolicy.ListAll(client, lifecyclepolicy.ListOpts{NeedVolumes: true})
	ct, err := result.Extract()
	require.NoError(t, err)
	require.Equal(t, policies, ct)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("lifecycle_policy", "v1")

	for _, policy := range policies {
		th.Mux.HandleFunc(deleteURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodDelete)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			w.WriteHeader(http.StatusNoContent)
		})

		err := lifecyclepolicy.Delete(client, policy.ID)
		require.NoError(t, err)
	}
}

func TestCreate(t *testing.T) {
	for i, policy := range policies {
		th.SetupHTTP()
		client := fake.ServiceTokenClient("lifecycle_policy", "v1")

		th.Mux.HandleFunc(createURL(), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodPost)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, createRequests[i])

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, responsesWithoutVolumes[i])
			if err != nil {
				log.Error(err)
			}
		})

		result := lifecyclepolicy.Create(client, createOpts[i])
		ct, err := result.Extract()
		require.NoError(t, err)
		policy.Volumes = nil
		require.Equal(t, policy, *ct)

		th.TeardownHTTP()
	}
}

func TestUpdate(t *testing.T) {
	for i, policy := range policies {
		th.SetupHTTP()
		client := fake.ServiceTokenClient("lifecycle_policy", "v1")

		th.Mux.HandleFunc(updateURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodPatch)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, updateRequests[i])

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, responsesWithoutVolumes[i])
			if err != nil {
				log.Error(err)
			}
		})

		result := lifecyclepolicy.Update(client, policy.ID, updateOpts[i])
		ct, err := result.Extract()
		require.NoError(t, err)
		policy.Volumes = nil
		require.Equal(t, policy, *ct)

		th.TeardownHTTP()
	}
}

func TestAddVolumes(t *testing.T) {
	for i, policy := range policies {
		th.SetupHTTP()
		client := fake.ServiceTokenClient("lifecycle_policy", "v1")

		th.Mux.HandleFunc(addVolumesURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodPut)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, addVolumesRequests[i])

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, responsesWithoutVolumes[i])
			if err != nil {
				log.Error(err)
			}
		})

		result := lifecyclepolicy.AddVolumes(client, policy.ID, addVolumesOpts[i])
		ct, err := result.Extract()
		require.NoError(t, err)
		policy.Volumes = nil
		require.Equal(t, policy, *ct)

		th.TeardownHTTP()
	}
}

func TestRemoveVolumes(t *testing.T) {
	for i, policy := range policies {
		th.SetupHTTP()
		client := fake.ServiceTokenClient("lifecycle_policy", "v1")

		th.Mux.HandleFunc(removeVolumesURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, http.MethodPut)
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, removeVolumesRequests[i])

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, responsesWithoutVolumes[i])
			if err != nil {
				log.Error(err)
			}
		})

		result := lifecyclepolicy.RemoveVolumes(client, policy.ID, removeVolumesOpts[i])
		ct, err := result.Extract()
		require.NoError(t, err)
		policy.Volumes = nil
		require.Equal(t, policy, *ct)

		th.TeardownHTTP()
	}
}

func TestAddSchedules(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("lifecycle_policy", "v1")

	policy, response := policies[0], responsesWithoutVolumes[0]

	th.Mux.HandleFunc(addSchedulesURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, addSchedulesRequests)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, response)
		if err != nil {
			log.Error(err)
		}
	})

	result := lifecyclepolicy.AddSchedules(client, policy.ID, addSchedulesOpts)
	ct, err := result.Extract()
	require.NoError(t, err)
	policy.Volumes = nil
	require.Equal(t, policy, *ct)
}

func TestRemoveSchedules(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fake.ServiceTokenClient("lifecycle_policy", "v1")

	policy, response := policies[1], responsesWithoutVolumes[1]

	th.Mux.HandleFunc(removeSchedulesURL(policy.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, removeSchedulesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, response)
		if err != nil {
			log.Error(err)
		}
	})

	result := lifecyclepolicy.RemoveSchedules(client, policy.ID, removeSchedulesOpts)
	ct, err := result.Extract()
	require.NoError(t, err)
	policy.Volumes = nil
	require.Equal(t, policy, *ct)
}
