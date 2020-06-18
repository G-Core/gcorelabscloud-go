package testing

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"

	"github.com/stretchr/testify/require"
)

func TestMarshallCreateStructure(t *testing.T) {
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
	require.JSONEq(t, CreateRequest, string(s))

}

func TestUpdateOpts(t *testing.T) {
	opts := subnets.UpdateOpts{}
	_, err := opts.ToSubnetUpdateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	opts = subnets.UpdateOpts{
		DNSNameservers: []net.IP{net.ParseIP("10.0.0.1")},
	}
	_, err = opts.ToSubnetUpdateMap()
	require.NoError(t, err)
}
