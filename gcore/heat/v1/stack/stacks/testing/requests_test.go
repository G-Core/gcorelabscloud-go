package testing

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks"

	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
	th "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/heat/%d/%d/stacks/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func prepareUpdateTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	url := prepareListTestURL()
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)

	client := fake.ServiceTokenClient("heat", "v1")

	createOpts := stacks.CreateOpts{
		Name:            "stackcreated",
		Timeout:         60,
		TemplateOpts:    template,
		DisableRollback: gcorecloud.Disabled,
	}
	actual, err := stacks.Create(client, createOpts).Extract()
	require.NoError(t, err)
	expected := CreateExpected
	require.Equal(t, expected, actual)
}

func TestCreateStackMissingRequiredInOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)

	client := fake.ServiceTokenClient("heat", "v1")

	createOpts := stacks.CreateOpts{
		DisableRollback: gcorecloud.Disabled,
	}
	r := stacks.Create(client, createOpts)
	th.AssertEquals(t, "Missing input for argument [Name]", r.Err.Error())
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	url := prepareListTestURL()
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")
	count := 0

	opts := stacks.ListOpts{}

	err := stacks.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		stackList, err := stacks.ExtractStacks(page)
		require.NoError(t, err)
		stack := stackList[0]
		require.Equal(t, StackList1, stack)
		require.Equal(t, ExpectedStackList1, stackList)
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
	url := prepareListTestURL()
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")
	opts := stacks.ListOpts{}
	stackList, err := stacks.ListAll(client, opts)
	require.NoError(t, err)
	stack := stackList[0]
	require.Equal(t, StackList1, stack)
	require.Equal(t, ExpectedStackList1, stackList)
}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Stack1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("heat", "v1")

	ct, err := stacks.Get(client, Stack1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Stack1, *ct)
	require.Equal(t, creationTime, ct.CreationTime)
	require.Equal(t, updatedTime, *ct.UpdatedTime)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareUpdateTestURL(Stack1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	template := new(stacks.Template)
	template.Bin = []byte(`
		{
			"heat_template_version": "2013-05-23",
			"description": "Simple template to test heat commands",
			"parameters": {
				"flavor": {
					"default": "m1.tiny",
					"type": "string"
				}
			}
		}`)

	updateOpts := &stacks.UpdateOpts{
		TemplateOpts: template,
	}

	client := fake.ServiceTokenClient("heat", "v1")

	err := stacks.Update(client, Stack1.ID, updateOpts).ExtractErr()
	require.NoError(t, err)

}

func TestUpdateStackNoTemplate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareUpdateTestURL(Stack1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	parameters := make(map[string]interface{})
	parameters["flavor"] = "m1.tiny"

	updateOpts := &stacks.UpdateOpts{
		Parameters: parameters,
	}
	expected := stacks.ErrTemplateRequired{}

	err := stacks.Update(fake.ServiceClient(), Stack1.ID, updateOpts).ExtractErr()
	th.AssertEquals(t, expected, err)
}

func TestUpdatePatch(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareUpdateTestURL(Stack1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	client := fake.ServiceTokenClient("heat", "v1")

	parameters := make(map[string]interface{})
	parameters["flavor"] = "m1.tiny"

	updateOpts := &stacks.UpdateOpts{
		Parameters: parameters,
	}

	err := stacks.UpdatePatch(client, Stack1.ID, updateOpts).ExtractErr()
	require.NoError(t, err)

}
