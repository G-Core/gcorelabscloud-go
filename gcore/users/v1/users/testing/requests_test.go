package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/users/v1/users"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareCreateTestURLParams() string {
	return fmt.Sprintf("/internal/users")
}

func prepareCreateApiTokenTestURLParams() string {
	return fmt.Sprintf("/internal/permanent_api_token")
}

func prepareUserAssignmentsTestURLParams() string {
	return fmt.Sprintf("/v1/users/assignments")
}

func TestCreateUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateTestURLParams(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	opts := users.CreateUserOpts{
		Email:    "test@test.test",
		Password: "test",
	}

	client := fake.ServiceTokenClient("", "internal")
	user, err := users.CreateUser(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, User1, *user)
}

func TestCreateApiToken(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateApiTokenTestURLParams(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateApiTokenRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CreateApiTokenResponse)
		if err != nil {
			log.Error(err)
		}
	})

	opts := users.CreateApiTokenOpts{
		Email:            "test@test.test",
		Password:         "test",
		TokenName:        "test",
		TokenDescription: "test description",
	}

	client := fake.ServiceTokenClient("", "internal")
	token, err := users.CreateApiToken(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, Token1, *token)
}

func TestUserAssignments(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUserAssignmentsTestURLParams(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UserAssignmentsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UserAssignmentsResponse)
		if err != nil {
			log.Error(err)
		}
	})

	clientID = 8
	opts := users.UserAssignmentOpts{
		ClientID: &clientID,
		UserID:   777,
		Role:     "ClientAdministrator",
	}

	client := fake.ServiceTokenClient("users", "v1")
	ua, err := users.AssignUser(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, UA1, *ua)
}
