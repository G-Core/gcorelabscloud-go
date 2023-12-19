package testing

import (
	"fmt"
	"math/big"
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

var ListResponse = fmt.Sprintf(`
{
  	"count": 1,
  	"results": [
		{
		  "id": "e7944e55-f957-413d-aa56-fdc876543113",
		  "name": "subnet",
		  "ip_version": 4,
		  "enable_dhcp": true,
		  "cidr": "192.168.10.0/24",
		  "created_at": "2020-03-05T12:03:24+0000",
		  "updated_at": "2020-03-05T12:03:25+0000",
		  "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		  "task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
		  "creator_task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
		  "region": "RegionOne",
		  "available_ips": 18446744073709551999,
		  "total_ips": 18446744073709552000,
		  "project_id": 1,
		  "region_id": 1,
		  "dns_nameservers": [
			"10.0.0.13"
		  ],
		  "gateway_ip" : "10.0.0.1",
		  "has_router": true,
		  "host_routes": [
			{
			  "destination": "10.0.3.0/24",
			  "nexthop": "10.0.0.13"
			}
		  ],
          "metadata": [%s]
    	}
  	]
}
`, MetadataResponse)

var GetResponse = fmt.Sprintf(`
{
  "id": "e7944e55-f957-413d-aa56-fdc876543113",
  "name": "subnet",
  "ip_version": 4,
  "enable_dhcp": true,
  "cidr": "192.168.10.0/24",
  "created_at": "2020-03-05T12:03:24+0000",
  "updated_at": "2020-03-05T12:03:25+0000",
  "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
  "creator_task_id": "50f53a35-42ed-40c4-82b2-5a37fb3e00bc",
  "region": "RegionOne",
  "project_id": 1,
  "region_id": 1,
  "available_ips": 18446744073709551999,
  "total_ips": 18446744073709552000,
  "dns_nameservers": [
	"10.0.0.13"
  ],
  "gateway_ip" : "10.0.0.1",
  "has_router": true,
  "host_routes": [
    {
      "destination": "10.0.3.0/24",
      "nexthop": "10.0.0.13"
    }
  ],
  "metadata": [%s]
}
`, MetadataResponse)

const CreateRequestNoGW = `
{
	"name": "subnet",
	"enable_dhcp": true,
	"cidr": "192.168.10.0/24",
	"network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
	"connect_to_network_router": true,
	"gateway_ip": null
}
`

const CreateRequestGW = `
{
	"name": "subnet",
	"enable_dhcp": true,
	"cidr": "192.168.10.0/24",
	"network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
	"connect_to_network_router": true
}
`

const UpdateRequestNoGW = `
{
 	"name": "subnet",
 	"gateway_ip": null,
    "dns_nameservers": null,
    "host_routes": null,
    "enable_dhcp": false
}
`

const UpdateRequestGW = `
{
 	"name": "subnet",
    "enable_dhcp": true,
    "dns_nameservers": null,
    "host_routes": null
}
`

const UpdateRequestNoData = `
{
    "enable_dhcp": false,
    "dns_nameservers": [],
    "host_routes": [],
    "gateway_ip": null
}
`
const CreateResponse = `
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
const MetadataResponse = `
{
  "key": "some_key",
  "value": "some_val",
  "read_only": false
}
`

const MetadataCreateRequest = `
{
"test1": "test1", 
"test2": "test2"
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

var createdTimeString = "2020-03-05T12:03:24+0000"
var updatedTimeString = "2020-03-05T12:03:25+0000"
var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
var cidr, _ = gcorecloud.ParseCIDRString("192.168.10.0/24")
var taskID = "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
var ip = net.ParseIP("10.0.0.13")
var gwip = net.ParseIP("10.0.0.1")
var routeCidr, _ = gcorecloud.ParseCIDRString("10.0.3.0/24")

func IPDual(i string) *big.Int {
	b := big.NewInt(0)
	b.SetString(i, 10)
	return b
}

var (
	Subnet1 = subnets.Subnet{
		ID:             "e7944e55-f957-413d-aa56-fdc876543113",
		Name:           "subnet",
		IPVersion:      4,
		EnableDHCP:     true,
		CIDR:           *cidr,
		CreatedAt:      createdTime,
		UpdatedAt:      updatedTime,
		NetworkID:      "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		TaskID:         taskID,
		CreatorTaskID:  taskID,
		Region:         "RegionOne",
		ProjectID:      fake.ProjectID,
		RegionID:       fake.RegionID,
		AvailableIps:   IPDual("18446744073709551999"),
		TotalIps:       IPDual("18446744073709552000"),
		HasRouter:      true,
		DNSNameservers: []net.IP{ip},
		GatewayIP:      gwip,
		HostRoutes: []subnets.HostRoute{
			{Destination: *routeCidr, NextHop: ip},
		},
		Metadata: []metadata.Metadata{ResourceMetadataReadOnly},
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedSubnetSlice = []subnets.Subnet{Subnet1}

	ResourceMetadata = map[string]interface{}{
		"some_key": "some_val",
	}

	ResourceMetadataReadOnly = metadata.Metadata{
		Key:      "some_key",
		Value:    "some_val",
		ReadOnly: false,
	}

	Metadata1 = metadata.Metadata{
		Key:      "cost-center",
		Value:    "Atlanta",
		ReadOnly: false,
	}
	Metadata2 = metadata.Metadata{
		Key:      "data-center",
		Value:    "A",
		ReadOnly: false,
	}
	ExpectedMetadataList = []metadata.Metadata{Metadata1, Metadata2}
)
