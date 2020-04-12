package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/limits"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListTestURL() string {
	return "/v1/limits_request"
}

func prepareItemTestURL() string {
	return fmt.Sprintf("/v1/limits_request/%d", limitRequestID)
}

const limitRequestID = 1

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

	client := fake.ServiceTokenClient("limits_request", "v1")
	count := 0

	err := limits.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := limits.ExtractLimitResults(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, LimitRequest1, ct)
		require.Equal(t, ExpectedLimitRequestSlice, actual)
		return true, nil
	})

	require.NoError(t, err)

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

	client := fake.ServiceTokenClient("limits_request", "v1")

	actual, err := limits.ListAll(client)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, LimitRequest1, ct)
	require.Equal(t, ExpectedLimitRequestSlice, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL()

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

	client := fake.ServiceTokenClient("limits_request", "v1")

	limit, err := limits.Get(client, limitRequestID).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)

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

	options := limits.NewCreateOpts("test")
	options.RequestedQuotas.ExternalIPCountLimit = 4
	client := fake.ServiceTokenClient("limits_request", "v1")
	limit, err := limits.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL()

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

	client := fake.ServiceTokenClient("limits_request", "v1")

	options := limits.NewUpdateOpts()
	options.ExternalIPCountLimit = 4

	limit, err := limits.Update(client, limitRequestID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)

}

func TestStatus(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, StatusRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, StatusResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("limits_request", "v1")

	options := limits.StatusOpts{
		Status: types.LimitRequestDone,
	}

	limit, err := limits.Status(client, limitRequestID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)

}

func TestDelete(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, StatusResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("limits_request", "v1")
	err := limits.Delete(client, limitRequestID).ExtractErr()
	require.NoError(t, err)

}
