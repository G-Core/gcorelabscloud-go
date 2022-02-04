package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v2/limits"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func prepareListTestURL() string {
	return "/v2/limits_request"
}

func prepareItemTestURL(limitID int) string {
	return fmt.Sprintf("/v2/limits_request/%d", limitID)
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

	client := fake.ServiceTokenClient("limits_request", "v2")
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

	client := fake.ServiceTokenClient("limits_request", "v2")

	actual, err := limits.ListAll(client)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, LimitRequest1, ct)
	require.Equal(t, ExpectedLimitRequestSlice, actual)
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
	options.RequestedQuotas.RegionalLimits = []limits.RegionalLimits{
		{RegionID: 1, FirewallCountLimit: 13, CPUCountLimit: 1},
	}
	client := fake.ServiceTokenClient("limits_request", "v2")
	limit, err := limits.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL(limitRequestID)

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

	client := fake.ServiceTokenClient("limits_request", "v2")

	limit, err := limits.Get(client, limitRequestID).Extract()
	require.NoError(t, err)
	require.Equal(t, LimitRequest1, *limit)
}

func TestDelete(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL(limitRequestID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	client := fake.ServiceTokenClient("limits_request", "v2")
	err := limits.Delete(client, limitRequestID).ExtractErr()
	require.NoError(t, err)
}
