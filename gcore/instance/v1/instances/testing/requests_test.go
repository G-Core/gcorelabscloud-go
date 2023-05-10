package testing

import (
	"fmt"
	"net/http"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"

	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/instances/%d/%d", version, projectID, regionID)
}

func prepareGetTestURLParams(version string, projectID int, regionID int, id string) string {
	return fmt.Sprintf("/%s/instances/%d/%d/%s", version, projectID, regionID, id)
}

func prepareListInstanceLocationTestURL() string {
	return "/v1/instances/search"
}

func prepareGetActionTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/instances/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}

func prepareGetActionDetailsTestURLParams(version string, id string, action, actionID string) string { // nolint
	return fmt.Sprintf("/%s/instances/%d/%d/%s/%s/%s", version, fake.ProjectID, fake.RegionID, id, action, actionID)
}

func prepareListTestURL() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
}

func prepareListInterfacesTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "interfaces")
}

func prepareGetSpiceConsoleTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "get_spice_console")
}

func prepareGetInstanceConsoleTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "get_console")
}

func prepareListInstanceMetricsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metrics")
}

func prepareListAvailableFlavorsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "available_flavors")
}

func prepareListSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "securitygroups")
}

func prepareListPortsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "ports")
}

func prepareAssignSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "addsecuritygroup")
}

func prepareUnAssignSecurityGroupsTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "delsecuritygroup")
}

func prepareAttachInterfaceTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "attach_interface")
}

func prepareDetachInterfaceTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "detach_interface")
}

func prepareStartTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "start")
}

func prepareStopTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "stop")
}

func preparePowerCycleTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "powercycle")
}

func prepareRebootTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "reboot")
}

func prepareSuspendTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "suspend")
}

func prepareResumeTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "resume")
}

func prepareChangeFlavorTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "changeflavor")
}

func prepareMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metadata")
}

func prepareMetadataDetailsTestURL(id, key string) string {
	return prepareGetActionDetailsTestURLParams("v1", id, "metadata", key)
}

func prepareRenameInstanceTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareDeleteTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareCreateTestURLV2() string {
	return prepareListTestURLParams("v2", fake.ProjectID, fake.RegionID)
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

	client := fake.ServiceTokenClient("instances", "v1")
	count := 0

	opts := instances.ListOpts{}

	err := instances.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := instances.ExtractInstances(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Instance1, ct)
		require.Equal(t, ExpectedInstancesSlice, actual)
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

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.ListOpts{}

	actual, err := instances.ListAll(client, opts)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Instance1, ct)
	require.Equal(t, ExpectedInstancesSlice, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Instance1.ID)

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

	client := fake.ServiceTokenClient("instances", "v1")

	ct, err := instances.Get(client, Instance1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Instance1, *ct)

}

func TestRenameInstance(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareRenameInstanceTestURL(Instance1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestJSONRequest(t, r, RenameInstanceRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.RenameInstanceOpts{Name: "new-name"}
	ct, err := instances.RenameInstance(client, Instance1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Instance1, *ct)

}

func TestListAllInterfaces(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInterfacesTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, InterfacesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	interfaces, err := instances.ListInterfacesAll(client, instanceID)

	require.NoError(t, err)
	require.Len(t, interfaces, 1)
	require.Equal(t, PortID, interfaces[0].PortID)
	require.Equal(t, ExpectedInstanceInterfacesSlice, interfaces)
}

func TestListAllSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, SecurityGroupsListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	sgs, err := instances.ListSecurityGroupsAll(client, instanceID)

	require.NoError(t, err)
	require.Len(t, sgs, 1)
	require.Equal(t, SecurityGroup1, sgs[0])
	require.Equal(t, ExpectedSecurityGroupsSlice, sgs)
}

func TestListAllPorts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListPortsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, PortsListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	ports, err := instances.ListPortsAll(client, instanceID)

	require.NoError(t, err)
	require.Len(t, ports, 1)
	require.Equal(t, InstancePort1, ports[0])
	require.Equal(t, ExpectedPortsSlice, ports)
}

func TestUnAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUnAssignSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, UnAssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := instances.UnAssignSecurityGroup(client, instanceID, opts).ExtractErr()

	require.NoError(t, err)
}

func TestAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAssignSecurityGroupsTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, AssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := instances.AssignSecurityGroup(client, instanceID, opts).ExtractErr()

	require.NoError(t, err)
}

func TestAttachInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAttachInterfaceTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, AttachInterfaceRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.InterfaceOpts{
		Type:     types.SubnetInterfaceType,
		SubnetID: "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
	}

	tasks, err := instances.AttachInterface(client, instanceID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDetachInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDetachInterfaceTestURL(Instance1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, DetachInterfaceRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.InterfaceOpts{
		PortID:    "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		IpAddress: "192.168.0.23",
	}

	tasks, err := instances.DetachInterface(client, instanceID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateTestURLV2(), func(w http.ResponseWriter, r *http.Request) {
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

	options := instances.CreateOpts{
		Flavor:        "g1-standard-1-2",
		Names:         []string{"name"},
		NameTemplates: nil,
		Volumes: []instances.CreateVolumeOpts{{
			Source:     types.NewVolume,
			BootIndex:  0,
			Size:       10,
			TypeName:   volumes.Standard,
			Name:       "name",
			ImageID:    "",
			SnapshotID: "",
			VolumeID:   "",
			Metadata:   map[string]string{"meta1": "value1"},
		}},
		Interfaces: []instances.InterfaceInstanceCreateOpts{{InterfaceOpts: instances.InterfaceOpts{
			Type:      types.SubnetInterfaceType,
			NetworkID: "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			SubnetID:  "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			FloatingIP: &instances.CreateNewInterfaceFloatingIPOpts{
				Source:             types.ExistingFloatingIP,
				ExistingFloatingID: "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			},
		}, SecurityGroups: []gcorecloud.ItemID{}}},
		Keypair:  "keypair",
		Password: "password",
		Username: "username",
		UserData: "",
	}

	err := options.Validate()
	require.NoError(t, err)

	client := fake.ServiceTokenClient("instances", "v2")
	tasks, err := instances.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeleteTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := instances.DeleteOpts{
		Volumes:         nil,
		DeleteFloatings: true,
		FloatingIPs:     nil,
	}

	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("instances", "v1")
	tasks, err := instances.Delete(client, instanceID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestStart(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareStartTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.Start(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestStop(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareStopTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.Stop(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestPowerCycle(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(preparePowerCycleTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.PowerCycle(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestReboot(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRebootTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.Reboot(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestSuspend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareSuspendTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.Suspend(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestResume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResumeTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")
	instance, err := instances.Resume(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, Instance1, *instance)
}

func TestResize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareChangeFlavorTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ResizeResponse)
		if err != nil {
			log.Error(err)
		}
	})

	opts := instances.ChangeFlavorOpts{FlavorID: Instance1.Flavor.FlavorID}

	client := fake.ServiceTokenClient("instances", "v1")
	tasks, err := instances.Resize(client, instanceID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestMetadataListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	actual, err := instances.MetadataListAll(client, instanceID)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Metadata1, ct)
	require.Equal(t, ExpectedMetadataList, actual)
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataDetailsTestURL(instanceID, Metadata1.Key), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	actual, err := instances.MetadataGet(client, instanceID, Metadata1.Key).Extract()
	require.NoError(t, err)
	require.Equal(t, &Metadata1, actual)
}

func TestMetadataCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	opts := instances.MetadataSetOpts{
		Metadata: []instances.MetadataOpts{{
			Key:   "test1",
			Value: "test1",
		}, {
			Key:   "test2",
			Value: "test2",
		},
		}}

	client := fake.ServiceTokenClient("instances", "v1")
	err := instances.MetadataCreate(client, instanceID, opts).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	opts := instances.MetadataSetOpts{
		Metadata: []instances.MetadataOpts{{
			Key:   "test1",
			Value: "test1",
		}, {
			Key:   "test2",
			Value: "test2",
		},
		}}

	client := fake.ServiceTokenClient("instances", "v1")
	err := instances.MetadataUpdate(client, instanceID, opts).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataDetailsTestURL(instanceID, Metadata1.Key), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v1")
	err := instances.MetadataDelete(client, instanceID, Metadata1.Key).ExtractErr()
	require.NoError(t, err)
}

func TestGetSpiceConsole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetSpiceConsoleTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, InstanceConsoleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	actual, err := instances.GetSpiceConsole(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, &Console, actual)
}

func TestGetInstanceConsole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetInstanceConsoleTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, InstanceConsoleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	actual, err := instances.GetInstanceConsole(client, instanceID).Extract()
	require.NoError(t, err)
	require.Equal(t, &Console, actual)
}

func TestListInstanceMetrics(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstanceMetricsTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ListInstanceMetricsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListInstanceMetricsResponse)
		if err != nil {
			log.Error(err)
		}
	})

	opts := instances.ListMetricsOpts{
		TimeUnit:     types.DayMetricsTimeUnit,
		TimeInterval: 4,
	}

	client := fake.ServiceTokenClient("instances", "v1")
	result, err := instances.ListInstanceMetrics(client, instanceID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, len(result), 1)
	require.Equal(t, Metrics, result[0])
}

func TestListAvailableFlavors(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListAvailableFlavorsTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListAvailableFlavorsResponse)
		if err != nil {
			log.Error(err)
		}
	})

	includePrice := true
	opts := flavors.ListOpts{
		IncludePrices: &includePrice,
	}
	client := fake.ServiceTokenClient("instances", "v1")
	result, err := instances.ListAvailableFlavors(client, instanceID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, len(result), 1)
	require.Equal(t, AvailableFlavor, result[0])
}

func TestListInstanceLocation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstanceLocationTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListInstanceLocationResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v1")

	opts := instances.ListInstanceLocationOpts{}
	actual, err := instances.ListInstanceLocation(client, opts).Extract()
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, InstanceLocation, ct)
	require.Equal(t, ExpectedInstancesLocationSlice, actual)
}
