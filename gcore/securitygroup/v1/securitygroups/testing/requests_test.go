package testing

import (
	"fmt"
	"net/http"
	"testing"

	instancestesting "github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/securitygroups/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, id)
}

func prepareActionTestURLParams(projectID int, regionID int, id, action string) string {
	return fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, id, action)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareAddRuleTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "rules")
}

func prepareListInstancesTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "instances")
}

func prepareDeepCopyTestURL(id string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, id, "copy")
}

func prepareGetActionTestURLParams(version string, id string, action string) string { // nolint
	return fmt.Sprintf("/%s/securitygroups/%d/%d/%s/%s", version, fake.ProjectID, fake.RegionID, id, action)
}

func prepareMetadataTestURL(id string) string {
	return prepareGetActionTestURLParams("v1", id, "metadata")
}

func prepareMetadataItemTestURL(id string) string {
	return fmt.Sprintf("/%s/securitygroups/%d/%d/%s/%s", "v1", fake.ProjectID, fake.RegionID, id, "metadata_item")
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

	client := fake.ServiceTokenClient("securitygroups", "v1")
	count := 0

	err := securitygroups.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := securitygroups.ExtractSecurityGroups(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, SecurityGroup1, ct)
		require.Equal(t, ExpectedSecurityGroupSlice, actual)
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

	client := fake.ServiceTokenClient("securitygroups", "v1")

	groups, err := securitygroups.ListAll(client, nil)
	require.NoError(t, err)
	ct := groups[0]
	require.Equal(t, SecurityGroup1, ct)
	require.Equal(t, ExpectedSecurityGroupSlice, groups)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(SecurityGroup1.ID)

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

	client := fake.ServiceTokenClient("securitygroups", "v1")

	ct, err := securitygroups.Get(client, SecurityGroup1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, SecurityGroup1, *ct)
	require.Equal(t, groupCreatedTime, ct.CreatedAt)
	require.Equal(t, groupUpdatedTime, *ct.UpdatedAt)

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

	options := securitygroups.CreateOpts{
		SecurityGroup: securitygroups.CreateSecurityGroupOpts{
			Name:               groupName,
			Description:        &groupDescription,
			SecurityGroupRules: []securitygroups.CreateSecurityGroupRuleOpts{},
		},
		Instances: []string{
			"8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
			"a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
		},
	}

	client := fake.ServiceTokenClient("securitygroups", "v1")
	group, err := securitygroups.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, SecurityGroup1, *group)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(SecurityGroup1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")
	err := securitygroups.Delete(client, SecurityGroup1.ID).ExtractErr()
	require.NoError(t, err)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(SecurityGroup1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
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

	client := fake.ServiceTokenClient("securitygroups", "v1")

	opts := securitygroups.UpdateOpts{
		Name: SecurityGroup1.Name,
	}

	ct, err := securitygroups.Update(client, SecurityGroup1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, SecurityGroup1, *ct)

}

func TestCreateRule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAddRuleTestURL(SecurityGroup1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRuleRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := securitygroups.CreateSecurityGroupRuleOpts{
		Direction:   direction,
		EtherType:   eitherType,
		Protocol:    protocol,
		Description: &groupDescription,
	}

	client := fake.ServiceTokenClient("securitygroups", "v1")
	rule, err := securitygroups.AddRule(client, SecurityGroup1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, securityGroupRule1, *rule)
}

func TestListAllInstances(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListInstancesTestURL(SecurityGroup1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, instancestesting.ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")

	results, err := securitygroups.ListAllInstances(client, SecurityGroup1.ID)
	require.NoError(t, err)
	instance := results[0]
	require.Equal(t, instancestesting.Instance1, instance)
	require.Equal(t, instancestesting.ExpectedInstancesSlice, results)
}

func TestIDFromName(t *testing.T) {
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

	client := fake.ServiceTokenClient("securitygroups", "v1")

	groupID, err := securitygroups.IDFromName(client, SecurityGroup1.Name)
	require.NoError(t, err)
	require.Equal(t, SecurityGroup1.ID, groupID)

}

func TestDeepCopy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeepCopyTestURL(SecurityGroup1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.WriteHeader(http.StatusOK)
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")

	opts := securitygroups.DeepCopyOpts{Name: "default"}
	err := securitygroups.DeepCopy(client, SecurityGroup1.ID, opts).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(groupID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")

	actual, err := securitygroups.MetadataListAll(client, groupID)
	require.NoError(t, err)
	ct := actual[0]
	require.Equal(t, Metadata1, ct)
	require.Equal(t, ExpectedMetadataList, actual)
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(groupID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")

	actual, err := securitygroups.MetadataGet(client, groupID, Metadata1.Key).Extract()
	require.NoError(t, err)
	require.Equal(t, &Metadata1, actual)
}

func TestMetadataCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(groupID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")
	err := securitygroups.MetadataCreateOrUpdate(client, groupID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataTestURL(groupID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, MetadataCreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")
	err := securitygroups.MetadataReplace(client, groupID, map[string]interface{}{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()
	require.NoError(t, err)
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareMetadataItemTestURL(groupID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("securitygroups", "v1")
	err := securitygroups.MetadataDelete(client, groupID, Metadata1.Key).ExtractErr()
	require.NoError(t, err)
}
