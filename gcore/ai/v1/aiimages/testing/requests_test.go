package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiimages"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListTestURLParams(version string, projectID int, regionID int) string {
	return fmt.Sprintf("/%s/ai/images/%d/%d", version, projectID, regionID)
}

func prepareListTestURL() string {
	return prepareListTestURLParams("v1", fake.ProjectID, fake.RegionID)
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

	client := fake.ServiceTokenClient("ai/images", "v1")
	count := 0

	opts := aiimages.AIImageListOpts{}

	err := aiimages.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := aiimages.ExtractAIImages(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, AIImage1, ct)
		require.Equal(t, ExpectedAIImageSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
