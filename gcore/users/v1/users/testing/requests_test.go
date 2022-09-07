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

	opts := users.CreateOpts{
		Email:    "test@test.test",
		Password: "test",
	}

	client := fake.ServiceTokenClient("", "internal")
	user, err := users.Create(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, User1, *user)
}
