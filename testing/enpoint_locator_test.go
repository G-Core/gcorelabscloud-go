package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/stretchr/testify/require"
)

func TestEndpointLocationWithoutRegionAndProject(t *testing.T) {
	baseEndpoint := "http://test.com"

	eo := gcorecloud.EndpointOpts{
		Type:    "test",
		Name:    "test",
		Region:  0,
		Project: 0,
		Version: "v1",
	}

	el := gcorecloud.DefaultEndpointLocator(baseEndpoint)

	url, err := el(eo)
	require.NoError(t, err)
	require.Equal(t, "http://test.com/v1/test///test", url)
}
