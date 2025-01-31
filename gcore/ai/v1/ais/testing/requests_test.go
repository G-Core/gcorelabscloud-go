package testing

import (
	"fmt"
	"net/http"
	"testing"

	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/ai/clusters/%d/%d", version, projectID, regionID)
}

func prepareGetTestURLParams(version string, projectID int, regionID int, id string) string {
	return fmt.Sprintf("/%s/ai/clusters/%d/%d/%s", version, projectID, regionID, id)
}

func prepareGetActionTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/ai/clusters/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}

func prepareGetActionGPUTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/ai/clusters/gpu/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}

func prepareListTestURL() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
}

func prepareListInterfacesTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "interfaces")
}

func prepareGetInstanceConsoleTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "get_console")
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

func prepareInstancePowerCycleTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "powercycle")
}

func prepareClusterPowerCycleTestURL(id string) string {
	return prepareGetActionTestURLParams("v2", id, "powercycle")
}

func prepareInstanceRebootTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "reboot")
}
func prepareClusterRebootTestURL(id string) string {
	return prepareGetActionTestURLParams("v2", id, "reboot")
}

func prepareSuspendTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "suspend")
}

func prepareResumeTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "resume")
}

func prepareResizeTestURL(id string) string {
	return prepareGetActionGPUTestURLParams("v1", id, "resize")
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareDeleteTestURL(id string) string {
	return prepareGetTestURLParams("v1", fake.ProjectID, fake.RegionID, id)
}

func prepareCreateTestURL() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
}

func prepareListMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metadata")
}

func prepareMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v2", id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/ai/clusters/%d/%d/%s/%s", "v2", fake.ProjectID, fake.RegionID, id, "metadata_item")
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

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	actual, err := ai.ListAll(client)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, AICluster1, ct)
	require.Equal(t, ExpectedAIClusterSlice, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(AICluster1.ClusterID)

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

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	ct, err := ai.Get(client, AICluster1.ClusterID).Extract()

	require.NoError(t, err)
	require.Equal(t, AICluster1, *ct)

}

func TestListAllInterfaces(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInterfacesTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ClusterInterfacesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	interfaces, err := ai.ListInterfacesAll(client, AICluster1.ClusterID)

	require.NoError(t, err)
	require.Len(t, interfaces, 1)
	require.Equal(t, PortID, interfaces[0].PortID)
	require.Equal(t, ExpectedAIClusterInterfacesSlice, interfaces)
}

func TestListAllPorts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListPortsTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, PortsListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	ports, err := ai.ListPortsAll(client, AICluster1.ClusterID)

	require.NoError(t, err)
	require.Len(t, ports, 1)
	require.Equal(t, AIClusterPort1, ports[0])
	require.Equal(t, ExpectedPortsSlice, ports)
}

func TestUnAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUnAssignSecurityGroupsTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, UnAssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := ai.UnAssignSecurityGroup(client, AICluster1.ClusterID, opts).ExtractErr()

	require.NoError(t, err)
}

func TestAssignSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAssignSecurityGroupsTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestJSONRequest(t, r, AssignSecurityGroupsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	opts := instances.SecurityGroupOpts{
		Name: "Test",
	}

	err := ai.AssignSecurityGroup(client, AICluster1.ClusterID, opts).ExtractErr()

	require.NoError(t, err)
}

func TestAttachInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	aiInstanceID := AICluster1.PoplarServer[0].ID
	th.Mux.HandleFunc(prepareAttachInterfaceTestURL(AICluster1.PoplarServer[0].ID), func(w http.ResponseWriter, r *http.Request) {
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

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	opts := ai.AttachInterfaceOpts{
		Type:     types.SubnetInterfaceType,
		SubnetID: "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
	}

	tasks, err := ai.AttachAIInstanceInterface(client, aiInstanceID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDetachInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	aiInstanceID := AICluster1.PoplarServer[0].ID
	th.Mux.HandleFunc(prepareDetachInterfaceTestURL(aiInstanceID), func(w http.ResponseWriter, r *http.Request) {
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

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	opts := ai.DetachInterfaceOpts{
		PortID:    "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		IpAddress: "192.168.0.23",
	}

	tasks, err := ai.DetachAIInstanceInterface(client, aiInstanceID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	url := prepareCreateTestURL()
	fmt.Println(url)
	th.Mux.HandleFunc(prepareCreateTestURL(), func(w http.ResponseWriter, r *http.Request) {
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

	options := ai.CreateOpts{
		Flavor:  "g2a-ai-fake-v1pod-8",
		Name:    "ivandts",
		ImageID: "06e62653-1f88-4d38-9aa6-62833e812b4f",
		Volumes: []instances.CreateVolumeOpts{
			{
				Source:    types.Image,
				BootIndex: 0,
				Size:      20,
				TypeName:  volumes.Standard,
				ImageID:   "06e62653-1f88-4d38-9aa6-62833e812b4f",
			},
		},
		Interfaces: []instances.InterfaceInstanceCreateOpts{
			{
				InterfaceOpts: instances.InterfaceOpts{
					Type:      types.AnySubnetInterfaceType,
					NetworkID: "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
				},
			},
		},
		Password: "secret",
		Username: "useruser",
	}
	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("ai/clusters", "v1")
	tasks, err := ai.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestResize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	url := prepareCreateTestURL()
	fmt.Println(url)
	th.Mux.HandleFunc(prepareResizeTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResizeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := ai.ResizeGPUAIClusterOpts{
		InstancesCount: 2,
	}
	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("ai/clusters/gpu", "v1")
	tasks, err := ai.Resize(client, AICluster1.ClusterID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeleteTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
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

	options := ai.DeleteOpts{
		Volumes:          nil,
		DeleteFloatings:  true,
		FloatingIPs:      nil,
		ReservedFixedIPs: nil,
	}

	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("ai/clusters", "v1")
	tasks, err := ai.Delete(client, AICluster1.ClusterID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestPowerCycleCluster(t *testing.T) {
	th.SetupHTTP()
	th.Mux.HandleFunc(prepareClusterPowerCycleTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AIClusterPowercycleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")
	result, err := ai.PowerCycleAICluster(client, AICluster1.ClusterID).Extract()
	require.NoError(t, err)

	require.Equal(t, AICluster1.PoplarServer, result)
}

func TestPowerCycleInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareInstancePowerCycleTestURL(AICluster1.PoplarServer[0].ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AIInstancePowercycleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	instance, err := ai.PowerCycleAIInstance(client, AICluster1.PoplarServer[0].ID).Extract()
	require.NoError(t, err)
	require.Equal(t, &AICluster1.PoplarServer[0], instance)
}

func TestRebootCluster(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareClusterRebootTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AIClusterRebootResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")
	result, err := ai.RebootAICluster(client, AICluster1.ClusterID).Extract()
	require.NoError(t, err)
	require.Equal(t, AICluster1.PoplarServer, result)
}

func TestRebootInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareInstanceRebootTestURL(AICluster1.PoplarServer[0].ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AIInstanceRebootResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	instance, err := instances.Reboot(client, AICluster1.PoplarServer[0].ID).Extract()
	require.NoError(t, err)
	require.Equal(t, &AICluster1.PoplarServer[0], instance)
}

func TestSuspend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareSuspendTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	tasks, err := ai.Suspend(client, AICluster1.ClusterID).Extract()
	require.NoError(t, err)
	require.Equal(t, &Tasks1, tasks)
}

func TestResume(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareResumeTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")
	tasks, err := ai.Resume(client, AICluster1.ClusterID).Extract()
	require.NoError(t, err)
	require.Equal(t, &Tasks1, tasks)
}

func TestGetInstanceConsole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetInstanceConsoleTestURL(AICluster1.PoplarServer[0].ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, InstanceConsoleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	actual, err := ai.GetInstanceConsole(client, AICluster1.PoplarServer[0].ID).Extract()
	require.NoError(t, err)
	require.Equal(t, &Console, actual)
}

// Metadata tests

func TestMetadataListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListMetadataTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v1")

	actual, err := ai.MetadataListAll(client, AICluster1.ClusterID)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Metadata1, ct)
	require.Equal(t, ExpectedMetadataList, actual)
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")

	actual, err := ai.MetadataGet(client, AICluster1.ClusterID, Metadata1.Key).Extract()
	require.NoError(t, err)
	require.Equal(t, &Metadata1, actual)
}

func TestMetadataCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")
	err := ai.MetadataCreateOrUpdate(client, AICluster1.ClusterID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")
	err := ai.MetadataReplace(client, AICluster1.ClusterID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc(prepareMetadataItemTestURL(AICluster1.ClusterID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("ai/clusters", "v2")
	err := ai.MetadataDelete(client, AICluster1.ClusterID, Metadata1.Key).ExtractErr()
	require.NoError(t, err)
}
