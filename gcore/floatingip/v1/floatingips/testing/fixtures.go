package testing

import (
	"net"
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/flavor/v1/flavors"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/floatingip/v1/floatingips"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

const ListResponse = `
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
      "region_id": 1
    }
  ]
}`

const GetResponse = `
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
  "region_id": 1
}
`

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

const AssignResponse = GetResponse
const UnassignResponse = GetResponse

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

	floatingIPDetails = floatingips.FloatingIPDetail{
		FloatingIP: &floatingIP,
		Instance:   instance,
	}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedFloatingIPSlice = []floatingips.FloatingIPDetail{floatingIPDetails}
)
