package testing

import (
	"encoding/json"
	"net"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/stretchr/testify/require"
)

func TestCreateOpts(t *testing.T) {
	snat := true
	options := routers.CreateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
		},
		Interfaces: []routers.Interface{
			{
				Type:     types.SubnetInterfaceType,
				SubnetID: "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
			},
		},
	}
	mp, err := options.ToRouterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NetworkID")

	options = routers.CreateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		},
		Interfaces: []routers.Interface{
			{
				Type:     "type",
				SubnetID: "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
			},
		},
	}
	mp, err = options.ToRouterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type")

	var gccidr gcorecloud.CIDR

	_, netIPNet, _ := net.ParseCIDR("10.0.3.0/24")
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask

	options = routers.CreateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		},
		Interfaces: []routers.Interface{
			{
				Type: "type",
			},
		},
		Routes: []subnets.HostRoute{
			{
				Destination: gccidr,
				NextHop:     net.ParseIP("10.0.0.13"),
			},
		},
	}
	mp, err = options.ToRouterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type")
	require.Contains(t, err.Error(), "SubnetID")

	_, netIPNet, _ = net.ParseCIDR("10.0.3.0/24")
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask

	options = routers.CreateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		},
		Interfaces: []routers.Interface{
			{
				Type:     types.SubnetInterfaceType,
				SubnetID: "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
			},
		},
		Routes: []subnets.HostRoute{
			{
				Destination: gccidr,
				NextHop:     net.ParseIP("10.0.0.13"),
			},
		},
	}
	mp, err = options.ToRouterCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, CreateRequest, string(s))
}

func TestUpdateOpts(t *testing.T) {
	snat := false
	opts := routers.UpdateOpts{
		Name: Router1.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       "type",
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		},
	}

	mp, err := opts.ToRouterUpdateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type")

	var gccidr gcorecloud.CIDR
	_, netIPNet, _ := net.ParseCIDR("10.0.4.0/24")
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask

	opts = routers.UpdateOpts{
		Name: Router2.Name,
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       types.ManualGateway,
			EnableSNat: &snat,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f2",
		},
		Routes: []subnets.HostRoute{
			{
				Destination: gccidr,
				NextHop:     net.ParseIP("10.0.0.14"),
			},
		},
	}
	mp, err = opts.ToRouterUpdateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, UpdateRequest, string(s))
}
