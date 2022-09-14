package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	gtesting "github.com/G-Core/gcorelabscloud-go/gcore/utils/testing"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string { // nolint
	return fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareUpdateTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareAttachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "attach")
}

func prepareDetachTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "detach")
}

func prepareRetypeTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "retype")
}

func prepareExtendTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "extend")
}

func prepareRevertTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "revert")
}

func prepareGetActionTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/floatingips/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}
func prepareMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/floatingips/%d/%d/%s/%s", "v1", fake.ProjectID, fake.RegionID, id, "metadata_item")
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

	client := fake.ServiceTokenClient("volumes", "v1")
	count := 0

	opts := volumes.ListOpts{}

	err := volumes.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := volumes.ExtractVolumes(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Volume1, ct)
		require.Equal(t, ExpectedVolumeSlice, actual)
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

	client := fake.ServiceTokenClient("volumes", "v1")
	opts := volumes.ListOpts{}

	actual, err := volumes.ListAll(client, opts)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Volume1, ct)
	require.Equal(t, ExpectedVolumeSlice, actual)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Volume1.ID)

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

	client := fake.ServiceTokenClient("volumes", "v1")

	ct, err := volumes.Get(client, Volume1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Volume1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)

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

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	size := 10
	typeName := volumes.SsdHiIops
	instanceIDToAttachTo := "88f3e0bd-ca86-4cf7-be8b-dd2988e23c2d"

	options := volumes.CreateOpts{
		Source:               "new-volume",
		Name:                 "TestVM5 Ubuntu volume",
		Size:                 size,
		TypeName:             typeName,
		ImageID:              "",
		SnapshotID:           "",
		InstanceIDToAttachTo: instanceIDToAttachTo,
	}

	err := options.Validate()
	require.NoError(t, err)
	client := fake.ServiceTokenClient("volumes", "v1")
	tasks, err := volumes.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("volumes", "v1")

	opts := volumes.DeleteOpts{Snapshots: []string{"x", "y"}}

	tasks, err := volumes.Delete(client, Volume1.ID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUpdateTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.UpdateOpts{
		Name: "updated",
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Update(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
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

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Attach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
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

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.InstanceOperationOpts{
		InstanceID: Volume1.Attachments[0].ServerID,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Detach(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
}

func TestRetype(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRetypeTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RetypeRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := volumes.VolumeTypePropertyOperationOpts{
		VolumeType: volumes.SsdHiIops,
	}

	client := fake.ServiceTokenClient("volumes", "v1")
	volume, err := volumes.Retype(client, Volume1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Volume1, *volume)
}

func TestExtend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareExtendTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ExtendRequest)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ExtendResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("volumes", "v1")
	opts := volumes.SizePropertyOperationOpts{Size: 16}

	tasks, err := volumes.Extend(client, Volume1.ID, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestRevert(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRevertTestURL(Volume1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ExtendResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("volumes", "v1")
	tasks, err := volumes.Revert(client, Volume1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestMetadataListAll(t *testing.T) {
	gtesting.BuildTestMetadataListAll("volumes", Volume1.ID)(t)
}

func TestMetadataGet(t *testing.T) {
	gtesting.BuildTestMetadataGet("volumes", Volume1.ID)(t)
}

func TestMetadataCreate(t *testing.T) {
	gtesting.BuildTestMetadataCreate("volumes", Volume1.ID)(t)
}

func TestMetadataUpdate(t *testing.T) {
	gtesting.BuildTestMetadataUpdate("volumes", Volume1.ID)(t)

}

func TestMetadataDelete(t *testing.T) {
	gtesting.BuildTestMetadataDelete("volumes", Volume1.ID)(t)
}
