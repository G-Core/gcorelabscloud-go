package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams() string {
	return "/v2/keypairs"
}

func prepareGetTestURLParams(id string) string {
	return fmt.Sprintf("/v2/keypairs/%s", id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams()
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(id)
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

	client := fake.ServiceTokenClient("keypairs", "v2")
	count := 0

	err := keypairs.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := keypairs.ExtractKeyPairs(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, KeyPair1, ct)
		require.Equal(t, ExpectedKeyPairsSlice, actual)
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

	testURL := prepareGetTestURL(KeyPair1.ID)

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

	client := fake.ServiceTokenClient("keypairs", "v2")

	ct, err := keypairs.Get(client, KeyPair1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, KeyPair1, *ct)

}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	url := prepareListTestURL()

	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
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

	options := keypairs.CreateOpts{
		Name:      KeyPair1.Name,
		PublicKey: KeyPair1.PublicKey,
		ProjectID: KeyPair1.ProjectID,
	}

	client := fake.ServiceTokenClient("keypairs", "v2")
	tasks, err := keypairs.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, KeyPair1, *tasks)

}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(KeyPair1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
	})

	client := fake.ServiceTokenClient("keypairs", "v2")
	err := keypairs.Delete(client, KeyPair1.ID).ExtractErr()
	require.NoError(t, err)
}
