package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/snapshot/v1/snapshots"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/snapshots/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/snapshots/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
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

	client := fake.ServiceTokenClient("snapshots", "v1")
	count := 0

	err := snapshots.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := snapshots.ExtractSnapshots(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Snapshot1, ct)
		require.Equal(t, ExpectedSnapshotSlice, actual)
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

	client := fake.ServiceTokenClient("snapshots", "v1")

	actual, err := snapshots.ListAll(client, nil)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Snapshot1, ct)
	require.Equal(t, ExpectedSnapshotSlice, actual)

}

func TestListQuery(t *testing.T) {
	opts := snapshots.ListOpts{
		VolumeID:   "",
		InstanceID: "",
	}
	res, err := opts.ToSnapshotListQuery()
	require.NoError(t, err)
	require.Equal(t, "", res)
}

func TestCreateOptsToMapWithoutMandatory(t *testing.T) {
	opts := snapshots.CreateOpts{
		VolumeID:    "",
		Name:        "",
		Description: "",
	}
	_, err := opts.ToSnapshotCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "VolumeID")
}

func TestCreateOptsToMap(t *testing.T) {
	opts := snapshots.CreateOpts{
		VolumeID:    "x",
		Name:        "y",
		Description: "",
	}
	res, err := opts.ToSnapshotCreateMap()
	require.NoError(t, err)
	require.Contains(t, res, "volume_id")
	require.NotContains(t, res, "description")
	require.Contains(t, res, "name")
	require.Len(t, res, 2)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Snapshot1.ID)

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

	client := fake.ServiceTokenClient("snapshots", "v1")

	ct, err := snapshots.Get(client, Snapshot1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Snapshot1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, *ct.UpdatedAt)

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

	options := snapshots.CreateOpts{
		VolumeID:    "67baa7d1-08ea-4fc5-bef2-6b2465b7d227",
		Name:        Snapshot1.Name,
		Description: "after boot",
	}

	client := fake.ServiceTokenClient("snapshots", "v1")
	tasks, err := snapshots.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Snapshot1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("snapshots", "v1")
	tasks, err := snapshots.Delete(client, Snapshot1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestFindByName(t *testing.T) {
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

	client := fake.ServiceTokenClient("snapshots", "v1")

	opts := snapshots.ListOpts{
		VolumeID:   "",
		InstanceID: "",
	}

	snapShopID, err := snapshots.IDFromName(client, Snapshot1.Name, opts)
	require.NoError(t, err)
	require.Equal(t, Snapshot1.ID, snapShopID)
	_, err = snapshots.IDFromName(client, "X", opts)
	require.Error(t, err)

}
