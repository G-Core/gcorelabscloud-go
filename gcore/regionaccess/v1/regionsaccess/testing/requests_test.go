package testing

import (
	"fmt"
	"net/http"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/regionaccess/v1/regionsaccess"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareRootTestURL() string {
	return "/v1/reseller_region"
}

func prepareDeleteTestURL(id int) string {
	return fmt.Sprintf("/v1/reseller_region/%d", id)
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRootTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reseller_region", "v1")

	results, err := regionsaccess.ListAll(client, nil)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, RegionAccess1, ct)
	require.Equal(t, ExpectedRegionAccessSlice, results)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRootTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	opts := regionsaccess.CreateOpts{
		AccessAllEdgeRegions: true,
		RegionIDs:            []int{1, 2, 3},
		ClientID:             &clientID,
		ResellerID:           &resellerID,
	}

	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)

	client := fake.ServiceTokenClient("reseller_region", "v1")
	region, err := regionsaccess.Create(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, RegionAccessCreated1, *region)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeleteTestURL(resellerID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("reseller_region", "v1")
	err := regionsaccess.Delete(client, resellerID).ExtractErr()
	require.NoError(t, err)

}
