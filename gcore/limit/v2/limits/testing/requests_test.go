package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v1/limits"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	"github.com/stretchr/testify/require"
)

func prepareItemTestURL(limitID int) string {
	return fmt.Sprintf("/v2/limits_request/%d", limitID)
}

const limitRequestID = 1

func TestDelete(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareItemTestURL(limitRequestID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	client := fake.ServiceTokenClient("limits_request", "v2")
	err := limits.Delete(client, limitRequestID).ExtractErr()
	require.NoError(t, err)
}
