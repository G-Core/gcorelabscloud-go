package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/pools"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/clusters"

	"github.com/stretchr/testify/require"
)

func TestResizeOpts(t *testing.T) {
	options := clusters.ResizeOpts{
		NodeCount:     0,
		NodesToRemove: nil,
		Pool:          "",
	}

	_, err := options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodeCount")

	options = clusters.ResizeOpts{
		NodeCount:     1,
		NodesToRemove: []string{"1"},
		Pool:          "",
	}

	_, err = options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodesToRemove")
	require.Contains(t, err.Error(), "Pool")

	options = clusters.ResizeOpts{
		NodeCount: 1,
		Pool:      "",
	}

	_, err = options.ToClusterResizeMap()
	require.NoError(t, err)

}

func TestCreateOptions(t *testing.T) {
	options := clusters.CreateOpts{
		Name:         Cluster1.Name,
		MasterCount:  1,
		KeyPair:      "",
		FixedNetwork: "",
		FixedSubnet:  "",
		Version:      "",
		Pools:        []pools.CreateOpts{},
	}

	_, err := options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "FixedNetwork")
	require.Contains(t, err.Error(), "FixedSubnet")
	require.Contains(t, err.Error(), "Pools")

	options = clusters.CreateOpts{
		Name:         Cluster1.Name,
		MasterCount:  4,
		KeyPair:      "",
		FixedNetwork: fixedNetwork,
		FixedSubnet:  fixedNetwork,
		Version:      "111",
	}

	_, err = options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Version")
	require.Contains(t, err.Error(), "Pools")

}
