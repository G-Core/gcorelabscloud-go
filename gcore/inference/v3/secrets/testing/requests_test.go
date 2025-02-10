package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/secrets"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams(projectID int) string {
	return fmt.Sprintf("/v3/inference/%d/secrets", projectID)
}

func prepareGetTestURLParams(projectID int, name string) string {
	return fmt.Sprintf("/v3/inference/%d/secrets/%s", projectID, name)
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
	crs, err := secrets.ListAll(client)
	th.AssertNoErr(t, err)

	if len(crs) != 1 {
		t.Errorf("Expected 1 page, got %d", len(crs))
	}

	ct := crs[0]
	require.Equal(t, Secret1, ct)
	require.Equal(t, SecretSlice, crs)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Secret1.Name)
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

	ct, err := secrets.Get(client, Secret1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Secret1, *ct)
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

	options := secrets.CreateInferenceSecretOpts{
		Name: Secret1.Name,
		Type: Secret1.Type,
		Data: secrets.CreateSecretData{
			AWSSecretKeyID:     Secret1.Data.AWSSecretKeyID,
			AWSSecretAccessKey: Secret1.Data.AWSSecretAccessKey,
		},
	}

	client := fake.ServiceTokenClient("inferences", "v3")
	cr, err := secrets.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Secret1, *cr)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Secret1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("inferences", "v3")
	err := secrets.Delete(client, Secret1.Name).ExtractErr()
	require.NoError(t, err)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Secret1.Name)

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
	opts := secrets.UpdateInferenceSecretOpts{
		Type: Secret1.Type,
		Data: secrets.CreateSecretData{
			AWSSecretKeyID:     Secret1.Data.AWSSecretKeyID,
			AWSSecretAccessKey: Secret1.Data.AWSSecretAccessKey,
		},
	}

	cr, err := secrets.Update(client, Secret1.Name, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Secret1, *cr)

}
