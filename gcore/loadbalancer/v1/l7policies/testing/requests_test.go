package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

//v1/l7policies/{project_id}/{region_id}/{l7policy_id}/rules
func prepareListPolicyTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/l7policies/%d/%d", projectID, regionID)
}

func prepareGetPolicyTestURLParams(projectID int, regionID int, pid string) string {
	return fmt.Sprintf("/v1/l7policies/%d/%d/%s", projectID, regionID, pid)
}

func prepareListTestURL() string {
	return prepareListPolicyTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetPolicyTestURL(id string) string {
	return prepareGetPolicyTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareListRuleTestURLParams(projectID int, regionID int, pid string) string {
	return fmt.Sprintf("/v1/l7policies/%d/%d/%s/rules", projectID, regionID, pid)
}

func prepareGetRuleTestURLParams(projectID int, regionID int, pid, rid string) string {
	return fmt.Sprintf("/v1/l7policies/%d/%d/%s/rules/%s", projectID, regionID, pid, rid)
}

func prepareListRuleTestURL(pid string) string {
	return prepareListRuleTestURLParams(fake.ProjectID, fake.RegionID, pid)
}

func prepareGetRuleTestURL(pid, rid string) string {
	return prepareGetRuleTestURLParams(fake.ProjectID, fake.RegionID, pid, rid)
}

func TestPolicyList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListPolicyResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")
	count := 0

	err := l7policies.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := l7policies.ExtractL7Polices(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Policy, ct)
		require.Equal(t, ExpectedPolicySlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestPolicyListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListPolicyResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	policies, err := l7policies.ListAll(client)
	require.NoError(t, err)
	ct := policies[0]
	require.Equal(t, Policy, ct)
	require.Equal(t, ExpectedPolicySlice, policies)

	th.AssertNoErr(t, err)

	if len(policies) != 1 {
		t.Errorf("Expected 1 page, got %d", len(policies))
	}
}

func TestPolicyGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetPolicyTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, GetPolicyResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	p, err := l7policies.Get(client, pid).Extract()
	th.AssertNoErr(t, err)

	require.Equal(t, Policy, *p)
}

func TestPolicyCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePolicyRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := l7policies.CreateOpts{
		Action:           l7policies.ActionRedirectToURL,
		ListenerID:       "023f2e34-7806-443b-bfae-16c324569a3d",
		Name:             "redirect-example.com",
		Position:         1,
		RedirectHTTPCode: 301,
		RedirectURL:      "http://www.example.com",
		Tags:             []string{"test_tag"},
	}

	client := fake.ServiceTokenClient("l7policies", "v1")
	tasks, err := l7policies.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestPolicyReplace(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetPolicyTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ReplacePolicyRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := l7policies.ReplaceOpts{
		Action:           l7policies.ActionRedirectToURL,
		Name:             "redirect-example.com",
		Position:         1,
		RedirectHTTPCode: 301,
		RedirectURL:      "http://www.example.com",
		Tags:             []string{"test_tag"},
	}

	client := fake.ServiceTokenClient("l7policies", "v1")
	tasks, err := l7policies.Replace(client, pid, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestPolicyDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetPolicyTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	task, err := l7policies.Delete(client, pid).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *task)
}

func TestRuleList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListRuleTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")
	count := 0

	err := l7policies.ListRule(client, pid).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := l7policies.ExtractL7Rules(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Rule, ct)
		require.Equal(t, ExpectedRuleSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestRuleListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListRuleTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	rules, err := l7policies.ListAllRule(client, pid)
	require.NoError(t, err)
	ct := rules[0]
	require.Equal(t, Rule, ct)
	require.Equal(t, ExpectedRuleSlice, rules)

	if len(rules) != 1 {
		t.Errorf("Expected 1 page, got %d", len(rules))
	}
}

func TestRuleGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetRuleTestURL(pid, rid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, GetRuleResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	r, err := l7policies.GetRule(client, pid, rid).Extract()
	th.AssertNoErr(t, err)

	require.Equal(t, Rule, *r)
}

func TestRuleCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListRuleTestURL(pid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRuleRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := l7policies.CreateRuleOpts{
		CompareType: l7policies.CompareTypeRegex,
		Invert:      false,
		Type:        l7policies.TypePath,
		Value:       "/images*",
		Tags:        []string{"test_tag"},
	}

	client := fake.ServiceTokenClient("l7policies", "v1")
	tasks, err := l7policies.CreateRule(client, pid, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestRuleDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetRuleTestURL(pid, rid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("l7policies", "v1")

	task, err := l7policies.DeleteRule(client, pid, rid).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *task)
}

func TestRuleReplace(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetRuleTestURL(pid, rid), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRuleRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TaskResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := l7policies.CreateRuleOpts{
		CompareType: l7policies.CompareTypeRegex,
		Invert:      false,
		Type:        l7policies.TypePath,
		Value:       "/images*",
		Tags:        []string{"test_tag"},
	}

	client := fake.ServiceTokenClient("l7policies", "v1")
	tasks, err := l7policies.ReplaceRule(client, pid, rid, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}
