package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/quota/v1/quotas"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListTestURL() string {
	return "/v1/client_quotas"
}

func prepareGetTestURL(id int) string {
	return fmt.Sprintf("/v1/client_quotas/%d", id)
}

const clientID = 1

func TestOwnQuota(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("client_quotas", "v1")

	quota, err := quotas.OwnQuota(client).Extract()
	require.NoError(t, err)
	require.Equal(t, Quota1, *quota)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(clientID)

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

	client := fake.ServiceTokenClient("client_quotas", "v1")

	quota, err := quotas.Get(client, clientID).Extract()
	require.NoError(t, err)
	require.Equal(t, Quota1, *quota)

}

func TestReplace(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(clientID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ReplaceRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, ReplaceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := quotas.ReplaceOpts{Quota: Quota1}

	client := fake.ServiceTokenClient("client_quotas", "v1")
	quota, err := quotas.Replace(client, clientID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Quota1, *quota)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(clientID)

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

	client := fake.ServiceTokenClient("client_quotas", "v1")

	options := quotas.UpdateOpts{Quota: quotas.NewQuota()}
	options.ExternalIPCountUsage = 4

	quota, err := quotas.Update(client, clientID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Quota1, *quota)

}
