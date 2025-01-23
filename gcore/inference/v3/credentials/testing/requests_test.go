package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/credentials"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams(projectID int) string {
	return fmt.Sprintf("/v3/inference/%d/registry_credentials", projectID)
}

func prepareGetTestURLParams(projectID int, name string) string {
	return fmt.Sprintf("/v3/inference/%d/registry_credentials/%s", projectID, name)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID)
}

func prepareGetTestURL(name string) string {
	return prepareGetTestURLParams(fake.ProjectID, name)
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

	client := fake.ServiceTokenClient("inferences", "v3")
	crs, err := credentials.ListAll(client)
	th.AssertNoErr(t, err)

	if len(crs) != 1 {
		t.Errorf("Expected 1 page, got %d", len(crs))
	}

	ct := crs[0]
	require.Equal(t, Creds1, ct)
	require.Equal(t, CredsSlice, crs)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Creds1.Name)
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

	client := fake.ServiceTokenClient("inferences", "v3")

	ct, err := credentials.Get(client, Creds1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Creds1, *ct)
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
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := credentials.CreateRegistryCredentialOpts{
		Name:        Creds1.Name,
		Username:    Creds1.Username,
		Password:    Creds1.Password,
		RegistryURL: Creds1.RegistryURL,
	}

	client := fake.ServiceTokenClient("inferences", "v3")
	cr, err := credentials.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Creds1, *cr)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Creds1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("inferences", "v3")
	err := credentials.Delete(client, Creds1.Name).ExtractErr()
	require.NoError(t, err)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Creds1.Name)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("inferences", "v3")
	opts := credentials.UpdateRegistryCredentialOpts{
		Username:    Creds1.Username,
		Password:    Creds1.Password,
		RegistryURL: Creds1.RegistryURL,
	}

	cr, err := credentials.Update(client, Creds1.Name, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Creds1, *cr)

}
