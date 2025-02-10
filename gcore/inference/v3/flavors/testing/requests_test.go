package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/flavors"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams() string {
	return fmt.Sprintf("/v3/inference/flavors")
}

func prepareGetTestURLParams(name string) string {
	return fmt.Sprintf("/v3/inference/flavors/%s", name)
}

func prepareListTestURL() string {
	return prepareListTestURLParams()
}

func prepareGetTestURL(name string) string {
	return prepareGetTestURLParams(name)
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
	fls, err := flavors.ListAllFlavor(client)
	th.AssertNoErr(t, err)

	if len(fls) != 1 {
		t.Errorf("Expected 1 page, got %d", len(fls))
	}

	ct := fls[0]
	require.Equal(t, Flavor1, ct)
	require.Equal(t, FlavorSlice, fls)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Flavor1.Name)
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

	ct, err := flavors.GetFlavor(client, Flavor1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Flavor1, *ct)
}
