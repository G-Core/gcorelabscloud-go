package testing

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListTestURLParams(projectID int, regionID int) string {
	return fmt.Sprintf("/v1/subnets/%d/%d", projectID, regionID)
}

func prepareGetTestURLParams(projectID int, regionID int, id string) string {
	return fmt.Sprintf("/v1/subnets/%d/%d/%s", projectID, regionID, id)
}

func prepareListTestURL() string {
	return prepareListTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return prepareGetTestURLParams(fake.ProjectID, fake.RegionID, id)
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

	client := fake.ServiceTokenClient("subnets", "v1")
	count := 0

	opts := subnets.ListOpts{
		NetworkID: Subnet1.NetworkID,
	}

	err := subnets.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := subnets.ExtractSubnets(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Subnet1, ct)
		require.Equal(t, ExpectedSubnetSlice, actual)
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

	client := fake.ServiceTokenClient("subnets", "v1")

	opts := subnets.ListOpts{
		NetworkID: Subnet1.NetworkID,
	}

	results, err := subnets.ListAll(client, opts)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, Subnet1, ct)
	require.Equal(t, ExpectedSubnetSlice, results)

}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

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

	client := fake.ServiceTokenClient("subnets", "v1")

	ct, err := subnets.Get(client, Subnet1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)

}

func TestCreateNoGW(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequestNoGW)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             boolPtr(true),
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
	}

	client := fake.ServiceTokenClient("subnets", "v1")
	tasks, err := subnets.Create(client, options, nil).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestCreateGW(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequestGW)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	gw := net.IP{}
	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             boolPtr(true),
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
		GatewayIP:              &gw,
	}

	client := fake.ServiceTokenClient("subnets", "v1")
	tasks, err := subnets.Create(client, options, nil).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetTestURL(Subnet1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, DeleteResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")
	tasks, err := subnets.Delete(client, Subnet1.ID, nil).Extract()
	require.NoError(t, err)
	require.Equal(t, Tasks1, *tasks)

}

func TestUpdateNoGW(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequestNoGW)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")

	opts := subnets.UpdateOpts{
		Name: Subnet1.Name,
	}

	ct, err := subnets.Update(client, Subnet1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
	require.Equal(t, Subnet1.Name, ct.Name)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)

}

func TestUpdateGW(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequestGW)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")

	gw := net.IP{}
	opts := subnets.UpdateOpts{
		Name:       Subnet1.Name,
		GatewayIP:  &gw,
		EnableDHCP: boolPtr(true),
	}

	ct, err := subnets.Update(client, Subnet1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
	require.Equal(t, Subnet1.Name, ct.Name)
	require.Equal(t, createdTime, ct.CreatedAt)
	require.Equal(t, updatedTime, ct.UpdatedAt)
}

func TestUpdateTags(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Subnet1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequestTags)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("subnets", "v1")

	value := "value1"
	opts := subnets.UpdateOpts{
		DNSNameservers: []net.IP{},
		HostRoutes:     []subnets.HostRoute{},
		Tags: map[string]*string{
			"key1": &value,
			"key2": nil,
		},
	}

	ct, err := subnets.Update(client, Subnet1.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, Subnet1, *ct)
}

// requestLog records the methods of the requests that a test server gets
type requestLog struct {
	mu      sync.Mutex
	methods []string
}

func (l *requestLog) add(method string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.methods = append(l.methods, method)
}

func (l *requestLog) list() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.methods...)
}

// closeTracker counts response bodies that the client did not close; an open
// body keeps its connection busy
type closeTracker struct {
	mu   sync.Mutex
	open int
}

type trackedBody struct {
	io.ReadCloser
	tracker *closeTracker
	once    sync.Once
}

func (b *trackedBody) Close() error {
	b.once.Do(func() {
		b.tracker.mu.Lock()
		defer b.tracker.mu.Unlock()
		b.tracker.open--
	})
	return b.ReadCloser.Close()
}

func (c *closeTracker) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.open++
	resp.Body = &trackedBody{ReadCloser: resp.Body, tracker: c}
	return resp, nil
}

// newTagsClient returns a client that fails the test if a response body
// stays open
func newTagsClient(t *testing.T) *gcorecloud.ServiceClient {
	client := fake.ServiceTokenClient("subnets", "v1")
	tracker := &closeTracker{}
	client.HTTPClient.Transport = tracker
	t.Cleanup(func() {
		tracker.mu.Lock()
		defer tracker.mu.Unlock()
		require.Zero(t, tracker.open, "response bodies left open")
	})
	return client
}

// handleSubnetTags serves GET and PATCH of Subnet1 and checks that the PATCH
// body is patchRequest; other paths, such as .../metadata, get 404
func handleSubnetTags(t *testing.T, patchRequest string) *requestLog {
	requests := &requestLog{}
	th.Mux.HandleFunc(prepareGetTestURL(Subnet1.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		requests.add(r.Method)
		switch r.Method {
		case http.MethodGet:
		case http.MethodPatch:
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestJSONRequest(t, r, patchRequest)
		default:
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, MetadataGetResponse)
		if err != nil {
			log.Error(err)
		}
	})
	return requests
}

func TestMetadataListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	requests := handleSubnetTags(t, "")

	client := newTagsClient(t)
	actual, err := subnets.MetadataListAll(client, Subnet1.ID)

	require.NoError(t, err)
	require.Equal(t, ExpectedMetadataList, actual)
	require.Equal(t, []string{http.MethodGet}, requests.list())
}

func TestMetadataGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	requests := handleSubnetTags(t, "")

	client := newTagsClient(t)
	actual, err := subnets.MetadataGet(client, Subnet1.ID, Metadata2.Key).Extract()

	require.NoError(t, err)
	require.Equal(t, &Metadata2, actual)
	require.Equal(t, []string{http.MethodGet}, requests.list())
}

func TestMetadataGetMissingKey(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	handleSubnetTags(t, "")

	client := newTagsClient(t)
	_, err := subnets.MetadataGet(client, Subnet1.ID, "missing").Extract()

	require.IsType(t, gcorecloud.ErrResourceNotFound{}, err)
}

func TestMetadataCreateOrUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	requests := handleSubnetTags(t, MetadataCreateRequest)

	client := newTagsClient(t)
	err := subnets.MetadataCreateOrUpdate(client, Subnet1.ID, map[string]string{
		"test1": "test1",
		"test2": "test2",
	}).ExtractErr()

	require.NoError(t, err)
	require.Equal(t, []string{http.MethodPatch}, requests.list())
}

func TestMetadataReplace(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	requests := handleSubnetTags(t, MetadataReplaceRequest)

	client := newTagsClient(t)
	err := subnets.MetadataReplace(client, Subnet1.ID, map[string]string{
		"cost-center": "Boston",
		"test1":       "test1",
	}).ExtractErr()

	require.NoError(t, err)
	require.Equal(t, []string{http.MethodGet, http.MethodPatch}, requests.list())
}

func TestMetadataDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	requests := handleSubnetTags(t, MetadataDeleteRequest)

	client := newTagsClient(t)
	err := subnets.MetadataDelete(client, Subnet1.ID, Metadata1.Key).ExtractErr()

	require.NoError(t, err)
	require.Equal(t, []string{http.MethodPatch}, requests.list())
}
