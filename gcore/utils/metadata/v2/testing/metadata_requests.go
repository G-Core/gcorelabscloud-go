package testing

import (
	"fmt"
	metadataV2 "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1"
	"net/http"
	"net/url"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	metadataTesting "github.com/G-Core/gcorelabscloud-go/gcore/utils/testing"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	"github.com/stretchr/testify/require"
)

type TestFunc func(t *testing.T)

type UrlFunc func(c *gcorecloud.ServiceClient, id string, args ...string)

func prepareTestParams(resourceName string, urlFunc func(c *gcorecloud.ServiceClient) string, extraParams ...string) (client *gcorecloud.ServiceClient, relativeUrl string) {
	version := "v2"
	if extraParams != nil {
		version = extraParams[0]
	}

	client = fake.ServiceTokenClient(resourceName, version)

	resourceUrl := ""
	if urlFunc == nil {
		resourceUrl = client.ResourceBaseURL()

	} else {
		resourceUrl = urlFunc(client)
	}

	parsedUrl, err := url.Parse(resourceUrl)

	if err != nil {
		panic(err)
	}

	relativeUrl = parsedUrl.Path
	return
}

func BuildTestMetadataCreate(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeUrl := prepareTestParams(resourceName, func(c *gcorecloud.ServiceClient) string {
			return metadata.MetadataURL(c, resourceID)
		}, extraParams...)

		th.Mux.HandleFunc(relativeUrl, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, metadataTesting.MetadataCreateRequest)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadataV2.MetadataCreateOrUpdate(client, resourceID, map[string]string{
			"test1": "test1",
			"test2": "test2",
		}).ExtractErr()
		require.NoError(t, err)
	}
}

func BuildTestMetadataUpdate(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeUrl := prepareTestParams(resourceName, func(c *gcorecloud.ServiceClient) string {
			return metadata.MetadataURL(c, resourceID)
		}, extraParams...)

		th.Mux.HandleFunc(relativeUrl, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, metadataTesting.MetadataCreateRequest)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadataV2.MetadataReplace(client, resourceID, map[string]string{
			"test1": "test1",
			"test2": "test2",
		}).ExtractErr()
		require.NoError(t, err)
	}
}

func BuildTestMetadataDelete(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeUrl := prepareTestParams(resourceName, func(c *gcorecloud.ServiceClient) string {
			return metadata.MetadataItemURL(c, resourceID, metadataTesting.ResourceMetadataReadOnly.Key)
		}, extraParams...)

		th.Mux.HandleFunc(relativeUrl, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Accept", "application/json")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadataV2.MetadataDelete(client, resourceID, metadataTesting.Metadata1.Key).ExtractErr()
		require.NoError(t, err)
	}
}
