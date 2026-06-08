package testing

import (
	"fmt"
	"net/http"
	"testing"

	instancesV2 "github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/types"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/instances/%d/%d/%s/%s", "v2", fake.ProjectID, fake.RegionID, id, "metadata_item")
}

func prepareActionTestURL(id string) string {
	return fmt.Sprintf("/%s/instances/%d/%d/%s/%s", "v2", fake.ProjectID, fake.RegionID, id, "action")
}

func TestAction(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareActionTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{"action": "start"}`)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ActionResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v2")
	opts := instancesV2.ActionOpts{Action: types.InstanceActionTypeStart}
	actual, err := instancesV2.Action(client, instanceID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, &ActionTasks, actual)
}

func TestMetadataItemGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("instances", "v2")

	actual, err := instancesV2.MetadataItemGet(client, instanceID, instancesV2.MetadataItemOpts{Key: Metadata.Key}).Extract()
	require.NoError(t, err)
	require.Equal(t, &Metadata, actual)
}

func TestMetadataItemDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("instances", "v2")
	err := instancesV2.MetadataItemDelete(client, instanceID, instancesV2.MetadataItemOpts{Key: Metadata.Key}).ExtractErr()
	require.NoError(t, err)
}
