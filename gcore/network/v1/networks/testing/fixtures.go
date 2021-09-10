package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "creator_task_id": null,
      "region": "RegionOne",
      "name": "private",
      "mtu": 1450,
      "id": "e7944e55-f957-413d-aa56-fdc876543113",
      "updated_at": "2020-03-05T12:03:25+0000",
      "created_at": "2020-03-05T12:03:24+0000",
      "task_id": null,
      "region_id": 1,
      "shared": false,
      "subnets": [
        "3730b4d3-9337-4a60-a35e-7e1620aabe6f"
      ],
      "external": false,
      "project_id": 1
	}
  ]
}
`

const ListInstancePortResponse = `
{
  "count": 1,
  "results": [
    {
      "id": "8e009163-d526-4351-9266-93d9fd8fa8ef",
      "instance_id": "bfc7824b-31b6-4a28-a0c4-7df137139215",
      "instance_name": "instance_1"
    }
  ]
}
`

const GetResponse = `
{
  "creator_task_id": null,
  "region": "RegionOne",
  "name": "private",
  "mtu": 1450,
  "id": "e7944e55-f957-413d-aa56-fdc876543113",
  "updated_at": "2020-03-05T12:03:25+0000",
  "created_at": "2020-03-05T12:03:24+0000",
  "task_id": null,
  "region_id": 1,
  "shared": false,
  "subnets": [
    "3730b4d3-9337-4a60-a35e-7e1620aabe6f"
  ],
  "external": false,
  "project_id": 1
}
`

const CreateRequest = `
{
	"name": "private",
	"mtu": 1450,
	"create_router": true
}	
`

const UpdateRequest = `
{
	"name": "private"
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

var createdTimeString = "2020-03-05T12:03:24+0000"
var updatedTimeString = "2020-03-05T12:03:25+0000"
var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
var createdTime = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}

var (
	Network1 = networks.Network{
		Name: "private",
		ID:   "e7944e55-f957-413d-aa56-fdc876543113",
		Subnets: []string{
			"3730b4d3-9337-4a60-a35e-7e1620aabe6f",
		},
		MTU:       1450,
		CreatedAt: createdTime,
		UpdatedAt: &updatedTime,
		External:  false,
		Default:   false,
		Shared:    false,
		ProjectID: fake.ProjectID,
		RegionID:  fake.RegionID,
		Region:    "RegionOne",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedNetworkSlice = []networks.Network{Network1}

	InstancePort1 = networks.InstancePort{
		ID:           "8e009163-d526-4351-9266-93d9fd8fa8ef",
		InstanceID:   "bfc7824b-31b6-4a28-a0c4-7df137139215",
		InstanceName: "instance_1",
	}
	ExpectedInstancePortSlice = []networks.InstancePort{InstancePort1}
)
