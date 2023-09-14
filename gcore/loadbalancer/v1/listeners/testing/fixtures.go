package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
      "name": "lbaas_test_listener",
      "task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
      "pool_count": 1,
      "operating_status": "ONLINE",
      "protocol_port": 80,
      "id": "43658ea9-54bd-4807-90b1-925921c9a0d1",
      "protocol": "TCP",
      "provisioning_status": "ACTIVE",
      "allowed_cidrs": ["10.10.0.0/24"]
    }
  ]
}
`

const GetResponse = `
{
  "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
  "name": "lbaas_test_listener",
  "task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
  "pool_count": 1,
  "operating_status": "ONLINE",
  "protocol_port": 80,
  "id": "43658ea9-54bd-4807-90b1-925921c9a0d1",
  "protocol": "TCP",
  "provisioning_status": "ACTIVE",
  "allowed_cidrs": ["10.10.0.0/24"]
}
`

const CreateRequest = `
{
  "name": "lbaas_test_listener",
  "protocol_port": 80,
  "protocol": "TCP",
  "loadbalancer_id": "43658ea9-54bd-4807-90b1-925921c9a0d1",
  "insert_x_forwarded": false,
  "allowed_cidrs": ["10.10.0.0/24"]
}
`

const UpdateRequest = `
{
	"name": "lbaas_test_listener"
}	
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`
const UpdateResponse = `
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

var (
	creatorTaskID = "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22"
	taskID        = "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22"

	Listener1 = listeners.Listener{
		PoolCount:          1,
		ProtocolPort:       80,
		Protocol:           types.ProtocolTypeTCP,
		Name:               "lbaas_test_listener",
		ID:                 "43658ea9-54bd-4807-90b1-925921c9a0d1",
		ProvisioningStatus: types.ProvisioningStatusActive,
		OperationStatus:    types.OperatingStatusOnline,
		CreatorTaskID:      &creatorTaskID,
		TaskID:             &taskID,
    AllowedCIDRS:       []string{"10.10.0.0/24"},
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedListenersSlice = []listeners.Listener{Listener1}
)
