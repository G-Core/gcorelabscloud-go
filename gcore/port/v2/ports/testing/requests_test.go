package testing

import (
	"fmt"
	"net/http"
	"testing"

	ports1 "github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	ports2 "github.com/G-Core/gcorelabscloud-go/gcore/port/v2/ports"

	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string {
	return fmt.Sprintf("/v2/ports/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareAllowedAddressPairsTestURLV2(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "allow_address_pairs")
}

func TestAllowAddressPairsV2(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAllowedAddressPairsTestURLV2(PortID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, allowedAddressPairsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, allowedAddressPairsResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ports", "v2")

	opts := ports1.AllowAddressPairsOpts{
		AllowedAddressPairs: []reservedfixedips.AllowedAddressPairs{{
			IPAddress:  PortIPRaw1,
			MacAddress: "00:16:3e:f2:87:16",
		}},
	}
	tasks, err := ports2.AllowAddressPairs(client, PortID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}
