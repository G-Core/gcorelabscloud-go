package testing

import (
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"

	"github.com/stretchr/testify/require"
)

func TestCreateOptions(t *testing.T) {
	options := clusters.CreateOpts{
		Name:         "",
		FixedNetwork: "",
		FixedSubnet:  "",
		KeyPair:      "",
		Version:      "",
	}

	_, err := options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "KeyPair")
	require.Contains(t, err.Error(), "FixedNetwork")
	require.Contains(t, err.Error(), "FixedSubnet")
	require.Contains(t, err.Error(), "Version")
	require.Contains(t, err.Error(), "Pools")
	require.NotContains(t, err.Error(), "PodsIPPool")
	require.NotContains(t, err.Error(), "ServicesIPPool")

	options = clusters.CreateOpts{
		Name:           "cluster-1-pool-1-machine-template",
		FixedNetwork:   fixedNetwork,
		FixedSubnet:    fixedSubnet,
		PodsIPPool:     &gcorecloud.CIDR{IPNet: *ipPool},
		ServicesIPPool: &gcorecloud.CIDR{IPNet: *ipPool},
		KeyPair:        "keypair",
		Version:        "v1.26.7",
		Pools:          []pools.CreateOpts{},
	}

	_, err = options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "Pools")
	require.NotContains(t, err.Error(), "FixedNetwork")
	require.NotContains(t, err.Error(), "FixedSubnet")
	require.NotContains(t, err.Error(), "PodsIPPool")
	require.NotContains(t, err.Error(), "ServicesIPPool")
	require.NotContains(t, err.Error(), "KeyPair")
	require.NotContains(t, err.Error(), "Version")
}

func TestUpgradeOptions(t *testing.T) {
	options := clusters.UpgradeOpts{
		Version: "",
	}

	_, err := options.ToClusterUpgradeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Version")

	options = clusters.UpgradeOpts{
		Version: "v1.26.7",
	}

	_, err = options.ToClusterUpgradeMap()
	require.NoError(t, err)
}
