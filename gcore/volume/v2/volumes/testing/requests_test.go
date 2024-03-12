package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v2/volumes"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string { // nolint
	return fmt.Sprintf("/v2/volumes/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareAttachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "attach")
}

func prepareDetachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "detach")
}

func TestAttach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAttachTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AttachDetachRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, AttachResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v2")
	tasks, err := volumes.Attach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDetach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDetachTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AttachDetachRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, DetachResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v2")
	tasks, err := volumes.Detach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}
