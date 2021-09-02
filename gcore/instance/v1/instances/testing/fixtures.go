package testing

import (
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/shopspring/decimal"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
      "task_state": null,
      "instance_description": "Testing",
      "instance_name": "Testing",
      "status": "ACTIVE",
      "instance_created": "2019-07-11T06:58:48Z",
      "vm_state": "active",
      "volumes": [
        {
          "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
          "delete_on_termination": false
        }
      ],
      "security_groups": [
        {
          "name": "default"
        }
      ],
      "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
      "task_id": "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02",
      "creator_task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
      "addresses": {
        "net1": [
          {
            "type": "fixed",
            "addr": "10.0.0.17"
          },
          {
            "type": "floating",
            "addr": "92.38.157.215"
          }
        ],
        "net2": [
          {
            "type": "fixed",
            "addr": "192.168.68.68"
          }
        ]
      },
      "metadata": {
        "os_distro": "centos",
        "os_version": "1711-x64",
        "image_name": "cirros-0.3.5-x86_64-disk",
        "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
        "snapshot_name": "test_snapshot",
        "snapshot_id": "c286cd13-fba9-4302-9cdb-4351a05a56ea",
        "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
      },
      "flavor": {
        "flavor_name": "g1s-shared-1-0.5",
        "disk": 0,
        "flavor_id": "g1s-shared-1-0.5",
        "vcpus": 1,
        "ram": 512
      },
      "project_id": 1,
      "region_id": 1,
	  "region": "RegionOne"	
    }
  ]
}
`

const GetResponse = `
{
  "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
  "task_state": null,
  "instance_description": "Testing",
  "instance_name": "Testing",
  "status": "ACTIVE",
  "instance_created": "2019-07-11T06:58:48Z",
  "vm_state": "active",
  "volumes": [
    {
      "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
      "delete_on_termination": false
    }
  ],
  "security_groups": [
    {
      "name": "default"
    }
  ],
  "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
  "task_id": "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02",
  "creator_task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
  "addresses": {
    "net1": [
      {
        "type": "fixed",
        "addr": "10.0.0.17"
      },
      {
        "type": "floating",
        "addr": "92.38.157.215"
      }
    ],
    "net2": [
      {
        "type": "fixed",
        "addr": "192.168.68.68"
      }
    ]
  },
  "metadata": {
    "os_distro": "centos",
    "os_version": "1711-x64",
    "image_name": "cirros-0.3.5-x86_64-disk",
    "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
    "snapshot_name": "test_snapshot",
    "snapshot_id": "c286cd13-fba9-4302-9cdb-4351a05a56ea",
    "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
  },
  "flavor": {
    "flavor_name": "g1s-shared-1-0.5",
    "disk": 0,
    "flavor_id": "g1s-shared-1-0.5",
    "vcpus": 1,
    "ram": 512
  },
  "project_id": 1,
  "region_id": 1,
  "region": "RegionOne"	
}
`

const InterfacesResponse = `
{
  "count": 1,
  "results": [
    {
      "port_id": "1f0ca628-a73b-42c0-bdac-7b10d023e097",
      "mac_address": "00:16:3e:f2:87:16",
      "network_id": "bc688791-f1b0-44eb-97d4-07697294b1e1",
	  "port_security_enabled": true,
      "ip_assignments": [
        {
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
        "subnets": [
          {
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
      "floatingip_details": [
        {
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
  ]
}
`

const CreateRequest = `
{
  "flavor": "g1-standard-1-2",
  "interfaces": [
	{
	  "floating_ip": {
		"existing_floating_id": "127.0.0.1",
		"source": "existing"
	  },
	  "network_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
	  "subnet_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
	  "type": "subnet"
	}
  ],
  "keypair_name": "keypair",
  "names": [
	"name"
  ],
  "password": "password",
  "security_groups": [
    {
      "id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c"
    }
  ],
  "user_data": "",
  "username": "username",
  "volumes": [
	{
	  "boot_index": 0,
	  "name": "name",
	  "size": 10,
	  "source": "new-volume",
	  "type_name": "standard"
	}
  ]
}
`

const SecurityGroupsListResponse = `
{
  "count": 1,
  "results": [
    {
      "name": "Test",
      "id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c"
    }
  ]
}
`

const AssignSecurityGroupsRequest = `
{
  "name": "Test"
}
`

const UnAssignSecurityGroupsRequest = `
{
  "name": "Test"
}
`

const AttachInterfaceRequest = `
{
  "type": "subnet",
  "subnet_id": "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
}
`

const DetachInterfaceRequest = `
{
  "ip_address": "192.168.0.23",
  "port_id": "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
}
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ResizeResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const MetadataListResponse = `
{
  "count": 2,
  "results": [
    {
      "key": "cost-center",
      "value": "Atlanta",
      "read_only": false
    },
    {
      "key": "data-center",
      "value": "A",
      "read_only": false
    }
  ]
}
`

const MetadataResponse = `
{
  "key": "cost-center",
  "value": "Atlanta",
  "read_only": false
}
`

const MetadataCreateRequest = `
{
"test1": "test1", 
"test2": "test2"
}
`

const InstanceConsoleResponse = `
{
    "remote_console": 

    {
        "url": "https://92.38.157.200:6082/spice_auto.html?token=369b3c64-541b-4cd7-bc47-a82cdb98228f",
        "type": "spice-html5",
        "protocol": "spice"
    }
}
`

const ListInstanceMetricsRequest = `
{
    "time_interval": 4,
    "time_unit": "day"
}
`

const ListInstanceMetricsResponse = `
{
  "count": 1,
  "results": [
    {
      "cpu_util": 8,
      "disks": [
        {
          "disk_Bps_read": 16384,
          "disk_Bps_write": 86016,
          "disk_iops_read": 3,
          "disk_iops_write": 12,
          "disk_name": "sda"
        }
      ],
      "memory_util": 33.28411162695459,
      "network_Bps_egress": 102,
      "network_Bps_ingress": 102,
      "network_pps_egress": 0.7,
      "network_pps_ingress": 0.7,
      "time": "2020-07-07T12:57:00Z"
    }
  ]
}
`

const ListAvailableFlavorsResponse = `
{
  "count": 1,
  "results": [
    {
      "currency_code": "USD",
      "flavor_id": "g1-gpu-1-2-1",
      "flavor_name": "g1-gpu-1-2-1",
      "hardware_description": {
        "gpu": "1x NVIDIA 11GB"
      },
      "price_per_hour": 1,
      "price_per_month": 720,
      "price_status": "show",
      "ram": 2048,
      "vcpus": 1
    }
  ]
}
`

var (
	ip1                 = net.ParseIP("10.0.0.17")
	ip2                 = net.ParseIP("92.38.157.215")
	ip3                 = net.ParseIP("192.168.68.68")
	tm, _               = time.Parse(gcorecloud.RFC3339ZZ, "2019-07-11T06:58:48Z")
	createdTime         = gcorecloud.JSONRFC3339ZZ{Time: tm}
	instanceID          = "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
	instanceName        = "Testing"
	instanceDescription = "Testing"
	taskID              = "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02"
	creatorTaskID       = "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"

	Instance1 = instances.Instance{
		ID:               instanceID,
		Name:             instanceName,
		Description:      instanceDescription,
		CreatedAt:        createdTime,
		Status:           "ACTIVE",
		VMState:          "active",
		TaskState:        nil,
		AvailabilityZone: instances.DefaultAvailabilityZone,
		Flavor: flavors.Flavor{
			FlavorID:   "g1s-shared-1-0.5",
			FlavorName: "g1s-shared-1-0.5",
			RAM:        512,
			VCPUS:      1,
		},
		Metadata: map[string]interface{}{
			"os_distro":     "centos",
			"os_version":    "1711-x64",
			"image_name":    "cirros-0.3.5-x86_64-disk",
			"image_id":      "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
			"snapshot_name": "test_snapshot",
			"snapshot_id":   "c286cd13-fba9-4302-9cdb-4351a05a56ea",
			"task_id":       "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
		},
		Volumes: []instances.InstanceVolume{{
			ID:                  "28bfe198-a003-4283-8dca-ab5da4a71b62",
			DeleteOnTermination: false,
		}},
		Addresses: map[string][]instances.InstanceAddress{
			"net1": {{
				Type:    types.AddressTypeFixed,
				Address: ip1,
			},
				{
					Type:    types.AddressTypeFloating,
					Address: ip2,
				},
			},
			"net2": {{
				Type:    types.AddressTypeFixed,
				Address: ip3,
			}},
		},
		SecurityGroups: []gcorecloud.ItemName{{
			Name: "default",
		}},
		CreatorTaskID: &creatorTaskID,
		TaskID:        &taskID,
		ProjectID:     1,
		RegionID:      1,
		Region:        "RegionOne",
	}
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
	MetricsTime, _                 = time.Parse(gcorecloud.RFC3339ZZ, "2020-07-07T12:57:00Z")
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
	InstanceInterface1 = instances.Interface{
		PortID:              PortID,
		MacAddress:          *PortMac,
		PortSecurityEnabled: true,
		NetworkID:           "bc688791-f1b0-44eb-97d4-07697294b1e1",
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
	ExpectedInstancesSlice          = []instances.Instance{Instance1}
	ExpectedInstanceInterfacesSlice = []instances.Interface{InstanceInterface1}
	ExpectedSecurityGroupsSlice     = []gcorecloud.ItemIDName{SecurityGroup1}
	Tasks1                          = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
	Metadata1 = instances.Metadata{
		Key:      "cost-center",
		Value:    "Atlanta",
		ReadOnly: false,
	}
	Metadata2 = instances.Metadata{
		Key:      "data-center",
		Value:    "A",
		ReadOnly: false,
	}
	ExpectedMetadataList = []instances.Metadata{Metadata1, Metadata2}
	Console              = instances.RemoteConsole{
		URL:      "https://92.38.157.200:6082/spice_auto.html?token=369b3c64-541b-4cd7-bc47-a82cdb98228f",
		Type:     "spice-html5",
		Protocol: "spice",
	}
	Metrics = instances.InstanceMetrics{
		CPUUtil: 8,
		Disks: []instances.DiskMetrics{
			{
				BpsRead:   16384,
				BpsWrite:  86016,
				IOPSRead:  3,
				IOPSWrite: 12,
				Name:      "sda",
			},
		},
		MemoryUtil:        33.28411162695459,
		NetworkBPSEgress:  102,
		NetworkBPSIngress: 102,
		NetworkPPSEgress:  0.7,
		NetworkPPSIngress: 0.7,
		Time:              gcorecloud.JSONRFC3339ZZ{Time: MetricsTime},
	}
	priceStatus     = "show"
	currencyCode, _ = gcorecloud.ParseCurrency("USD")
	pricePerHour    = decimal.NewFromInt(1)
	pricePerMonth   = decimal.NewFromInt(720)
	AvailableFlavor = flavors.Flavor{
		FlavorID:      "g1-gpu-1-2-1",
		FlavorName:    "g1-gpu-1-2-1",
		PriceStatus:   &priceStatus,
		CurrencyCode:  currencyCode,
		PricePerHour:  &pricePerHour,
		PricePerMonth: &pricePerMonth,
		RAM:           2048,
		VCPUS:         1,
	}
)
