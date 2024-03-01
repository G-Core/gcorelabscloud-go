package testing

import (
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "flavor_id": "g0-standard-2-4",
      "min_node_count": 1,
      "name": "pool-1",
      "created_at": "2023-08-28T09:40:39Z",
      "id": "f3446423-0a82-475a-a1bd-31ce788ace9e",
      "boot_volume_size": 50,
      "max_node_count": 2,
      "is_public_ipv4": false,
      "status": "Running",
      "node_count": 1,
      "boot_volume_type": "ssd_hiiops",
      "auto_healing_enabled": false
    }
  ]
}
`

const CreateRequest = `
{
  "name": "pool-1",
  "flavor_id": "g0-standard-2-4",
  "min_node_count": 1,
  "max_node_count": 2,
  "boot_volume_size": 50,
  "boot_volume_type": "ssd_hiiops",
  "servergroup_policy": "affinity"
}
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const GetResponse = `
{
  "flavor_id": "g0-standard-2-4",
  "min_node_count": 1,
  "name": "pool-1",
  "created_at": "2023-08-28T09:40:39Z",
  "id": "f3446423-0a82-475a-a1bd-31ce788ace9e",
  "boot_volume_size": 50,
  "max_node_count": 2,
  "is_public_ipv4": false,
  "status": "Running",
  "node_count": 1,
  "boot_volume_type": "ssd_hiiops",
  "auto_healing_enabled": false,
  "servergroup_policy": "affinity",
  "servergroup_id": "f3446423-0a82-475a-a1bd-31ce788ace9e",
  "servergroup_name": "affinity"
}
`

const UpdateRequest = `
{
  "min_mode_count": 1,
  "max_mode_count": 3,
}
`

const UpdateResponse = `
{
  "flavor_id": "g0-standard-2-4",
  "min_node_count": 1,
  "name": "pool-1",
  "created_at": "2023-08-28T09:40:39Z",
  "id": "f3446423-0a82-475a-a1bd-31ce788ace9e",
  "boot_volume_size": 50,
  "max_node_count": 3,
  "is_public_ipv4": false,
  "status": "Running",
  "node_count": 2,
  "boot_volume_type": "ssd_hiiops",
  "auto_healing_enabled": false,
  "servergroup_policy": "affinity",
  "servergroup_id": "f3446423-0a82-475a-a1bd-31ce788ace9e",
  "servergroup_name": "affinity"
}
`

const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ResizeRequest = `
{
  "node_count": 2
}
`

const ResizeResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ListInstancesResponse = `
{
  "count": 1,
  "results": [
    {
      "vm_state": "active",
      "region_id": 7,
      "volumes": [
        {
          "delete_on_termination": false,
          "id": "1ed838bb-2072-42a3-a5f6-d09777a3b023"
        }
      ],
      "region": "ED-10 Preprod",
      "status": "ACTIVE",
      "instance_created": "2023-08-28T09:40:39Z",
      "creator_task_id": "9640f68f-5748-4113-90bd-67a66e985e43",
      "keypair_name": "73a53a48-1f94-4f5c-9990-d44c8e60d992",
      "task_state": null,
      "project_id": 1234,
      "addresses": {
        "cluster-1": [
          {
            "addr": "10.42.42.179",
            "type": "fixed"
          }
        ]
      },
      "metadata": {
        "capgc/infra_machine_name": "cluster-1-pool-1-machine-template-j6c5g",
        "capgc/infra_machine_uid": "1a08139a-d441-4311-8dfb-442fe366be95",
        "capgc/kubernetes_version": "v1.26.7",
        "gcloud_cluster_name": "cluster-1",
        "gcloud_service": "k8s",
        "gcloud_service_type": "worker",
        "task_id": "9640f68f-5748-4113-90bd-67a66e985e43",
        "os_distro": "Ubuntu",
        "os_type": "linux",
        "os_version": "22.04",
        "image_name": "gcloud-k8s-v1.26.7-xUbuntu_22.04-worker-0.0.1.raw",
        "image_id": "d488fd8c-e70c-4bc2-b2b5-260960a083a2"
      },
      "instance_id": "2246207d-fb9f-4ea4-acea-5b2cf77ff46b",
      "flavor": {
        "flavor_id": "g0-standard-2-4",
        "flavor_name": "g0-standard-2-4",
        "vcpus": 2,
        "ram": 4096
      },
      "security_groups": [
        {
          "name": "cluster-1-7-1234-worker"
        }
      ],
      "instance_description": null,
      "instance_name": "cluster-1-pool-1-machine-deployment-56bc6958d-jz6m4",
      "metadata_detailed": [
        {
          "key": "capgc/infra_machine_name",
          "value": "gbernady-pool-1-machine-template-j6c5g",
          "read_only": false
        },
        {
          "key": "capgc/infra_machine_uid",
          "value": "1a08139a-d441-4311-8dfb-442fe366be95",
          "read_only": false
        },
        {
          "key": "capgc/kubernetes_version",
          "value": "v1.26.7",
          "read_only": false
        },
        {
          "key": "image_id",
          "value": "d488fd8c-e70c-4bc2-b2b5-260960a083a2",
          "read_only": true
        },
        {
          "key": "image_name",
          "value": "gcloud-k8s-v1.27.4-xUbuntu_22.04-worker-0.0.1.raw",
          "read_only": true
        },
        {
          "key": "os_distro",
          "value": "Ubuntu",
          "read_only": true
        },
        {
          "key": "os_type",
          "value": "linux",
          "read_only": true
        },
        {
          "key": "os_version",
          "value": "22.04",
          "read_only": true
        },
        {
          "key": "task_id",
          "value": "9640f68f-5748-4113-90bd-67a66e985e43",
          "read_only": true
        }
      ]
    }
  ]
}
`

var (
	createdTime, _ = time.Parse(time.RFC3339, "2023-08-28T09:40:39Z")
	creatorTaskID  = "9640f68f-5748-4113-90bd-67a66e985e43"
	Cluster1Name   = "cluster-1"
	Pool1          = pools.ClusterPool{
		ID:                 "f3446423-0a82-475a-a1bd-31ce788ace9e",
		Name:               "pool-1",
		FlavorID:           "g0-standard-2-4",
		NodeCount:          1,
		MinNodeCount:       1,
		MaxNodeCount:       2,
		Status:             "Running",
		BootVolumeType:     volumes.SsdHiIops,
		BootVolumeSize:     50,
		AutoHealingEnabled: false,
		CreatedAt:          createdTime,
		IsPublicIPv4:       false,
		ServerGroupPolicy:  "affinity",
		ServerGroupID:      "f3446423-0a82-475a-a1bd-31ce788ace9e",
		ServerGroupName:    "affinity",
	}
	Instance1 = instances.Instance{
		ID:        "2246207d-fb9f-4ea4-acea-5b2cf77ff46b",
		Name:      "cluster-1-pool-1-machine-deployment-56bc6958d-jz6m4",
		CreatedAt: gcorecloud.JSONRFC3339ZZ{Time: createdTime},
		Status:    "ACTIVE",
		VMState:   "active",
		Flavor: flavors.Flavor{
			FlavorID:   "g0-standard-2-4",
			FlavorName: "g0-standard-2-4",
			VCPUS:      2,
			RAM:        4096,
		},
		Metadata: map[string]interface{}{
			"capgc/infra_machine_name": "cluster-1-pool-1-machine-template-j6c5g",
			"capgc/infra_machine_uid":  "1a08139a-d441-4311-8dfb-442fe366be95",
			"capgc/kubernetes_version": "v1.26.7",
			"gcloud_cluster_name":      "cluster-1",
			"gcloud_service":           "k8s",
			"gcloud_service_type":      "worker",
			"task_id":                  "9640f68f-5748-4113-90bd-67a66e985e43",
			"os_distro":                "Ubuntu",
			"os_type":                  "linux",
			"os_version":               "22.04",
			"image_name":               "gcloud-k8s-v1.26.7-xUbuntu_22.04-worker-0.0.1.raw",
			"image_id":                 "d488fd8c-e70c-4bc2-b2b5-260960a083a2",
		},
		Volumes: []instances.InstanceVolume{
			{
				ID:                  "1ed838bb-2072-42a3-a5f6-d09777a3b023",
				DeleteOnTermination: false,
			},
		},
		Addresses: map[string][]instances.InstanceAddress{
			"cluster-1": {
				{
					Address: net.ParseIP("10.42.42.179"),
					Type:    types.AddressTypeFixed,
				},
			},
		},
		SecurityGroups: []gcorecloud.ItemName{
			{
				Name: "cluster-1-7-1234-worker",
			},
		},
		CreatorTaskID:    &creatorTaskID,
		ProjectID:        1234,
		RegionID:         7,
		Region:           "ED-10 Preprod",
		AvailabilityZone: "nova",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
	ExpectedClusterPoolListSlice = []pools.ClusterPool{Pool1}
	ExpectedInstancesSlice       = []instances.Instance{Instance1}
)
