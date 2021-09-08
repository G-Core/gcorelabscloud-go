package testing

import (
	"net"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const EnableResponse = `
{
	"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
	"mac_address": "00:16:3e:f2:87:16",
	"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
	"port_security_enabled": true,
	"ip_assignments": [{
			"ip_address": "192.168.123.20",
			"subnet_id": "351b0dd7-ca09-431c-be53-935db3785067"
		},
		{
			"ip_address": "192.168.120.16",
			"subnet_id": "0a641ef8-62dc-4146-91e5-6ab4b464df6d"
		}
	],
	"network_details": {
		"external": false,
		"region_id": 1,
		"id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
		"mtu": 1450,
		"task_id": "a10dd116-07f5-4225-abb7-f42da5cb78fc",
		"creator_task_id": "a10dd116-07f5-4225-abb7-f42da5cb78fb",
		"name": "test network",
		"updated_at": "2020-02-26T08:44:44+0000",
		"project_id": 1,
		"shared": false,
		"subnets": [{
				"created_at": "2019-07-18T12:07:00+0000",
				"enable_dhcp": true,
				"id": "351b0dd7-ca09-431c-be53-935db3785067",
				"cidr": "192.168.123.0/24",
				"updated_at": "2019-07-22T10:55:45+0000",
				"ip_version": 4,
				"name": "test 2",
				"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
				"project_id": 1,
				"region_id": 1
			},
			{
				"created_at": "2019-07-22T15:15:05+0000",
				"enable_dhcp": true,
				"id": "0a641ef8-62dc-4146-91e5-6ab4b464df6d",
				"cidr": "192.168.120.0/24",
				"updated_at": "2019-07-22T15:15:05+0000",
				"ip_version": 4,
				"name": "string",
				"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
				"project_id": 1,
				"region_id": 1
			}
		],
		"created_at": "2020-02-26T08:44:08+0000",
		"region": "ED-8"
	},
	"floatingip_details": [{
			"region_id": 1,
			"id": "f32fe70c-f2ce-492e-858a-621bdc234885",
			"status": "ACTIVE",
			"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			"updated_at": "2020-02-26T08:47:23+0000",
			"project_id": 1,
			"fixed_ip_address": "192.168.123.20",
			"floating_ip_address": "5.188.135.29",
			"creator_task_id": "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b",
			"created_at": "2020-02-26T08:47:19+0000",
			"region": "ED-8",
			"router_id": "bf231ab1-769f-44db-bcb1-7f4199a1cba8"
		},
		{
			"region_id": 1,
			"id": "f32fe70c-f2ce-492e-858a-621bdc234441",
			"status": "ACTIVE",
			"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			"updated_at": "2020-02-26T08:47:23+0000",
			"project_id": 1,
			"fixed_ip_address": "192.168.120.16",
			"floating_ip_address": "5.188.135.30",
			"creator_task_id": "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b",
			"created_at": "2020-02-26T08:47:18+0000",
			"region": "ED-8",
			"router_id": "bf231ab1-769f-44db-bcb1-7f4199a1cba8"
		}
	]
}
`

const DisableResponse = `
{
	"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
	"mac_address": "00:16:3e:f2:87:16",
	"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
	"port_security_enabled": false,
	"ip_assignments": [{
			"ip_address": "192.168.123.20",
			"subnet_id": "351b0dd7-ca09-431c-be53-935db3785067"
		},
		{
			"ip_address": "192.168.120.16",
			"subnet_id": "0a641ef8-62dc-4146-91e5-6ab4b464df6d"
		}
	],
	"network_details": {
		"external": false,
		"region_id": 1,
		"id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
		"mtu": 1450,
		"task_id": "a10dd116-07f5-4225-abb7-f42da5cb78fc",
		"creator_task_id": "a10dd116-07f5-4225-abb7-f42da5cb78fb",
		"name": "test network",
		"updated_at": "2020-02-26T08:44:44+0000",
		"project_id": 1,
		"shared": false,
		"subnets": [{
				"created_at": "2019-07-18T12:07:00+0000",
				"enable_dhcp": true,
				"id": "351b0dd7-ca09-431c-be53-935db3785067",
				"cidr": "192.168.123.0/24",
				"updated_at": "2019-07-22T10:55:45+0000",
				"ip_version": 4,
				"name": "test 2",
				"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
				"project_id": 1,
				"region_id": 1
			},
			{
				"created_at": "2019-07-22T15:15:05+0000",
				"enable_dhcp": true,
				"id": "0a641ef8-62dc-4146-91e5-6ab4b464df6d",
				"cidr": "192.168.120.0/24",
				"updated_at": "2019-07-22T15:15:05+0000",
				"ip_version": 4,
				"name": "string",
				"network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
				"project_id": 1,
				"region_id": 1
			}
		],
		"created_at": "2020-02-26T08:44:08+0000",
		"region": "ED-8"
	},
	"floatingip_details": [{
			"region_id": 1,
			"id": "f32fe70c-f2ce-492e-858a-621bdc234885",
			"status": "ACTIVE",
			"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			"updated_at": "2020-02-26T08:47:23+0000",
			"project_id": 1,
			"fixed_ip_address": "192.168.123.20",
			"floating_ip_address": "5.188.135.29",
			"creator_task_id": "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b",
			"created_at": "2020-02-26T08:47:19+0000",
			"region": "ED-8",
			"router_id": "bf231ab1-769f-44db-bcb1-7f4199a1cba8"
		},
		{
			"region_id": 1,
			"id": "f32fe70c-f2ce-492e-858a-621bdc234441",
			"status": "ACTIVE",
			"port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			"updated_at": "2020-02-26T08:47:23+0000",
			"project_id": 1,
			"fixed_ip_address": "192.168.120.16",
			"floating_ip_address": "5.188.135.30",
			"creator_task_id": "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b",
			"created_at": "2020-02-26T08:47:18+0000",
			"region": "ED-8",
			"router_id": "bf231ab1-769f-44db-bcb1-7f4199a1cba8"
		}
	]
}
`

const allowedAddressPairsRequest = `
{
  "allowed_address_pairs": [
    {
      "ip_address": "192.168.123.20",
      "mac_address": "00:16:3e:f2:87:16"
    }
  ]
}
`

const allowedAddressPairsResponse = `
	{
  "port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
  "instance_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
  "network_id": "351b0dd7-ca09-431c-be53-935db3785067",
  "allowed_address_pairs": [
    {
      "ip_address": "192.168.123.20",
      "mac_address": "00:16:3e:f2:87:16"
    }
  ]
}
`

var (
	PortID                         = "1f0ca628-a73b-42c0-bdac-7b10d023e097"
	PortMac, _                     = gcorecloud.ParseMacString("00:16:3e:f2:87:16")
	PortIP1                        = net.ParseIP("192.168.123.20")
	PortIP2                        = net.ParseIP("192.168.120.16")
	PortNetworkUpdatedAt, _        = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:44:44+0000")
	PortNetworkCreatedAt, _        = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:44:08+0000")
	PortNetworkSubnet1CreatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2019-07-18T12:07:00+0000")
	PortNetworkSubnet1UpdatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2019-07-22T10:55:45+0000")
	PortNetworkSubnet2CreatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2019-07-22T15:15:05+0000")
	PortNetworkSubnet2UpdatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2019-07-22T15:15:05+0000")
	PortNetworkSubnet1Cidr, _      = gcorecloud.ParseCIDRString("192.168.123.0/24")
	PortNetworkSubnet2Cidr, _      = gcorecloud.ParseCIDRString("192.168.120.0/24")
	FloatingIP1                    = net.ParseIP("5.188.135.29")
	FixedIP1                       = net.ParseIP("192.168.123.20")
	FloatingIP2                    = net.ParseIP("5.188.135.30")
	FixedIP2                       = net.ParseIP("192.168.120.16")
	PortFloatingIP1CreatedAt, _    = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:47:19+0000")
	PortFloatingIP1UpdatedAt, _    = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:47:23+0000")
	PortFloatingIP1CreatorTaskID   = "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b"
	PortFloatingIP2CreatedAt, _    = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:47:18+0000")
	PortFloatingIP2UpdatedAt, _    = time.Parse(gcorecloud.RFC3339Z, "2020-02-26T08:47:23+0000")
	PortFloatingIP2CreatorTaskID   = "d1c1fd65-6eef-4e64-96cb-705cbbdbc90b"
	NetworkDetailsCreatorTask      = "a10dd116-07f5-4225-abb7-f42da5cb78fb"
	NetworkDetailsTask             = "a10dd116-07f5-4225-abb7-f42da5cb78fc"
	SecurityGroup1                 = gcorecloud.ItemIDName{
		ID:   "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
		Name: "Test",
	}
	instanceInterface = instances.Interface{
		PortID:     PortID,
		MacAddress: *PortMac,
		NetworkID:  "bc688791-f1b0-44eb-97d4-07697294b1e1",
		IPAssignments: []instances.PortIP{{
			IPAddress: PortIP1,
			SubnetID:  "351b0dd7-ca09-431c-be53-935db3785067",
		}, {
			IPAddress: PortIP2,
			SubnetID:  "0a641ef8-62dc-4146-91e5-6ab4b464df6d",
		}},
		NetworkDetails: instances.NetworkDetail{
			Mtu:           1450,
			UpdatedAt:     &gcorecloud.JSONRFC3339Z{Time: PortNetworkUpdatedAt},
			CreatedAt:     gcorecloud.JSONRFC3339Z{Time: PortNetworkCreatedAt},
			ID:            "bc688791-f1b0-44eb-97d4-07697294b1e1",
			External:      false,
			Default:       false,
			Shared:        false,
			Name:          "test network",
			TaskID:        &NetworkDetailsTask,
			CreatorTaskID: &NetworkDetailsCreatorTask,
			Subnets: []instances.Subnet{{
				ID:         "351b0dd7-ca09-431c-be53-935db3785067",
				Name:       "test 2",
				IPVersion:  gcorecloud.IPv4,
				EnableDHCP: true,
				Cidr:       *PortNetworkSubnet1Cidr,
				CreatedAt:  gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet1CreatedAt},
				UpdatedAt:  &gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet1UpdatedAt},
				NetworkID:  "bc688791-f1b0-44eb-97d4-07697294b1e1",
				ProjectID:  1,
				RegionID:   1,
				Region:     "",
			}, {
				ID:         "0a641ef8-62dc-4146-91e5-6ab4b464df6d",
				Name:       "string",
				IPVersion:  gcorecloud.IPv4,
				EnableDHCP: true,
				Cidr:       *PortNetworkSubnet2Cidr,
				CreatedAt:  gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet2CreatedAt},
				UpdatedAt:  &gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet2UpdatedAt},
				NetworkID:  "bc688791-f1b0-44eb-97d4-07697294b1e1",
				ProjectID:  1,
				RegionID:   1,
				Region:     "",
			}},
			ProjectID: 1,
			RegionID:  1,
			Region:    "ED-8",
		},
		FloatingIPDetails: []instances.FloatingIP{{
			FloatingIPAddress: FloatingIP1,
			FixedIPAddress:    FixedIP1,
			Status:            "ACTIVE",
			RouterID:          "bf231ab1-769f-44db-bcb1-7f4199a1cba8",
			ID:                "f32fe70c-f2ce-492e-858a-621bdc234885",
			PortID:            "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			CreatedAt:         gcorecloud.JSONRFC3339Z{Time: PortFloatingIP1CreatedAt},
			UpdatedAt:         &gcorecloud.JSONRFC3339Z{Time: PortFloatingIP1UpdatedAt},
			CreatorTaskID:     &PortFloatingIP1CreatorTaskID,
			ProjectID:         1,
			RegionID:          1,
			Region:            "ED-8",
		}, {
			FloatingIPAddress: FloatingIP2,
			FixedIPAddress:    FixedIP2,
			Status:            "ACTIVE",
			RouterID:          "bf231ab1-769f-44db-bcb1-7f4199a1cba8",
			ID:                "f32fe70c-f2ce-492e-858a-621bdc234441",
			PortID:            "1f0ca628-a73b-42c0-bdac-7b10d023e097",
			CreatedAt:         gcorecloud.JSONRFC3339Z{Time: PortFloatingIP2CreatedAt},
			UpdatedAt:         &gcorecloud.JSONRFC3339Z{Time: PortFloatingIP2UpdatedAt},
			CreatorTaskID:     &PortFloatingIP2CreatorTaskID,
			ProjectID:         1,
			RegionID:          1,
			Region:            "ED-8",
		}},
	}
	addrPairs1 = ports.InstancePort{
		NetworkID: "351b0dd7-ca09-431c-be53-935db3785067",
		AllowedAddressPairs: []reservedfixedips.AllowedAddressPairs{
			{
				IPAddress:  PortIP1,
				MacAddress: "00:16:3e:f2:87:16",
			},
		},
		InstanceID: "bc688791-f1b0-44eb-97d4-07697294b1e1",
		PortID:     "1f0ca628-a73b-42c0-bdac-7b10d023e097",
	}
)
