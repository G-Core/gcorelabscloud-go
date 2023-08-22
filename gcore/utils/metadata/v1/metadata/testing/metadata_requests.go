package testing

import (
	"fmt"
	metadataV1 "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
	"net/http"
	"net/url"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	metadataTesting "github.com/G-Core/gcorelabscloud-go/gcore/utils/testing"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type TestFunc func(t *testing.T)

type UrlFunc func(c *gcorecloud.ServiceClient, id string, args ...string)

func prepareTestParams(resourceName string, urlFunc func(c *gcorecloud.ServiceClient) string, extraParams ...string) (client *gcorecloud.ServiceClient, relativeUrl string) {
	version := "v1"
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

func BuildTestMetadataListAll(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeUrl := prepareTestParams(resourceName, nil, extraParams...)

		th.Mux.HandleFunc(relativeUrl, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, metadataTesting.MetadataListResponse)
			if err != nil {
				log.Error(err)
			}
		})

		actual, err := metadataV1.MetadataListAll(client, resourceID)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, metadataTesting.Metadata1, ct)
		require.Equal(t, metadataTesting.ExpectedMetadataList, actual)
	}
}

func BuildTestMetadataGet(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeUrl := prepareTestParams(resourceName, func(c *gcorecloud.ServiceClient) string {
			return metadata.MetadataItemURL(c, resourceID, metadataTesting.ResourceMetadataReadOnly.Key)
		}, extraParams...)

		th.Mux.HandleFunc(relativeUrl, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, metadataTesting.MetadataResponse)
			if err != nil {
				log.Error(err)
			}
		})

		actual, err := metadataV1.MetadataGet(client, resourceID, metadataTesting.ResourceMetadataReadOnly.Key).Extract()
		require.NoError(t, err)
		require.Equal(t, &metadataTesting.ResourceMetadataReadOnly, actual)
	}
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

		err := metadataV1.MetadataCreateOrUpdate(client, resourceID, map[string]string{
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

		err := metadataV1.MetadataReplace(client, resourceID, map[string]string{
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

		err := metadataV1.MetadataDelete(client, resourceID, metadataTesting.Metadata1.Key).ExtractErr()
		require.NoError(t, err)
	}
}
