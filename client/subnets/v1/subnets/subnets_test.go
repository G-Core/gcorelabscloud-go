package subnets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const testSubnetID = "e7944e55-f957-413d-aa56-fdc876543113"

const testSubnetResponse = `
{
  "id": "e7944e55-f957-413d-aa56-fdc876543113",
  "metadata": [
    {"key": "a", "value": "1", "read_only": false},
    {"key": "b", "value": "2", "read_only": false},
    {"key": "owner", "value": "gcore", "read_only": true}
  ]
}
`

type testRequest struct {
	method string
	body   string
}

// testSubnetServer serves GET and PATCH of the test subnet and records them;
// any other request, such as .../metadata, fails the test
func testSubnetServer(t *testing.T) func() []testRequest {
	var mu sync.Mutex
	var requests []testRequest

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		w.WriteHeader(http.StatusNotFound)
	})
	subnetURL := fmt.Sprintf("/v1/subnets/%d/%d/%s", fake.ProjectID, fake.RegionID, testSubnetID)
	th.Mux.HandleFunc(subnetURL, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("cannot read request body: %v", err)
		}
		mu.Lock()
		requests = append(requests, testRequest{method: r.Method, body: string(body)})
		mu.Unlock()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, testSubnetResponse)
	})

	return func() []testRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]testRequest(nil), requests...)
	}
}

func runMetadataCommand(args ...string) (string, error) {
	cc := func(*cli.Context) (*gcorecloud.ServiceClient, error) {
		return fake.ServiceTokenClient("subnets", "v1"), nil
	}
	var out bytes.Buffer
	app := &cli.App{
		Name:     "gcoreclient",
		Writer:   &out,
		Commands: []*cli.Command{newMetadataCommand(cc)},
		// return errors instead of calling os.Exit
		ExitErrHandler: func(*cli.Context, error) {},
	}
	err := app.Run(append([]string{"gcoreclient", "metadata"}, args...))
	return out.String(), err
}

func TestMetadataCommandsUseSubnetPatch(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantOutput string
		wantMethod []string
		wantPatch  string
	}{
		{
			name:       "list",
			args:       []string{"list", "-t", testSubnetID},
			wantOutput: `[{"key":"a","value":"1","read_only":false},{"key":"b","value":"2","read_only":false},{"key":"owner","value":"gcore","read_only":true}]`,
			wantMethod: []string{http.MethodGet},
		},
		{
			name:       "show",
			args:       []string{"show", "-t", "-m", "b", testSubnetID},
			wantOutput: `{"key":"b","value":"2","read_only":false}`,
			wantMethod: []string{http.MethodGet},
		},
		{
			name:       "create",
			args:       []string{"create", "-m", "c=3", testSubnetID},
			wantMethod: []string{http.MethodPatch},
			wantPatch:  `{"tags": {"c": "3"}}`,
		},
		{
			name:       "update",
			args:       []string{"update", "-m", "a=10", testSubnetID},
			wantMethod: []string{http.MethodPatch},
			wantPatch:  `{"tags": {"a": "10"}}`,
		},
		{
			name:       "replace",
			args:       []string{"replace", "-m", "a=10", "-m", "c=3", testSubnetID},
			wantMethod: []string{http.MethodGet, http.MethodPatch},
			wantPatch:  `{"tags": {"a": "10", "b": null, "c": "3"}}`,
		},
		{
			name:       "delete",
			args:       []string{"delete", "-m", "a", testSubnetID},
			wantMethod: []string{http.MethodPatch},
			wantPatch:  `{"tags": {"a": null}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			th.SetupHTTP()
			defer th.TeardownHTTP()
			requests := testSubnetServer(t)

			output, err := runMetadataCommand(tc.args...)
			require.NoError(t, err)

			got := requests()
			methods := make([]string, 0, len(got))
			for _, r := range got {
				methods = append(methods, r.method)
			}
			require.Equal(t, tc.wantMethod, methods)
			if tc.wantPatch != "" {
				require.JSONEq(t, tc.wantPatch, got[len(got)-1].body)
			}
			if tc.wantOutput != "" {
				require.JSONEq(t, tc.wantOutput, output)
			}
		})
	}
}

func TestMetadataShowMissingKey(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	testSubnetServer(t)

	_, err := runMetadataCommand("show", "-m", "missing", testSubnetID)

	require.EqualError(t, err, "Unable to find subnet tag with name missing")
}
