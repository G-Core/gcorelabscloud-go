package testing

import (
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets/testing"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "port_id": "817c8a3d-bb67-4b88-a0d1-aec980318ff1",
      "name": "Reserved fixed ip 10.100.179.44",
      "created_at": "2020-09-14T14:45:30+0000",
      "updated_at": "2020-09-14T14:45:31+0000",
      "status": "DOWN",
      "fixed_ip_address": "10.100.179.44",
      "subnet_id": "747db04a-2aac-4fda-9492-d9b85a798c09",
      "creator_task_id": "30378aea-9343-4ff6-be38-9756094e05da",
      "task_id": null,
      "is_external": false,
      "is_vip": false,
      "reservation": {
        "status": "available",
        "resource_type": null,
        "resource_id": null
      },
      "region": "ED-10",
      "region_id": 3,
      "project_id": 1,
      "allowed_address_pairs": [
        {
          "ip_address": "192.168.123.0/24",
          "mac_address": "00:16:3e:f2:87:16"
        }
      ],
      "network_id": "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
      "network": {
        "name": "public",
        "id": "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
        "subnets": [
          "747db04a-2aac-4fda-9492-d9b85a798c09"
        ],
        "mtu": 1500,
		"created_at": "2020-09-14T14:45:30+0000",
		"updated_at": "2020-09-14T14:45:31+0000",
        "external": true,
        "default": true,
        "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
        "project_id": 1,
        "region_id": 3
      }
    }
  ]
}
`

const GetResponse = `
{
  "port_id": "817c8a3d-bb67-4b88-a0d1-aec980318ff1",
  "name": "Reserved fixed ip 10.100.179.44",
  "created_at": "2020-09-14T14:45:30+0000",
  "updated_at": "2020-09-14T14:45:31+0000",
  "status": "DOWN",
  "fixed_ip_address": "10.100.179.44",
  "subnet_id": "747db04a-2aac-4fda-9492-d9b85a798c09",
  "creator_task_id": "30378aea-9343-4ff6-be38-9756094e05da",
  "task_id": null,
  "is_external": false,
  "is_vip": false,
  "reservation": {
	"status": "available",
	"resource_type": null,
	"resource_id": null
  },
  "region": "ED-10",
  "region_id": 3,
  "project_id": 1,
  "allowed_address_pairs": [
	{
	  "ip_address": "192.168.123.0/24",
	  "mac_address": "00:16:3e:f2:87:16"
	}
  ],
  "network_id": "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
  "network": {
	"name": "public",
	"id": "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
	"subnets": [
	  "747db04a-2aac-4fda-9492-d9b85a798c09"
	],
	"mtu": 1500,
	"created_at": "2020-09-14T14:45:30+0000",
	"updated_at": "2020-09-14T14:45:31+0000",
	"external": true,
	"default": true,
	"task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
	"project_id": 1,
	"region_id": 3
  }
}
`

const CreateRequest = `
{
    "type": "ip_address",
	"network_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
	"ip_address": "192.168.1.2",
    "is_vip": true
}
`

const TaskResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DeviceResponse = `
{
  "count": 1,
  "results": [
    {
      "port_id": "351b0dd7-ca09-431c-be53-935db3785067",
      "ip_assignments": [
        {
          "ip_address": "192.168.123.20",
          "subnet_id": "b39792c3-3160-4356-912e-ba396c95cdcf",
          "subnet": {
            "region": "Luxembourg 1",
            "total_ips": 253,
            "network_id": "b30d0de7-bca2-4c83-9c57-9e645bd2cc92",
            "task_id": null,
            "dns_nameservers": [
              "8.8.8.8",
              "8.8.4.4"
            ],
            "region_id": 3,
            "name": "subnet_3",
            "host_routes": [],
            "id": "b39792c3-3160-4356-912e-ba396c95cdcf",
            "creator_task_id": "5cc890da-d031-4a23-ac31-625edfa22812",
            "gateway_ip": "192.168.13.1",
            "ip_version": 4,
            "updated_at": "2020-09-14T14:45:31+0000",
            "created_at": "2020-09-14T14:45:30+0000",
            "cidr": "192.168.13.0/24",
            "enable_dhcp": true,
            "project_id": 1,
            "available_ips": 250,
            "has_router": false
          }
        }
      ],
      "instance_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
      "instance_name": "Virtual Machine 1",
      "network": {
        "name": "public",
        "id": "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
        "subnets": [
          "f00624ab-41bc-4d54-a723-1673ce32d997",
          "41e0f698-4d39-483b-b77a-18eb070e4c09"
        ],
        "mtu": 1500,
		"updated_at": "2020-09-14T14:45:31+0000",
		"created_at": "2020-09-14T14:45:30+0000",
        "external": true,
        "default": true,
        "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
        "project_id": 1,
        "region_id": 1
      }
    }
  ]
}
`

const SwitchVIPRequest = `
{
  "is_vip": true
}
`

const PortsRequest = `
{
  "port_ids": [
    "351b0dd7-ca09-431c-be53-935db3785067",
    "bc688791-f1b0-44eb-97d4-07697294b1e1"
  ]
}
`

var (
	createdTimeString    = "2020-09-14T14:45:30+0000"
	updatedTimeString    = "2020-09-14T14:45:31+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
	updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
	updatedTime          = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}

	cidr, _ = gcorecloud.ParseCIDRString("192.168.13.0/24")

	networkTaskID = "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"

	ReservedFixedIP1 = reservedfixedips.ReservedFixedIP{
		PortID:         "817c8a3d-bb67-4b88-a0d1-aec980318ff1",
		Name:           "Reserved fixed ip 10.100.179.44",
		CreatedAt:      createdTime,
		UpdatedAt:      updatedTime,
		Status:         "DOWN",
		FixedIPAddress: net.ParseIP("10.100.179.44"),
		SubnetID:       "747db04a-2aac-4fda-9492-d9b85a798c09",
		CreatorTaskID:  "30378aea-9343-4ff6-be38-9756094e05da",
		TaskID:         nil,
		IsExternal:     false,
		IsVip:          false,
		Reservation: reservedfixedips.IPReservation{
			Status:       "available",
			ResourceType: nil,
			ResourceID:   nil,
		},
		Region:    "ED-10",
		RegionID:  3,
		ProjectID: 1,
		AllowedAddressPairs: []reservedfixedips.AllowedAddressPairs{
			{
				IPAddress:  "192.168.123.0/24",
				MacAddress: "00:16:3e:f2:87:16",
			},
		},
		NetworkID: "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
		Network: networks.Network{
			Name:      "public",
			ID:        "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
			Subnets:   []string{"747db04a-2aac-4fda-9492-d9b85a798c09"},
			MTU:       1500,
			CreatedAt: createdTime,
			UpdatedAt: &updatedTime,
			External:  true,
			Default:   true,
			TaskID:    &networkTaskID,
			ProjectID: 1,
			RegionID:  3,
		},
	}

	ExpectedIPsSlice = []reservedfixedips.ReservedFixedIP{ReservedFixedIP1}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	Device1 = reservedfixedips.Device{
		PortID: "351b0dd7-ca09-431c-be53-935db3785067",
		IPAssignments: []reservedfixedips.IPAssignment{
			{
				IPAddress: net.ParseIP("192.168.123.20"),
				SubnetID:  "b39792c3-3160-4356-912e-ba396c95cdcf",
				Subnet: subnets.Subnet{
					ID:            "b39792c3-3160-4356-912e-ba396c95cdcf",
					Name:          "subnet_3",
					IPVersion:     4,
					EnableDHCP:    true,
					CIDR:          *cidr,
					CreatedAt:     createdTime,
					UpdatedAt:     updatedTime,
					NetworkID:     "b30d0de7-bca2-4c83-9c57-9e645bd2cc92",
					TaskID:        "",
					CreatorTaskID: "5cc890da-d031-4a23-ac31-625edfa22812",
					Region:        "Luxembourg 1",
					ProjectID:     1,
					RegionID:      3,
					AvailableIps:  testing.IPDual("250"),
					TotalIps:      testing.IPDual("253"),
					HasRouter:     false,
					DNSNameservers: []net.IP{
						net.ParseIP("8.8.8.8"),
						net.ParseIP("8.8.4.4"),
					},
					HostRoutes: []subnets.HostRoute{},
					GatewayIP:  net.ParseIP("192.168.13.1"),
				},
			},
		},
		InstanceID:   "bc688791-f1b0-44eb-97d4-07697294b1e1",
		InstanceName: "Virtual Machine 1",
		Network: networks.Network{
			Name: "public",
			ID:   "eed97610-708d-43a5-a9a5-caebd2b7b4ee",
			Subnets: []string{
				"f00624ab-41bc-4d54-a723-1673ce32d997",
				"41e0f698-4d39-483b-b77a-18eb070e4c09",
			},
			MTU:       1500,
			CreatedAt: createdTime,
			UpdatedAt: &updatedTime,
			External:  true,
			Default:   true,
			TaskID:    &networkTaskID,
			ProjectID: 1,
			RegionID:  1,
		},
	}
	ExpectedDevicesSlice = []reservedfixedips.Device{Device1}
)
