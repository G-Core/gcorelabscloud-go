package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const fileSharePath = "file_shares"

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareActionTestURLParams(projectID, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareExtendActionTestURL(id string) string { // nolint
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "extend")
}

func prepareListAccessRuleTestURLParams(projectID int, regionID int, shareID string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s", projectID, regionID, shareID, "access_rule")
}

func prepareGetAccessRuleTestURLParams(projectID int, regionID int, shareID string, ruleID string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s/%s", projectID, regionID, shareID, "access_rule", ruleID)
}

func prepareListAccessRuleTestURL(shareID string) string {
	return prepareListAccessRuleTestURLParams(fake.ProjectID, fake.RegionID, shareID)
}

func prepareGetAccessRuleTestURL(shareID string, ruleID string) string {
	return prepareGetAccessRuleTestURLParams(fake.ProjectID, fake.RegionID, shareID, ruleID)
}

func prepareMetadataTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/file_shares/%d/%d/%s/%s", "v1", fake.ProjectID, fake.RegionID, id, "metadata_item")
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

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	count := 0

	err := file_shares.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := file_shares.ExtractFileShares(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, ListFileShare1, ct)
		require.Equal(t, ExpectedFileShareSlice, actual)
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

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	actual, err := file_shares.ListAll(client, nil)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, ListFileShare1, ct)
	require.Equal(t, ExpectedFileShareSlice, actual)

}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(FileShare1.ID)

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

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	ct, err := file_shares.Get(client, FileShare1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, FileShare1, *ct)
	require.Equal(t, &createdTime, ct.CreatedAt)
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

	options := file_shares.CreateOpts{
		Name:     FileShare1.Name,
		Protocol: "NFS",
		Size:     13,
		Network: file_shares.FileShareNetworkOpts{
			NetworkID: "9b17dd07-1281-4fe0-8c13-d80c5725e297",
			SubnetID:  "221f8318-cf2d-47a7-90f7-97acfa4ef165",
		},
		Metadata: map[string]string{"qqq": "that"},
	}

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	tasks, err := file_shares.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	tasks, err := file_shares.Delete(client, FileShare1.ID).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(FileShare1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UpdateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	opts := file_shares.UpdateOpts{
		Name: "myshareqqq",
	}
	ct, err := file_shares.Update(client, FileShare1.ID, opts).Extract()

	FileShare1.Name = opts.Name
	require.NoError(t, err)
	require.Equal(t, FileShare1, *ct)
	require.Equal(t, FileShare1.Name, ct.Name)
	require.Equal(t, &createdTime, ct.CreatedAt)
}

func TestExtend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareExtendActionTestURL(FileShare1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ExtendRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ExtendResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	opts := file_shares.ExtendOpts{
		Size: 15,
	}
	ct, err := file_shares.Extend(client, FileShare1.ID, opts).Extract()

	FileShare1.Size = opts.Size
	require.NoError(t, err)
	require.Equal(t, Tasks1, *ct)
}

// Test file share access rules

func TestListAccessRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListAccessRuleTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListAccessRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	count := 0

	err := file_shares.ListAccessRules(client, FileShare1.ID, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := file_shares.ExtractAccessRule(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, AccessRule1, ct)
		require.Equal(t, ExpectedAccessRuleSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateAccessRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListAccessRuleTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateAccessRuleRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateAccessRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := file_shares.CreateAccessRuleOpts{
		IPAddress:  "10.100.100.0/24",
		AccessMode: "rw",
	}

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	result, err := file_shares.CreateAccessRule(client, FileShare1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, &CreatedAccessRule, result)
}

func TestDeleteAccessRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetAccessRuleTestURL(FileShare1.ID, AccessRule1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	err := file_shares.DeleteAccessRule(client, FileShare1.ID, AccessRule1.ID).ExtractErr()
	require.NoError(t, err)
}

// Metadata tests

func TestMetadataListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	actual, err := metadata.MetadataListAll(client, FileShare1.ID)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Metadata1, ct)
	require.Equal(t, ExpectedMetadataList, actual)
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")

	actual, err := metadata.MetadataGet(client, FileShare1.ID, Metadata2.Key).Extract()
	require.NoError(t, err)
	require.Equal(t, &Metadata2, actual)
}

func TestMetadataCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	err := metadata.MetadataCreateOrUpdate(client, FileShare1.ID, map[string]string{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	err := metadata.MetadataReplace(client, FileShare1.ID, map[string]string{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(FileShare1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient(fileSharePath, "v1")
	err := metadata.MetadataDelete(client, FileShare1.ID, Metadata1.Key).ExtractErr()
	require.NoError(t, err)
}
