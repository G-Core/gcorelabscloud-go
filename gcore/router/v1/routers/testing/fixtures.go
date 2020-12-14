package testing

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const CreateRequest = `
{
	"name": "router",
	"external_gateway_info": {
		"type": "manual",
		"enable_snat": true,
        "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
	},
	"interfaces": [
		{
			"type": "subnet",
			"subnet_id": "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
		}
	],
	"routes": [
		{
			"destination": "10.0.3.0/24",
			"nexthop": "10.0.0.13"
		}
	]
}
`

const UpdateRequest = `
{
	"name": "update_router",
	"external_gateway_info": {
		"type": "manual",
		"enable_snat": false,
        "network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f2"
	},
	"routes": [
		{
			"destination": "10.0.4.0/24",
			"nexthop": "10.0.0.14"
		}
	]
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

var taskID = "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"

var (
	Router1 = routers.Router{
		ID:     "e7944e55-f957-413d-aa56-fdc876543113",
		Name:   "router",
		Status: "ACTIVE",
		ExternalGatewayInfo: routers.ExtGatewayInfo{
			EnableSNat: true,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f1",
			ExternalFixedIPs: []routers.ExtFixedIPs{
				{
					IPAddress: "172.24.4.6",
					SubnetID:  "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
				},
			},
		},
		TaskID:        taskID,
		CreatorTaskID: taskID,
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
	}

	Router2 = routers.Router{
		ID:     "e7944e55-f957-413d-aa56-fdc876543113",
		Name:   "update_router",
		Status: "ACTIVE",
		ExternalGatewayInfo: routers.ExtGatewayInfo{
			EnableSNat: false,
			NetworkID:  "ee2402d0-f0cd-4503-9b75-69be1d11c5f2",
			ExternalFixedIPs: []routers.ExtFixedIPs{
				{
					IPAddress: "172.24.4.7",
					SubnetID:  "b930d7f6-ceb7-40a0-8b81-a425dd994ccf",
				},
			},
		},
		TaskID:        taskID,
		CreatorTaskID: taskID,
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
	}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedRouterSlice = []routers.Router{Router1}
)

var ListResponse = fmt.Sprintf(`
{
	"count": 1,
	"results": [
		{
			"id": "e7944e55-f957-413d-aa56-fdc876543113",
			"name": "router",
			"status": "ACTIVE",
			"external_gateway_info": {
				"enable_snat": true,
				"external_fixed_ips": [
					{
						"ip_address": "172.24.4.6",
						"subnet_id": "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
					}
				],
			"network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
			},
			"task_id": "%[1]s",
			"creator_task_id": "%[1]s",
			"project_id": %d,
			"region_id": %d
    	}
  	]
}
`, taskID, fake.ProjectID, fake.RegionID)

var GetResponse = fmt.Sprintf(`
{
	"id": "e7944e55-f957-413d-aa56-fdc876543113",
	"name": "router",
	"status": "ACTIVE",
	"external_gateway_info": {
		"enable_snat": true,
		"external_fixed_ips": [
			{
				"ip_address": "172.24.4.6",
				"subnet_id": "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
			}
		],
	"network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
	},
	"task_id": "%[1]s",
	"creator_task_id": "%[1]s",
	"project_id": %d,
	"region_id": %d
}
`, taskID, fake.ProjectID, fake.RegionID)

var GetResponseUpdate = fmt.Sprintf(`
{
	"id": "e7944e55-f957-413d-aa56-fdc876543113",
	"name": "update_router",
	"status": "ACTIVE",
	"external_gateway_info": {
		"enable_snat": false,
		"external_fixed_ips": [
			{
				"ip_address": "172.24.4.7",
				"subnet_id": "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
			}
		],
	"network_id": "ee2402d0-f0cd-4503-9b75-69be1d11c5f2"
	},
	"task_id": "%[1]s",
	"creator_task_id": "%[1]s",
	"project_id": %d,
	"region_id": %d
}
`, taskID, fake.ProjectID, fake.RegionID)
