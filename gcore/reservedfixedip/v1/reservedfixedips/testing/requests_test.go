package testing

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d/%s", projectID, regionID, id)
}

func prepareSwitchVIPTestURL(id string) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d/%s", fake.ProjectID, fake.RegionID, id)
}

func prepareConnectedDeviceTestURL(id string) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d/%s/connected_devices", fake.ProjectID, fake.RegionID, id)
}

func prepareAvailableDeviceTestURL(id string) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d/%s/available_devices", fake.ProjectID, fake.RegionID, id)
}

func preparePortsToShareVIPTestURL(id string) string {
	return fmt.Sprintf("/v1/reserved_fixed_ips/%d/%d/%s/connected_devices", fake.ProjectID, fake.RegionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareDeleteTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
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

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	var count int
	opts := reservedfixedips.ListOpts{}
	err := reservedfixedips.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := reservedfixedips.ExtractReservedFixedIPs(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, ReservedFixedIP1, ct)
		require.Equal(t, ExpectedIPsSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

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

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	opts := reservedfixedips.ListOpts{}
	ips, err := reservedfixedips.ListAll(client, opts)
	require.NoError(t, err)

	require.Equal(t, ReservedFixedIP1, ips[0])
	require.Equal(t, ExpectedIPsSlice, ips)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(ReservedFixedIP1.PortID)

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

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	ct, err := reservedfixedips.Get(client, ReservedFixedIP1.PortID).Extract()

	require.NoError(t, err)
	require.Equal(t, ReservedFixedIP1, *ct)
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

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := reservedfixedips.CreateOpts{
		Type:      reservedfixedips.IPAddress,
		NetworkID: networkTaskID,
		IPAddress: net.ParseIP("192.168.1.2"),
		IsVip:     true,
	}

	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")
	tasks, err := reservedfixedips.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareDeleteTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	task, err := reservedfixedips.Delete(client, ReservedFixedIP1.PortID).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *task)
}

func TestListConnectedDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareConnectedDeviceTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	var count int
	err := reservedfixedips.ListConnectedDevice(client, ReservedFixedIP1.PortID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := reservedfixedips.ExtractDevices(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Device1, ct)
		require.Equal(t, ExpectedDevicesSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAllConnectedDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareConnectedDeviceTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	devices, err := reservedfixedips.ListAllConnectedDevice(client, ReservedFixedIP1.PortID)
	require.NoError(t, err)
	require.Equal(t, Device1, devices[0])
	require.Equal(t, ExpectedDevicesSlice, devices)
}

func TestListAvailableDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareAvailableDeviceTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	var count int
	err := reservedfixedips.ListAvailableDevice(client, ReservedFixedIP1.PortID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := reservedfixedips.ExtractDevices(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Device1, ct)
		require.Equal(t, ExpectedDevicesSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAllAvailableDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareAvailableDeviceTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")

	devices, err := reservedfixedips.ListAllAvailableDevice(client, ReservedFixedIP1.PortID)
	require.NoError(t, err)
	require.Equal(t, Device1, devices[0])
	require.Equal(t, ExpectedDevicesSlice, devices)
}

func TestSwitchVIP(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareSwitchVIPTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SwitchVIPRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := reservedfixedips.SwitchVIPOpts{
		IsVip: true,
	}

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")
	ct, err := reservedfixedips.SwitchVIP(client, ReservedFixedIP1.PortID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, ReservedFixedIP1, *ct)
}

func TestAddPortsToShareVIP(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := preparePortsToShareVIPTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, PortsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := reservedfixedips.PortsToShareVIPOpts{
		PortIDs: []string{
			"351b0dd7-ca09-431c-be53-935db3785067",
			"bc688791-f1b0-44eb-97d4-07697294b1e1",
		},
	}

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")
	actual, err := reservedfixedips.AddPortsToShareVIP(client, ReservedFixedIP1.PortID, options).Extract()
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Device1, ct)
	require.Equal(t, ExpectedDevicesSlice, actual)

	th.AssertNoErr(t, err)
}

func TestReplacePortsToShareVIP(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := preparePortsToShareVIPTestURL(ReservedFixedIP1.PortID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, PortsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeviceResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := reservedfixedips.PortsToShareVIPOpts{
		PortIDs: []string{
			"351b0dd7-ca09-431c-be53-935db3785067",
			"bc688791-f1b0-44eb-97d4-07697294b1e1",
		},
	}

	client := fake.ServiceTokenClient("reserved_fixed_ips", "v1")
	actual, err := reservedfixedips.ReplacePortsToShareVIP(client, ReservedFixedIP1.PortID, options).Extract()
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Device1, ct)
	require.Equal(t, ExpectedDevicesSlice, actual)

	th.AssertNoErr(t, err)
}
