package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	th "github.com/G-Core/gcorelabscloud-go/testhelper"
)

func prepareSectionTestURLParams(projectID int, regionID int, section string) string { // nolint
	return fmt.Sprintf("/v1/laas/%d/%d/%s", projectID, regionID, section)
}

func prepareActionTestURLParams(projectID int, regionID int, section, action string) string { // nolint
	return fmt.Sprintf("/v1/laas/%d/%d/%s/%s", projectID, regionID, section, action)
}

func prepareStatusTestURL() string {
	return prepareSectionTestURLParams(fake.ProjectID, fake.RegionID, "status")
}

func prepareUsersTestURL() string {
	return prepareSectionTestURLParams(fake.ProjectID, fake.RegionID, "users")
}

func prepareTopicsTestURL() string {
	return prepareSectionTestURLParams(fake.ProjectID, fake.RegionID, "topics")
}

func prepareTopicsDeleteTestURL(topicName string) string {
	return prepareActionTestURLParams(fake.ProjectID, fake.RegionID, "topics", topicName)
}

func prepareOpenSearchTestURL() string { // nolint
	return fmt.Sprintf("/v1/laas/%d/opensearch_hosts", fake.RegionID)
}

func prepareKafkaTestURL() string { // nolint
	return fmt.Sprintf("/v1/laas/%d/kafka_hosts", fake.RegionID)
}

func TestGetStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareStatusTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, getStatusResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	status, err := laas.GetStatus(client).Extract()
	require.NoError(t, err)
	require.Equal(t, expectedStatus, *status)
}

func TestUpdateStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareStatusTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, updateRequest)
		w.Header().Add("Content-Type", "application/json")

		_, err := fmt.Fprint(w, getStatusResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	opts := laas.UpdateStatusOpts{IsInitialized: true}
	status, err := laas.UpdateStatus(client, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, expectedStatus, *status)
}

func TestListKafkaHost(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareKafkaTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, getKafkaResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	hosts, err := laas.ListKafkaHosts(client).Extract()
	require.NoError(t, err)
	require.Equal(t, expectedKafkaHosts, *hosts)
}

func TestListOpensearchHost(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareOpenSearchTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, getOpensearchResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	hosts, err := laas.ListOpenSearchHosts(client).Extract()
	require.NoError(t, err)
	require.Equal(t, expectedOpensearchHosts, *hosts)
}

func TestListTopic(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareTopicsTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, listTopicResponse)
		if err != nil {
			log.Error(err)
		}
	})

	var count int
	client := fake.ServiceTokenClient("laas", "v1")
	err := laas.ListTopic(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := laas.ExtractTopics(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, topic, ct)
		require.Equal(t, expectedTopicSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListTopicAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareTopicsTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, listTopicResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	topics, err := laas.ListTopicAll(client)
	require.NoError(t, err)

	ct := topics[0]
	require.Equal(t, topic, ct)
	require.Equal(t, expectedTopicSlice, topics)
}

func TestCreateTopic(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareTopicsTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createTopicRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, topicResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	opts := laas.CreateTopicOpts{Name: "string"}
	topic1, err := laas.CreateTopic(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, topic, *topic1)
}

func TestDeleteTopic(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareTopicsDeleteTestURL(topicName), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.WriteHeader(http.StatusNoContent)
	})

	client := fake.ServiceTokenClient("laas", "v1")
	err := laas.DeleteTopic(client, topicName).ExtractErr()
	require.NoError(t, err)
}

func TestRegenerateUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUsersTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, regenerateUserResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("laas", "v1")
	creds, err := laas.RegenerateUser(client).Extract()
	require.NoError(t, err)
	require.Equal(t, expectedCreds, *creds)
}
