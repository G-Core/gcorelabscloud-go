package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bminstances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/bminstances/%d/%d", version, projectID, regionID)
}

func prepareCreateTestURLV1() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
}

func prepareRebuildTestURLParams(id string) string { // nolint
	return fmt.Sprintf("/v1/bminstances/%d/%d/%s/rebuild", fake.ProjectID, fake.RegionID, id)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateTestURLV1(), func(w http.ResponseWriter, r *http.Request) {
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

	options := bminstances.CreateOpts{
		Flavor:        "bm1-infrastructure-small",
		Names:         []string{"name"},
		NameTemplates: nil,
		ImageID:       "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
		Interfaces: []bminstances.InterfaceOpts{{
			Type:      types.SubnetInterfaceType,
			NetworkID: "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			SubnetID:  "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			FloatingIP: &bminstances.CreateNewInterfaceFloatingIPOpts{
				Source:             types.ExistingFloatingIP,
				ExistingFloatingID: "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
			},
		}},
		Keypair:  "keypair",
		Password: "password",
		Username: "username",
	}

	err := options.Validate()
	require.NoError(t, err)

	client := fake.ServiceTokenClient("bminstances", "v1")
	tasks, err := bminstances.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestRebuild(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRebuildTestURLParams(instanceID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RebuildRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("bminstances", "v1")
	opts := bminstances.RebuildInstanceOpts{
		ImageID: imageID,
	}
	tasks, err := bminstances.Rebuild(client, instanceID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}
