package testing

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/stretchr/testify/require"
)

func TestCreateOptsNoGW(t *testing.T) {
	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             true,
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
	}
	mp, err := options.ToSubnetCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, CreateRequestNoGW, string(s))
}

func TestCreateOptsGW(t *testing.T) {
	gw := net.IP{}
	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             true,
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
		GatewayIP:              &gw,
	}
	mp, err := options.ToSubnetCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, CreateRequestGW, string(s))
}

func TestUpdateOpts(t *testing.T) {
	opts := subnets.UpdateOpts{
		DNSNameservers: []net.IP{},
		HostRoutes:     []subnets.HostRoute{},
	}
	mp, err := opts.ToSubnetUpdateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, UpdateRequestNoData, string(s))
}
