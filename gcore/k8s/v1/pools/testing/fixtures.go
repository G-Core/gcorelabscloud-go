package testing

import (
	"net"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/pools"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
      "uuid": "908338b2-9217-4673-af0e-f0093139fbac",
      "image_id": "fedora-coreos",
      "flavor_id": "g1-standard-1-2",
      "min_node_count": 1,
      "role": "worker",
      "status": "CREATE_COMPLETE",
      "is_default": true,
      "max_node_count": null,
      "node_count": 1,
      "name": "test1"
    }
  ]
}
`

const GetResponse1 = `
{
  "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
  "uuid": "908338b2-9217-4673-af0e-f0093139fbac",
  "image_id": "fedora-coreos",
  "flavor_id": "g1-standard-1-2",
  "docker_volume_size": 10,
  "min_node_count": 1,
  "labels": {
	"gcloud_project_id":    "1",
	"gcloud_region_id":     "1",
	"gcloud_access_token":  "token",
	"gcloud_refresh_token": "token"
  },
  "role": "worker",
  "project_id": "46beed3938e6474390b530fefd6173d2",
  "cluster_id": "5e09faed-e742-404f-8a75-0ea5eb3c435f",
  "status_reason": "Stack CREATE completed successfully",
  "status": "CREATE_COMPLETE",
  "is_default": true,
  "node_addresses": [
    "192.168.0.5"
  ],
  "max_node_count": null,
  "node_count": 1,
  "name": "test1"
}
`

const UpdateResponse = GetResponse1

const CreateRequest = `
{
  "docker_volume_size": 5,
  "name": "test1",
  "node_count": 1,
  "flavor_id": "g1-standard-1-2"
}
`

const UpdateRequest = `
{
	"min_mode_count": 3,
	"max_mode_count": 4,
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

var (
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
	nodeAddress = net.ParseIP("192.168.0.5")
	labels      = map[string]string{
		"gcloud_project_id":    "1",
		"gcloud_region_id":     "1",
		"gcloud_access_token":  "token",
		"gcloud_refresh_token": "token",
	}
	PoolList1 = pools.ClusterListPool{
		UUID:         "908338b2-9217-4673-af0e-f0093139fbac",
		Name:         "test1",
		FlavorID:     "g1-standard-1-2",
		ImageID:      "fedora-coreos",
		NodeCount:    1,
		MinNodeCount: 1,
		MaxNodeCount: nil,
		IsDefault:    true,
		StackID:      "2f0d5d97-fb3c-4218-9201-34f804299510",
		Status:       "CREATE_COMPLETE",
		Role:         "worker",
	}
	Pool1 = pools.ClusterPool{
		ClusterID:        "5e09faed-e742-404f-8a75-0ea5eb3c435f",
		ProjectID:        "46beed3938e6474390b530fefd6173d2",
		Labels:           labels,
		NodeAddresses:    []net.IP{nodeAddress},
		StatusReason:     "Stack CREATE completed successfully",
		DockerVolumeSize: 10,
		ClusterListPool:  &PoolList1,
	}
	UpdatedPool1                 = Pool1
	ExpectedClusterListPoolSlice = []pools.ClusterListPool{PoolList1}
)
