package testing

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
)

func TestCreateOpts(t *testing.T) {
	options := reservedfixedips.CreateOpts{
		Type: reservedfixedips.External,
	}

	_, err := options.ToReservedFixedIPCreateMap()
	require.NoError(t, err)

	options = reservedfixedips.CreateOpts{
		Type: reservedfixedips.Subnet,
	}

	_, err = options.ToReservedFixedIPCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "SubnetID")

	options.SubnetID = networkTaskID
	_, err = options.ToReservedFixedIPCreateMap()
	require.NoError(t, err)

	options = reservedfixedips.CreateOpts{
		Type: reservedfixedips.AnySubnet,
	}

	_, err = options.ToReservedFixedIPCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NetworkID")

	options.NetworkID = networkTaskID
	_, err = options.ToReservedFixedIPCreateMap()
	require.NoError(t, err)

	options = reservedfixedips.CreateOpts{
		Type:  reservedfixedips.IPAddress,
		IsVip: true,
	}

	_, err = options.ToReservedFixedIPCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NetworkID")
	require.Contains(t, err.Error(), "IPAddress")

	options.NetworkID = networkTaskID
	_, err = options.ToReservedFixedIPCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "IPAddress")

	options.IPAddress = net.ParseIP("192.168.1.2")
	mp, err := options.ToReservedFixedIPCreateMap()
	require.NoError(t, err)

	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, CreateRequest, string(s))
}

func TestPortsToShareVIPOpts(t *testing.T) {
	options := reservedfixedips.PortsToShareVIPOpts{
		PortIDs: []string{
			"351b0dd7-ca09-431c-be53-935db3785067",
			"bc688791-f1b0-44eb-97d4-07697294b1e1",
		},
	}

	mp, err := options.ToPortsToShareVIPOptsMap()
	require.NoError(t, err)

	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, PortsRequest, string(s))

	options = reservedfixedips.PortsToShareVIPOpts{
		PortIDs: []string{},
	}

	_, err = options.ToPortsToShareVIPOptsMap()
	require.NoError(t, err)

	options = reservedfixedips.PortsToShareVIPOpts{
		PortIDs: []string{
			"351b0dd7-ca09-431c-be53-935db3785067",
			"12322",
		},
	}

	_, err = options.ToPortsToShareVIPOptsMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "PortIDs")
}
