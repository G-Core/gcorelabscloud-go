package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/inferences"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func prepareListTestURLParams(projectID int) string {
	return fmt.Sprintf("/v3/inference/%d/deployments", projectID)
}

func prepareGetTestURLParams(projectID int, name string) string {
	return fmt.Sprintf("/v3/inference/%d/deployments/%s", projectID, name)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID)
}

func prepareGetTestURL(name string) string {
	return prepareGetTestURLParams(fake.ProjectID, name)
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

	client := fake.ServiceTokenClient("inferences", "v3")
	infs, err := inferences.ListAllInferenceDeployments(client)
	th.AssertNoErr(t, err)

	if len(infs) != 1 {
		t.Errorf("Expected 1 page, got %d", len(infs))
	}

	ct := infs[0]
	require.Equal(t, Inference1, ct)
	require.Equal(t, InferencesSlice, infs)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Inference1.Name)
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

	client := fake.ServiceTokenClient("inferences", "v3")

	ct, err := inferences.GetInferenceDeployment(client, Inference1.Name).Extract()

	require.NoError(t, err)
	require.Equal(t, Inference1, *ct)
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

		_, err := fmt.Fprint(w, TasksResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := inferences.CreateInferenceDeploymentOpts{
		Name:          "test-inf",
		Image:         "nginx:latest",
		ListeningPort: 8080,
		AuthEnabled:   false,
		Containers: []inferences.CreateContainerOpts{
			{
				RegionID: regionID,
				Scale: inferences.ContainerScale{
					CooldownPeriod: &cooldownPeriod,
					Max:            3,
					Min:            1,
					Triggers: inferences.ContainerScaleTrigger{
						Cpu: &inferences.ScaleTriggerThreshold{
							Threshold: 80,
						},
						Memory: &inferences.ScaleTriggerThreshold{
							Threshold: 70,
						},
					},
				},
			},
		},
		Timeout: &timeout,
		Envs: map[string]string{
			"DEBUG_MODE": "False",
			"KEY":        "12345",
		},
		FlavorName: "inference-16vcpu-232gib-1xh100-80gb",
		Command: []string{
			"nginx",
			"-g",
			"daemon off;",
		},
		CredentialsName: &credentialsName,
		Logging: &inferences.CreateLoggingOpts{
			DestinationRegionID: regionID,
			Enabled:             true,
			RetentionPolicy: inferences.LoggingRetentionPolicy{
				Period: &retentionPolicy,
			},
			TopicName: topicName,
		},
	}

	client := fake.ServiceTokenClient("inferences", "v3")
	tasks, err := inferences.CreateInferenceDeployment(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Inference1.Name), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, TasksResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("inferences", "v3")
	tasks, err := inferences.DeleteInferenceDeployment(client, Inference1.Name).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Inference1.Name)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, TasksResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("inferences", "v3")
	opts := inferences.UpdateInferenceDeploymentOpts{
		Image:         &image,
		ListeningPort: &listeningPort,
		AuthEnabled:   &enableAuth,
		Containers: []inferences.CreateContainerOpts{
			{
				RegionID: regionID,
				Scale: inferences.ContainerScale{
					CooldownPeriod: &cooldownPeriod,
					Max:            3,
					Min:            1,
					Triggers: inferences.ContainerScaleTrigger{
						Cpu: &inferences.ScaleTriggerThreshold{
							Threshold: 80,
						},
						Memory: &inferences.ScaleTriggerThreshold{
							Threshold: 70,
						},
					},
				},
			},
		},
		Timeout: &timeout,
		Envs: map[string]string{
			"DEBUG_MODE": "False",
			"KEY":        "12345",
		},
		FlavorName: &flavorName,
		Command: []string{
			"nginx",
			"-g",
			"daemon off;",
		},
		CredentialsName: &credentialsName,
		Logging: &inferences.CreateLoggingOpts{
			DestinationRegionID: regionID,
			Enabled:             true,
			RetentionPolicy: inferences.LoggingRetentionPolicy{
				Period: &retentionPolicy,
			},
			TopicName: topicName,
		},
	}

	tasks, err := inferences.UpdateInferenceDeployment(client, Inference1.Name, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}
