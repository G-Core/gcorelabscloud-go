package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/apitokens"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/types"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func prepareRootTestURL(clientID int) string {
	return fmt.Sprintf("/clients/%d/tokens", clientID)
}

func prepareResourceTestURL(clientID, tokenID int) string {
	return fmt.Sprintf("/clients/%d/tokens/%d", clientID, tokenID)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRootTestURL(clientID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceClient()

	results, err := apitokens.List(client, clientID, nil).Extract()
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, apiToken1, ct)
	require.Equal(t, ExpectedAPITokenSlice, results)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareResourceTestURL(clientID, apiToken1.ID)

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

	client := fake.ServiceClient()
	ct, err := apitokens.Get(client, clientID, apiToken1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, apiToken1, *ct)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResourceTestURL(clientID, apiToken1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceClient()
	err := apitokens.Delete(client, clientID, apiToken1.ID).ExtractErr()
	require.NoError(t, err)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRootTestURL(clientID), func(w http.ResponseWriter, r *http.Request) {
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

	options := apitokens.CreateOpts{
		Name:        "My token",
		Description: "It's my token",
		ClientUser: apitokens.CreateClientUser{
			Role: apitokens.ClientRole{
				ID:   types.RoleIDAdministrators,
				Name: types.RoleNameAdministrators,
			},
		},
	}

	client := fake.ServiceClient()
	ct, err := apitokens.Create(client, clientID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Token1, *ct)

}
