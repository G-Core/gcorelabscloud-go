package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "tenant_id": "fe5cc21270554c0d9d4cdc48ba574987",
      "task_state": null,
      "instance_description": "Testing",
      "instance_name": "Testing",
      "status": "ACTIVE",
      "instance_created": "2019-07-11T06:58:48Z",
      "vm_state": "active",
      "volumes": [
        {
          "id": "28bfe198-a003-4283-8dca-ab5da4a71b62",
          "delete_on_termination": false
        }
      ],
      "security_groups": [
        {
          "name": "default"
        }
      ],
      "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
      "task_id": "f28a4982-9be1-4e50-84e7-6d1a6d3f8a02",
      "creator_task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30",
      "addresses": {
        "net1": [
          {
            "type": "fixed",
            "addr": "10.0.0.17"
          },
          {
            "type": "floating",
            "addr": "92.38.157.215"
          }
        ],
        "net2": [
          {
            "type": "fixed",
            "addr": "192.168.68.68"
          }
        ]
      },
      "metadata": {
        "os_distro": "centos",
        "os_version": "1711-x64",
        "image_name": "cirros-0.3.5-x86_64-disk",
        "image_id": "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1",
        "snapshot_name": "test_snapshot",
        "snapshot_id": "c286cd13-fba9-4302-9cdb-4351a05a56ea",
        "task_id": "d1e1500b-e2be-40aa-9a4b-cc493fa1af30"
      },
      "flavor": {
        "flavor_name": "g1s-shared-1-0.5",
        "disk": 0,
        "flavor_id": "g1s-shared-1-0.5",
        "vcpus": 1,
        "ram": 512
      },
      "project_id": 1,
      "region_id": 1,
	  "region": "RegionOne"	
    }
  ]
}
`

const CreateRequest = `
{
  "flavor": "bm1-infrastructure-small",
  "interfaces": [
	{
	  "floating_ip": {
		"existing_floating_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
		"source": "existing"
	  },
	  "network_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
	  "subnet_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
	  "type": "subnet"
	}
  ],
  "keypair_name": "keypair",
  "names": [
	"name"
  ],
  "image_id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
  "password": "password",
  "username": "username"
}
`

const CreateResponse = `
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
)
