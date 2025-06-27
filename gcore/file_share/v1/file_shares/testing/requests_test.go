package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fakeclient "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	"github.comcom/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const fileSharePath = "/fileshares/v1/shares"
const fileShareListPath = "/fileshares/v1/shares"

var (
	client = client.NewClient()
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(client.ProjectID, client.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(client.ProjectID, client.RegionID, id)
}

func prepareActionTestURLParams(projectID, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareExtendActionTestURL(id string) string { // nolint
	return prepareActionTestURLParams(client.ProjectID, client.RegionID, id, "extend")
}

func prepareListAccessRuleTestURLParams(projectID int, regionID int, shareID string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s", projectID, regionID, shareID, "access_rule")
}

func prepareGetAccessRuleTestURLParams(projectID int, regionID int, shareID string, ruleID string) string {
	return fmt.Sprintf("/v1/file_shares/%d/%d/%s/%s/%s", projectID, regionID, shareID, "access_rule", ruleID)
}

func prepareListAccessRuleTestURL(shareID string) string {
	return prepareListAccessRuleTestURLParams(client.ProjectID, client.RegionID, shareID)
}

func prepareGetAccessRuleTestURL(shareID string, ruleID string) string {
	return prepareGetAccessRuleTestURLParams(client.ProjectID, client.RegionID, shareID, ruleID)
}

func prepareMetadataTestURL(id string) string {
	return prepareActionTestURLParams(client.ProjectID, client.RegionID, id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/file_shares/%d/%d/%s/%s", "v1", client.ProjectID, client.RegionID, id, "metadata_item")
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fileShareListPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ListResponse)
	})

	pages, err := file_shares.List(client).AllPages()
	th.AssertNoErr(t, err)

	actual, err := file_shares.ExtractFileShares(pages)
	th.AssertNoErr(t, err)

	expected := []file_shares.FileShare{FirstFileShare}

	th.CheckDeepEquals(t, expected, actual)
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fileShareListPath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ListResponse)
	})

	all, err := file_shares.ListAll(client)
	th.AssertNoErr(t, err)

	expected := []file_shares.FileShare{FirstFileShare}
	th.CheckDeepEquals(t, expected, all)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	actual, err := file_shares.Get(client, FirstFileShare.ID).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &SecondFileShare, actual)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fileSharePath, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, CreateResponse)
	})

	createOpts := file_shares.CreateOpts{
		Name:     "myshare",
		Protocol: "NFS",
		Size:     13,
		Network: &file_shares.FileShareNetworkOpts{
			NetworkID: "9b17dd07-1281-4fe0-8c13-d80c5725e297",
			SubnetID:  "221f8318-cf2d-47a7-90f7-97acfa4ef165",
		},
		Tags: map[string]string{
			"qqq": "that",
		},
	}
	res := file_shares.Create(client, createOpts)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	var expectedTasks = tasks.TaskResults{
		Tasks: []tasks.TaskID{
			"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}
	th.CheckDeepEquals(t, &expectedTasks, task)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, DeleteResponse)
	})

	res := file_shares.Delete(client, FirstFileShare.ID)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	var expectedTasks = tasks.TaskResults{
		Tasks: []tasks.TaskID{
			"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}
	th.CheckDeepEquals(t, &expectedTasks, task)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	updateOpts := file_shares.UpdateOpts{
		Name: "myshareqqq",
	}

	actual, err := file_shares.Update(client, FirstFileShare.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &UpdatedFileShare, actual)
}

func TestExtend(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/action", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))
		th.TestJSONRequest(t, r, ExtendRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, ExtendResponse)
	})

	extendOpts := file_shares.ExtendOpts{
		Size: 15,
	}

	res := file_shares.Extend(client, FirstFileShare.ID, extendOpts)
	th.AssertNoErr(t, res.Err)

	task, err := res.Extract()
	th.AssertNoErr(t, err)

	var expectedTasks = tasks.TaskResults{
		Tasks: []tasks.TaskID{
			"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54",
		},
	}
	th.CheckDeepEquals(t, &expectedTasks, task)
}

// Test file share access rules

func TestListAccessRules(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/access-rules", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListAccessRuleResponse)
	})

	pages, err := file_shares.ListAccessRules(client, FirstFileShare.ID).AllPages()
	th.AssertNoErr(t, err)

	actual, err := file_shares.ExtractAccessRule(pages)
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []file_shares.AccessRule{FirstAccessRule}, actual)
}

func TestCreateAccessRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/access-rules", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, CreateAccessRuleResponse)
	})

	createOpts := file_shares.CreateAccessRuleOpts{
		IPAddress:  "10.100.100.0/24",
		AccessMode: "rw",
	}

	actual, err := file_shares.CreateAccessRule(client, FirstFileShare.ID, createOpts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &SecondAccessRule, actual)
}

func TestDeleteAccessRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/access-rules/%s", fileSharePath, FirstFileShare.ID, FirstAccessRule.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.WriteHeader(http.StatusNoContent)
	})

	res := file_shares.DeleteAccessRule(client, FirstFileShare.ID, FirstAccessRule.ID)
	th.AssertNoErr(t, res.Err)
}

// Metadata tests

func TestMetadataList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/metadata", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MetadataListResponse)
	})

	pages, err := file_shares.MetadataList(client, FirstFileShare.ID).AllPages()
	th.AssertNoErr(t, err)

	actual, err := file_shares.ExtractMetadata(pages)
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []metadata.Metadata{FirstMetadata, SecondMetadata}, actual)
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/metadata/%s", fileSharePath, FirstFileShare.ID, SecondMetadata.Key), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MetadataResponse)
	})

	actual, err := file_shares.MetadataGet(client, FirstFileShare.ID, SecondMetadata.Key).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &SecondMetadata, actual)
}

func TestMetadataCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/metadata", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MetadataListResponse)
	})

	res := file_shares.MetadataCreateOrUpdate(client, FirstFileShare.ID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	})
	th.AssertNoErr(t, res.Err)
}

func TestMetadataUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/metadata", fileSharePath, FirstFileShare.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, MetadataListResponse)
	})

	res := file_shares.MetadataCreateOrUpdate(client, FirstFileShare.ID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	})
	th.AssertNoErr(t, res.Err)
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	client := fakeclient.ServiceTokenClient(fileSharePath, "v1")

	th.Mux.HandleFunc(fmt.Sprintf("%s/%s/metadata/%s", fileSharePath, FirstFileShare.ID, "test"), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fakeclient.AccessToken))

		w.WriteHeader(http.StatusNoContent)
	})

	res := file_shares.MetadataDelete(client, FirstFileShare.ID, "test")
	th.AssertNoErr(t, res.Err)
}
