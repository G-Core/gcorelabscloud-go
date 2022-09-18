package testing

import (
	"fmt"
	"net"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"

	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"

	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

var ListResponse = fmt.Sprintf(`
{
  "count": 1,
  "results": [
    {
      "floating_ip_address": "172.24.4.34",
      "router_id": "11005a33-c5ac-4c96-ab6f-8f2827cc7da6",
      "port_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
      "status": "ACTIVE",
      "id": "c64e5db1-5f1f-43ec-a8d9-5090df85b82d",
      "fixed_ip_address": null,
      "description": "FLOAT FLOAT FLOAT",
      "instance": {
        "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
        "task_state": null,
        "instance_description": "Testing",
        "instance_name": "Testing",
        "status": "ACTIVE",
        "instance_created": "2019-07-11T06:58:48Z",
        "vm_state": "active",
	    "region": "ED-8",
    	"volumes": [
          {
            "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
            "delete_on_termination": false
          }
        ],
        "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
        "metadata": {
          "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
        },
        "project_id": 1,
        "region_id": 1
      },
      "updated_at": "2019-06-13T13:58:12+0000",
      "created_at": "2019-06-13T13:58:12+0000",
	  "region": "ED-8",
      "project_id": 1,
      "region_id": 1,
      "metadata": [%s]
    }
  ]
}`, MetadataResponse)

var GetResponse = fmt.Sprintf(`
{
  "floating_ip_address": "172.24.4.34",
  "router_id": "11005a33-c5ac-4c96-ab6f-8f2827cc7da6",
  "port_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "status": "ACTIVE",
  "id": "c64e5db1-5f1f-43ec-a8d9-5090df85b82d",
  "fixed_ip_address": null,
  "description": "FLOAT FLOAT FLOAT",
  "updated_at": "2019-06-13T13:58:12+0000",
  "created_at": "2019-06-13T13:58:12+0000",
  "project_id": 1,
  "region": "ED-8",
  "region_id": 1,
 "metadata": [%s]
}`, MetadataResponse)

const CreateRequest = `
{
  "port_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "fixed_ip_address": "192.168.10.15"
}
`

const AssignRequest = `
{
  "port_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "fixed_ip_address": "192.168.10.15"
}
`

const UnassignRequest = `
{
  "port_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
  "fixed_ip_address": "192.168.10.15"
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

var AssignResponse = GetResponse
var UnassignResponse = GetResponse

var (
	floatingIPUpdatedAtParsed, _ = time.Parse(gcorecloud.RFC3339Z, "2019-06-13T13:58:12+0000")
	floatingIPCreatedAtParsed, _ = time.Parse(gcorecloud.RFC3339Z, "2019-06-13T13:58:12+0000")
	floatingIPAddress            = net.ParseIP("172.24.4.34")

	floatingIPCreatedAt = gcorecloud.JSONRFC3339Z{Time: floatingIPCreatedAtParsed}
	floatingIPUpdatedAt = gcorecloud.JSONRFC3339Z{Time: floatingIPUpdatedAtParsed}

	instanceCreatedAt, _ = time.Parse(gcorecloud.RFC3339ZZ, "2019-07-11T06:58:48Z")

	floatingIP = instances.FloatingIP{
		FloatingIPAddress: floatingIPAddress,
		RouterID:          "11005a33-c5ac-4c96-ab6f-8f2827cc7da6",
		SubnetID:          "",
		Status:            "ACTIVE",
		ID:                "c64e5db1-5f1f-43ec-a8d9-5090df85b82d",
		PortID:            "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		DNSDomain:         "",
		DNSName:           "",
		FixedIPAddress:    nil,
		UpdatedAt:         &floatingIPUpdatedAt,
		CreatedAt:         floatingIPCreatedAt,
		CreatorTaskID:     nil,
		ProjectID:         1,
		RegionID:          1,
		Region:            "ED-8",
		Metadata:          []metadata.Metadata{ResourceMetadataReadOnly},
	}

	instance = instances.Instance{
		ID:          "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
		Name:        "Testing",
		Description: "Testing",
		CreatedAt:   gcorecloud.JSONRFC3339ZZ{Time: instanceCreatedAt},
		Status:      "ACTIVE",
		VMState:     "active",
		TaskState:   nil,
		Flavor:      flavors.Flavor{},
		Metadata: map[string]interface{}{
			"task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
		},
		Volumes: []instances.InstanceVolume{{
			ID:                  "28bfe198-a003-4283-8dca-ab5da4a71b62",
			DeleteOnTermination: false,
		}},
		Addresses:        nil,
		SecurityGroups:   nil,
		CreatorTaskID:    nil,
		TaskID:           nil,
		ProjectID:        1,
		RegionID:         1,
		Region:           "ED-8",
		AvailabilityZone: instances.DefaultAvailabilityZone,
	}

	floatingIPDetails = floatingips.FloatingIPDetail{
		FloatingIPAddress: floatingIPAddress,
		RouterID:          "11005a33-c5ac-4c96-ab6f-8f2827cc7da6",
		SubnetID:          "",
		Status:            "ACTIVE",
		ID:                "c64e5db1-5f1f-43ec-a8d9-5090df85b82d",
		PortID:            "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
		DNSDomain:         "",
		DNSName:           "",
		FixedIPAddress:    nil,
		UpdatedAt:         &floatingIPUpdatedAt,
		CreatedAt:         floatingIPCreatedAt,
		CreatorTaskID:     nil,
		ProjectID:         1,
		RegionID:          1,
		Region:            "ED-8",
		Instance:          instance,
		Metadata:          []metadata.Metadata{ResourceMetadataReadOnly},
	}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedFloatingIPSlice = []floatingips.FloatingIPDetail{floatingIPDetails}

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
