package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
)

func TestCreateOpts(t *testing.T) {
	options := networks.CreateOpts{
		Name:         Network1.Name,
		Mtu:          1450,
		CreateRouter: true,
	}
	_, err := options.ToNetworkCreateMap()
	require.NoError(t, err)

	options = networks.CreateOpts{
		Name:         Network1.Name,
		Mtu:          1501,
		CreateRouter: true,
	}
	_, err = options.ToNetworkCreateMap()
	require.Error(t, err)

}
