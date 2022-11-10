package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const TaskResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var ListResponse = `
{
  "count": 1,
  "results": [
    {
      "instances": [
        {
          "instance_id": "6d14f194-6c1e-49b3-9fc7-50dc8401eb74",
          "instance_name": "test_ruslan_aa2"
        }
      ],
      "name": "example_server_group",
      "policy": "anti-affinity",
      "project_id": 1,
      "region": "Luxembourg 1",
      "region_id": 1,
      "servergroup_id": "47003067-550a-6f17-93b6-81ee16ba061e"
    }
  ]
}
`
var GetResponse = `
{
  "instances": [
    {
      "instance_id": "6d14f194-6c1e-49b3-9fc7-50dc8401eb74",
      "instance_name": "test_ruslan_aa2"
    }
  ],
  "name": "example_server_group",
  "policy": "anti-affinity",
  "project_id": 1,
  "region": "Luxembourg 1",
  "region_id": 1,
  "servergroup_id": "47003067-550a-6f17-93b6-81ee16ba061e"
}
`

var CreateRequest = `
{
  "name": "example_server_group",
  "policy": "anti-affinity"
}
`

var (
	Sg1 = servergroups.ServerGroup{
		ServerGroupID: "47003067-550a-6f17-93b6-81ee16ba061e",
		ProjectID:     1,
		RegionID:      1,
		Region:        "Luxembourg 1",
		Name:          "example_server_group",
		Instances: []servergroups.ServerGroupInstance{
			{
				InstanceID:   "6d14f194-6c1e-49b3-9fc7-50dc8401eb74",
				InstanceName: "test_ruslan_aa2",
			},
		},
		Policy: "anti-affinity",
	}
	ExpectedServerGroupSlice = []servergroups.ServerGroup{Sg1}
	Tasks1                   = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)
