package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygrouprules"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/securitygrouprules/%d/%d/%s", projectID, regionID, id)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(securityGroupRule1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusNoContent)

	})

	client := fake.ServiceTokenClient("securitygrouprules", "v1")
	err := securitygrouprules.Delete(client, securityGroupRule1.ID).ExtractErr()
	require.NoError(t, err)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(securityGroupRule1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ReplaceRuleRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, ReplaceRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("securitygrouprules", "v1")

	opts := securitygroups.CreateSecurityGroupRuleOpts{
		Direction:   direction,
		EtherType:   eitherType,
		Protocol:    protocol,
		Description: &groupDescription,
	}

	ct, err := securitygrouprules.Replace(client, securityGroupRule1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, securityGroupRule1, *ct)

}
