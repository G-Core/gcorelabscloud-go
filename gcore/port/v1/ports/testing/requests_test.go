package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"

	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareEnableTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "enable_port_security")
}

func prepareDisableTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "disable_port_security")
}

func prepareAllowedAddressPairsTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "allow_address_pairs")
}

func TestEnablePortSecurity(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareEnableTestURL(instanceInterface.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, EnableResponse)
		if err != nil {
			log.Error(err)
		}
	})

	instanceInterface.PortSecurityEnabled = true
	client := fake.ServiceTokenClient("ports", "v1")
	iface, err := ports.EnablePortSecurity(client, instanceInterface.PortID).Extract()
	require.NoError(t, err)
	require.Equal(t, instanceInterface, *iface)
}

func TestDisablePortSecurity(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareDisableTestURL(instanceInterface.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DisableResponse)
		if err != nil {
			log.Error(err)
		}
	})

	instanceInterface.PortSecurityEnabled = false
	client := fake.ServiceTokenClient("ports", "v1")
	iface, err := ports.DisablePortSecurity(client, instanceInterface.PortID).Extract()
	require.NoError(t, err)
	require.Equal(t, instanceInterface, *iface)
}

func TestAllowAddressPairs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAllowedAddressPairsTestURL(PortID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, allowedAddressPairsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, allowedAddressPairsResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ports", "v1")

	opts := ports.AllowAddressPairsOpts{
		AllowedAddressPairs: []reservedfixedips.AllowedAddressPairs{{
			IPAddress:  PortIPRaw1,
			MacAddress: "00:16:3e:f2:87:16",
		}},
	}
	result, err := ports.AllowAddressPairs(client, PortID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, addrPairs1, *result)
}
