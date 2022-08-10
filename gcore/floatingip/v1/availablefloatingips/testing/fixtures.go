package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"net"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"

	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const MetadataResponse = `
{
  "key": "some_key",
  "value": "some_val",
  "read_only": false
}
`

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
	ResourceMetadataReadOnly = metadata.Metadata{
		Key:      "some_key",
		Value:    "some_val",
		ReadOnly: false,
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

	ExpectedFloatingIPSlice = []floatingips.FloatingIPDetail{floatingIPDetails}
)
