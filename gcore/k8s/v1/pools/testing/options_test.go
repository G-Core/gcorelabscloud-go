package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"

	"github.com/stretchr/testify/require"
)

func TestUpdateOpts(t *testing.T) {

	options := pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err := options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")

	options = pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)

	options = pools.UpdateOpts{}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "Name")

}

func TestCreateOpts(t *testing.T) {

	options := pools.CreateOpts{
		Name:             "",
		FlavorID:         "",
		NodeCount:        0,
		DockerVolumeSize: 0,
		MinNodeCount:     5,
		MaxNodeCount:     3,
	}

	_, err := options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "FlavorID")
	require.Contains(t, err.Error(), "NodeCount")

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        5,
		DockerVolumeSize: 10,
		MinNodeCount:     4,
		MaxNodeCount:     3,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        5,
		DockerVolumeSize: 10,
		MinNodeCount:     6,
		MaxNodeCount:     8,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "NodeCount")

}
