package testing

import (
	"net"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

const ListResponse = `
{
  "count": 1,
  "results": [
      {
          "flavor": "g2a-ai-fake-v1pod-8",
          "task_id": "b34d8be3-73b2-402b-92c8-16e944d65f0c",
          "volumes": [
                {
                  "volume_image_metadata": {
                      "signature_verified": "False",
                      "os_distro": "poplar-ubuntu",
                      "os_type": "linux",
                      "os_version": "20.04",
                      "vipu_version": "1.18.0",
                      "build_ts": "1694789419",
                      "sdk_version": "3.0.0",
                      "release_version": "1.8.4",
                      "display_order": "2004300",
                      "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
                      "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                      "checksum": "dcb3767a59b4c1f0fbc09b439d8bc789",
                      "container_format": "bare",
                      "disk_format": "qcow2",
                      "min_disk": "0",
                      "min_ram": "0",
                      "size": "5703401472"
                  },
                  "updated_at": "2023-09-28T15:24:34+0000",
                  "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "project_id": 516070,
                  "region": "ED-10 Preprod",
                  "region_id": 7,
                  "name": "ivandts_bootvolume",
                  "created_at": "2023-09-28T15:23:04+0000",
                  "bootable": true,
                  "attachments": [
                      {
                          "attached_at": "2023-09-28T15:24:34+0000",
                          "attachment_id": "a1f35e2b-afae-4caf-9f09-386c136cec45",
                          "server_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
                          "volume_id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                          "device": "/dev/vda"
                      }
                  ],
                  "volume_type": "standard",
                  "size": 20,
                  "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                  "status": "in-use",
                  "limiter_stats": {
                      "MBps_base_limit": 10,
                      "iops_base_limit": 120,
                      "MBps_burst_limit": 100,
                      "iops_burst_limit": 1200
                  },
                  "metadata": {
                      "task_id": "e673bba0-fcef-44d9-904c-824546b608ec"
                  },
                  "metadata_detailed": [
                      {
                          "key": "task_id",
                          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                          "read_only": true
                      }
                  ]
              }
          ],
          "creator_task_id": "b34d8be3-73b2-402b-92c8-16e944d65f0c",
          "user_data": "#cloud-config\nssh_pwauth: True\nusers:\n  - name: kolya\n    passwd: $6$rounds=4096$jB/jrhCWrbx65sHb$e5eLHfdJZ/IhiB06N0i/wPepo1fS3Y2o//D7C.jnw66mEqgPUWFuhGAOShC3lYF3eVGJOnEoWZ6N2fRCHj/4W.\n    lock-passwd: False\n    sudo:  ALL=(ALL:ALL) ALL\n",
          "project_id": 516070,
          "region": "ED-10 Preprod",
          "region_id": 7,
          "cluster_metadata_detailed": null,
          "task_status": "FINISHED",
          "created_at": "2023-09-28 15:25:06.115000",
          "cluster_name": "ivandts",
          "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
          "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
          "cluster_metadata": null,
          "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
          "interfaces": [
              {
                  "type": "any_subnet",
                  "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b"
              }
          ],
          "poplar_servers": [
              {
                  "flavor": {
                      "flavor_id": "g2a-ai-fake-v1pod-8",
                      "os_type": null,
                      "architecture": null,
                      "vcpus": 1,
                      "ram": 2048,
                      "flavor_name": "g2a-ai-fake-v1pod-8",
                      "hardware_description": {
                          "network": "2x100G",
                          "cpu": "1 vCPU",
                          "ram": "2GB RAM",
                          "ipu": "vPOD-8 (Classic)"
                      }
                  },
                  "task_id": null,
                  "instance_created": "2023-09-28T15:24:32Z",
                  "volumes": [
                      {
                          "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                          "delete_on_termination": false
                      }
                  ],
                  "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "instance_description": null,
                  "project_id": 516070,
                  "region": "ED-10 Preprod",
                  "region_id": 7,
                  "instance_name": "ivandts",
                  "vm_state": "active",
                  "task_state": null,
                  "addresses": {
                      "qa-alex-network": [
                          {
                              "addr": "10.10.0.247",
                              "type": "fixed"
                          }
                      ],
                      "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
                          {
                              "addr": "10.191.167.5",
                              "type": "fixed"
                          }
                      ]
                  },
                  "status": "ACTIVE",
                  "security_groups": [
                      {
                          "name": "default"
                      },
                      {
                          "name": "ivandts FE"
                      }
                  ],
                  "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
                  "keypair_name": null,
                  "metadata": {
                      "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
                      "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
                      "vipu_version": "1.18.0",
                      "poplar_sdk_version": "3.0.0",
                      "os_distro": "poplar-ubuntu",
                      "os_type": "linux",
                      "os_version": "20.04",
                      "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                      "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
                  },
                  "metadata_detailed": [
                      {
                          "key": "cluster_id",
                          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                          "read_only": false
                      },
                      {
                          "key": "image_id",
                          "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
                          "read_only": true
                      },
                      {
                          "key": "image_name",
                          "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                          "read_only": true
                      },
                      {
                          "key": "os_distro",
                          "value": "poplar-ubuntu",
                          "read_only": true
                      },
                      {
                          "key": "os_type",
                          "value": "linux",
                          "read_only": true
                      },
                      {
                          "key": "os_version",
                          "value": "20.04",
                          "read_only": true
                      },
                      {
                          "key": "poplar_sdk_version",
                          "value": "3.0.0",
                          "read_only": false
                      },
                      {
                          "key": "task_id",
                          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                          "read_only": true
                      },
                      {
                          "key": "vipu_version",
                          "value": "1.18.0",
                          "read_only": false
                      }
                  ]
              }
          ],
          "security_groups": [
              {
                  "security_groups": [
                      {
                          "name": "security-group-1"
                      }
                  ],
                  "network_id": "bf572176-2d95-4fe0-9de0-f54a5307fbe6",
                  "port_id": "d7136b4d-c5f3-4d3b-bd86-aeb01942cfc8"
              },
              {
                  "security_groups": [
                      {
                          "name": "security-group-2"
                      }
                  ],
                  "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
                  "port_id": "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142"
              }
          ],
          "keypair_name": null,
          "cluster_status": "ACTIVE"
      }
  ]
}
`

const GetResponse = `
{
  "flavor": "g2a-ai-fake-v1pod-8",
  "task_id": "b34d8be3-73b2-402b-92c8-16e944d65f0c",
  "volumes": [
      {
          "volume_image_metadata": {
              "signature_verified": "False",
              "os_distro": "poplar-ubuntu",
              "os_type": "linux",
              "os_version": "20.04",
              "vipu_version": "1.18.0",
              "build_ts": "1694789419",
              "sdk_version": "3.0.0",
              "release_version": "1.8.4",
              "display_order": "2004300",
              "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
              "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
              "checksum": "dcb3767a59b4c1f0fbc09b439d8bc789",
              "container_format": "bare",
              "disk_format": "qcow2",
              "min_disk": "0",
              "min_ram": "0",
              "size": "5703401472"
          },
          "updated_at": "2023-09-28T15:24:34+0000",
          "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
          "project_id": 516070,
          "region": "ED-10 Preprod",
          "region_id": 7,
          "name": "ivandts_bootvolume",
          "created_at": "2023-09-28T15:23:04+0000",
          "bootable": true,
          "attachments": [
              {
                  "attached_at": "2023-09-28T15:24:34+0000",
                  "attachment_id": "a1f35e2b-afae-4caf-9f09-386c136cec45",
                  "server_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
                  "volume_id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                  "device": "/dev/vda"
              }
          ],
          "volume_type": "standard",
          "size": 20,
          "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
          "status": "in-use",
          "limiter_stats": {
              "MBps_base_limit": 10,
              "iops_base_limit": 120,
              "MBps_burst_limit": 100,
              "iops_burst_limit": 1200
          },
          "metadata": {
              "task_id": "e673bba0-fcef-44d9-904c-824546b608ec"
          },
          "metadata_detailed": [
              {
                  "key": "task_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": true
              }
          ]
      }
  ],
  "creator_task_id": "b34d8be3-73b2-402b-92c8-16e944d65f0c",
  "user_data": "#cloud-config\nssh_pwauth: True\nusers:\n  - name: kolya\n    passwd: $6$rounds=4096$jB/jrhCWrbx65sHb$e5eLHfdJZ/IhiB06N0i/wPepo1fS3Y2o//D7C.jnw66mEqgPUWFuhGAOShC3lYF3eVGJOnEoWZ6N2fRCHj/4W.\n    lock-passwd: False\n    sudo:  ALL=(ALL:ALL) ALL\n",
  "project_id": 516070,
  "region": "ED-10 Preprod",
  "region_id": 7,
  "cluster_metadata_detailed": null,
  "task_status": "FINISHED",
  "created_at": "2023-09-28 15:25:06.115000",
  "cluster_name": "ivandts",
  "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
  "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
  "cluster_metadata": null,
  "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
  "interfaces": [
      {
          "type": "any_subnet",
          "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b"
      }
  ],
  "poplar_servers": [
      {
          "flavor": {
              "flavor_id": "g2a-ai-fake-v1pod-8",
              "os_type": null,
              "architecture": null,
              "vcpus": 1,
              "ram": 2048,
              "flavor_name": "g2a-ai-fake-v1pod-8",
              "hardware_description": {
                  "network": "2x100G",
                  "cpu": "1 vCPU",
                  "ram": "2GB RAM",
                  "ipu": "vPOD-8 (Classic)"
              }
          },
          "task_id": null,
          "instance_created": "2023-09-28T15:24:32Z",
          "volumes": [
              {
                  "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                  "delete_on_termination": false
              }
          ],
          "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
          "instance_description": null,
          "project_id": 516070,
          "region": "ED-10 Preprod",
          "region_id": 7,
          "instance_name": "ivandts",
          "vm_state": "active",
          "task_state": null,
          "addresses": {
              "qa-alex-network": [
                  {
                      "addr": "10.10.0.247",
                      "type": "fixed"
                  }
              ],
              "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
                  {
                      "addr": "10.191.167.5",
                      "type": "fixed"
                  }
              ]
          },
          "status": "ACTIVE",
          "security_groups": [
              {
                  "name": "default"
              },
              {
                  "name": "ivandts FE"
              }
          ],
          "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
          "keypair_name": null,
          "metadata": {
              "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "vipu_version": "1.18.0",
              "poplar_sdk_version": "3.0.0",
              "os_distro": "poplar-ubuntu",
              "os_type": "linux",
              "os_version": "20.04",
              "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
              "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
          },
          "metadata_detailed": [
              {
                  "key": "cluster_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": false
              },
              {
                  "key": "image_id",
                  "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
                  "read_only": true
              },
              {
                  "key": "image_name",
                  "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                  "read_only": true
              },
              {
                  "key": "os_distro",
                  "value": "poplar-ubuntu",
                  "read_only": true
              },
              {
                  "key": "os_type",
                  "value": "linux",
                  "read_only": true
              },
              {
                  "key": "os_version",
                  "value": "20.04",
                  "read_only": true
              },
              {
                  "key": "poplar_sdk_version",
                  "value": "3.0.0",
                  "read_only": false
              },
              {
                  "key": "task_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": true
              },
              {
                  "key": "vipu_version",
                  "value": "1.18.0",
                  "read_only": false
              }
          ]
      }
  ],
  "security_groups": [
      {
          "security_groups": [
              {
                  "name": "security-group-1"
              }
          ],
          "network_id": "bf572176-2d95-4fe0-9de0-f54a5307fbe6",
          "port_id": "d7136b4d-c5f3-4d3b-bd86-aeb01942cfc8"
      },
      {
          "security_groups": [
              {
                  "name": "security-group-2"
              }
          ],
          "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
          "port_id": "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142"
      }
  ],
  "keypair_name": null,
  "cluster_status": "ACTIVE"
}
`
const ClusterInterfacesResponse = `
{
  "count": 1,
  "results": [
      {
          "port_security_enabled": true,
          "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
          "ip_assignments": [
              {
                  "ip_address": "10.10.0.247",
                  "subnet_id": "8a5d4b01-4d80-4c7e-ba88-96162e3781a4"
              }
          ],
          "network_details": {
              "mtu": 1500,
              "project_id": null,
              "region": null,
              "region_id": null,
              "updated_at": "2023-09-21T06:24:34+0000",
              "subnets": [
                  {
                      "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
                      "enable_dhcp": true,
                      "host_routes": [],
                      "updated_at": "2023-09-21T06:24:34+0000",
                      "creator_task_id": "58cb0400-13d9-4539-8e7c-bd5e66edde2c",
                      "gateway_ip": "10.10.0.1",
                      "project_id": null,
                      "region": null,
                      "region_id": null,
                      "name": "qa-alex-subnet",
                      "created_at": "2023-09-21T06:24:34+0000",
                      "id": "8a5d4b01-4d80-4c7e-ba88-96162e3781a4",
                      "cidr": "10.10.0.0/24",
                      "dns_nameservers": [
                          "8.8.8.8",
                          "8.8.4.4"
                      ],
                      "ip_version": 4,
                      "has_router": false,
                      "metadata": []
                  }
              ],
              "name": "qa-alex-network",
              "external": false,
              "shared": false,
              "segmentation_id": 338,
              "created_at": "2023-09-21T06:24:13+0000",
              "creator_task_id": "5f4dd40a-158b-49f2-b1c3-8bf764318ab1",
              "type": "vxlan",
              "id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
              "metadata": []
          },
          "port_id": "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142",
          "floatingip_details": [],
          "mac_address": "fa:16:3e:f5:f2:6b"
      }
  ]
}
`

const CreateRequest = `
{
  "flavor": "g2a-ai-fake-v1pod-8",
  "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
  "interfaces": [
      {
          "network_id": "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
          "type": "any_subnet"
      }
  ],
  "username": "useruser",
  "password": "secret",
  "metadata": {
      "foo": "bar"
  },
  "name": "ivandts",
  "volumes": [
      {
          "boot_index": 0,
          "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f",
          "size": 20,                                                                                                                                                                                                                                     
          "source": "image",                                                                                                                                                                                                                               
          "type_name": "standard"                                                                                                                                                                                                                        
      }                                                                                                                                                                                                                                                    
  ]
}       
`

const ResizeRequest = `
{
  "instances_count": 2
}
`

const PortsListResponse = `
{
  "count": 1,
  "results": [
      {
          "id": "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142",
          "name": "port for instance ivandts",
          "security_groups": [
              {
                  "id": "77ae0765-f262-493a-ba32-d9892436ddd0",
                  "name": "ivandts FE"
              }
          ]
      }
  ]
}
`

const AssignSecurityGroupsRequest = `
{
  "name": "Test"
}
`

const UnAssignSecurityGroupsRequest = `
{
  "name": "Test"
}
`

const AttachInterfaceRequest = `
{
  "type": "subnet",
  "subnet_id": "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
}
`

const DetachInterfaceRequest = `
{
  "ip_address": "192.168.0.23",
  "port_id": "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d"
}
`

const AIClusterPowercycleResponse = `
{
  "count": 1,
  "results": [
      {
          "region_id": 7,
          "region": "ED-10 Preprod",
          "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
          "keypair_name": null,
          "status": "ACTIVE",
          "addresses": {
              "qa-alex-network": [
                  {
                      "type": "fixed",
                      "addr": "10.10.0.247"
                  }
              ],
              "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
                  {
                      "type": "fixed",
                      "addr": "10.191.167.5"
                  }
              ]
          },
          "instance_created": "2023-09-28T15:24:32Z",
          "instance_description": null,
          "vm_state": "active",
          "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
          "flavor": {
              "os_type": null,
              "ram": 2048,
              "hardware_description": {
                  "cpu": "1 vCPU",
                  "ram": "2GB RAM",
                  "ipu": "vPOD-8 (Classic)",
                  "network": "2x100G"
              },
              "vcpus": 1,
              "flavor_name": "g2a-ai-fake-v1pod-8",
              "architecture": null,
              "flavor_id": "g2a-ai-fake-v1pod-8"
          },
          "volumes": [
              {
                  "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                  "delete_on_termination": false
              }
          ],
          "project_id": 516070,
          "security_groups": [
              {
                  "name": "default"
              },
              {
                  "name": "ivandts FE"
              }
          ],
          "task_state": null,
          "metadata": {
              "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "vipu_version": "1.18.0",
              "poplar_sdk_version": "3.0.0",
              "os_distro": "poplar-ubuntu",
              "os_type": "linux",
              "os_version": "20.04",
              "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
              "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
          },
          "instance_name": "ivandts",
          "task_id": null,
          "metadata_detailed": [
              {
                  "key": "cluster_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": false
              },
              {
                  "key": "image_id",
                  "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
                  "read_only": true
              },
              {
                  "key": "image_name",
                  "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                  "read_only": true
              },
              {
                  "key": "os_distro",
                  "value": "poplar-ubuntu",
                  "read_only": true
              },
              {
                  "key": "os_type",
                  "value": "linux",
                  "read_only": true
              },
              {
                  "key": "os_version",
                  "value": "20.04",
                  "read_only": true
              },
              {
                  "key": "poplar_sdk_version",
                  "value": "3.0.0",
                  "read_only": false
              },
              {
                  "key": "task_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": true
              },
              {
                  "key": "vipu_version",
                  "value": "1.18.0",
                  "read_only": false
              }
          ]
      }
  ]
}
`

const AIInstancePowercycleResponse = `
{
  "region_id": 7,
  "region": "ED-10 Preprod",
  "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
  "keypair_name": null,
  "status": "ACTIVE",
  "addresses": {
      "qa-alex-network": [
          {
              "type": "fixed",
              "addr": "10.10.0.247"
          }
      ],
      "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
          {
              "type": "fixed",
              "addr": "10.191.167.5"
          }
      ]
  },
  "instance_created": "2023-09-28T15:24:32Z",
  "instance_description": null,
  "vm_state": "active",
  "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
  "flavor": {
      "os_type": null,
      "ram": 2048,
      "hardware_description": {
          "cpu": "1 vCPU",
          "ram": "2GB RAM",
          "ipu": "vPOD-8 (Classic)",
          "network": "2x100G"
      },
      "vcpus": 1,
      "flavor_name": "g2a-ai-fake-v1pod-8",
      "architecture": null,
      "flavor_id": "g2a-ai-fake-v1pod-8"
  },
  "volumes": [
      {
          "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
          "delete_on_termination": false
      }
  ],
  "project_id": 516070,
  "security_groups": [
      {
          "name": "default"
      },
      {
          "name": "ivandts FE"
      }
  ],
  "task_state": null,
  "metadata": {
      "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
      "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
      "vipu_version": "1.18.0",
      "poplar_sdk_version": "3.0.0",
      "os_distro": "poplar-ubuntu",
      "os_type": "linux",
      "os_version": "20.04",
      "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
      "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
  },
  "instance_name": "ivandts",
  "task_id": null,
  "metadata_detailed": [
      {
          "key": "cluster_id",
          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
          "read_only": false
      },
      {
          "key": "image_id",
          "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
          "read_only": true
      },
      {
          "key": "image_name",
          "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
          "read_only": true
      },
      {
          "key": "os_distro",
          "value": "poplar-ubuntu",
          "read_only": true
      },
      {
          "key": "os_type",
          "value": "linux",
          "read_only": true
      },
      {
          "key": "os_version",
          "value": "20.04",
          "read_only": true
      },
      {
          "key": "poplar_sdk_version",
          "value": "3.0.0",
          "read_only": false
      },
      {
          "key": "task_id",
          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
          "read_only": true
      },
      {
          "key": "vipu_version",
          "value": "1.18.0",
          "read_only": false
      }
  ]
}
`

const AIInstanceRebootResponse = `
{
  "region_id": 7,
  "region": "ED-10 Preprod",
  "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
  "keypair_name": null,
  "status": "ACTIVE",
  "addresses": {
      "qa-alex-network": [
          {
              "type": "fixed",
              "addr": "10.10.0.247"
          }
      ],
      "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
          {
              "type": "fixed",
              "addr": "10.191.167.5"
          }
      ]
  },
  "instance_created": "2023-09-28T15:24:32Z",
  "instance_description": null,
  "vm_state": "active",
  "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
  "flavor": {
      "os_type": null,
      "ram": 2048,
      "hardware_description": {
          "cpu": "1 vCPU",
          "ram": "2GB RAM",
          "ipu": "vPOD-8 (Classic)",
          "network": "2x100G"
      },
      "vcpus": 1,
      "flavor_name": "g2a-ai-fake-v1pod-8",
      "architecture": null,
      "flavor_id": "g2a-ai-fake-v1pod-8"
  },
  "volumes": [
      {
          "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
          "delete_on_termination": false
      }
  ],
  "project_id": 516070,
  "security_groups": [
      {
          "name": "default"
      },
      {
          "name": "ivandts FE"
      }
  ],
  "task_state": null,
  "metadata": {
      "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
      "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
      "vipu_version": "1.18.0",
      "poplar_sdk_version": "3.0.0",
      "os_distro": "poplar-ubuntu",
      "os_type": "linux",
      "os_version": "20.04",
      "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
      "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
  },
  "instance_name": "ivandts",
  "task_id": null,
  "metadata_detailed": [
      {
          "key": "cluster_id",
          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
          "read_only": false
      },
      {
          "key": "image_id",
          "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
          "read_only": true
      },
      {
          "key": "image_name",
          "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
          "read_only": true
      },
      {
          "key": "os_distro",
          "value": "poplar-ubuntu",
          "read_only": true
      },
      {
          "key": "os_type",
          "value": "linux",
          "read_only": true
      },
      {
          "key": "os_version",
          "value": "20.04",
          "read_only": true
      },
      {
          "key": "poplar_sdk_version",
          "value": "3.0.0",
          "read_only": false
      },
      {
          "key": "task_id",
          "value": "e673bba0-fcef-44d9-904c-824546b608ec",
          "read_only": true
      },
      {
          "key": "vipu_version",
          "value": "1.18.0",
          "read_only": false
      }
  ]
}
`

const AIClusterRebootResponse = `
{
  "count": 1,
  "results": [
      {
          "region_id": 7,
          "region": "ED-10 Preprod",
          "instance_id": "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
          "keypair_name": null,
          "status": "ACTIVE",
          "addresses": {
              "qa-alex-network": [
                  {
                      "type": "fixed",
                      "addr": "10.10.0.247"
                  }
              ],
              "ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": [
                  {
                      "type": "fixed",
                      "addr": "10.191.167.5"
                  }
              ]
          },
          "instance_created": "2023-09-28T15:24:32Z",
          "instance_description": null,
          "vm_state": "active",
          "creator_task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
          "flavor": {
              "os_type": null,
              "ram": 2048,
              "hardware_description": {
                  "cpu": "1 vCPU",
                  "ram": "2GB RAM",
                  "ipu": "vPOD-8 (Classic)",
                  "network": "2x100G"
              },
              "vcpus": 1,
              "flavor_name": "g2a-ai-fake-v1pod-8",
              "architecture": null,
              "flavor_id": "g2a-ai-fake-v1pod-8"
          },
          "volumes": [
              {
                  "id": "459bf28d-df63-45d2-a462-6c216e571ddc",
                  "delete_on_termination": false
              }
          ],
          "project_id": 516070,
          "security_groups": [
              {
                  "name": "default"
              },
              {
                  "name": "ivandts FE"
              }
          ],
          "task_state": null,
          "metadata": {
              "task_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "cluster_id": "e673bba0-fcef-44d9-904c-824546b608ec",
              "vipu_version": "1.18.0",
              "poplar_sdk_version": "3.0.0",
              "os_distro": "poplar-ubuntu",
              "os_type": "linux",
              "os_version": "20.04",
              "image_name": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
              "image_id": "06e62653-1f88-4d38-9aa6-62833e812b4f"
          },
          "instance_name": "ivandts",
          "task_id": null,
          "metadata_detailed": [
              {
                  "key": "cluster_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": false
              },
              {
                  "key": "image_id",
                  "value": "06e62653-1f88-4d38-9aa6-62833e812b4f",
                  "read_only": true
              },
              {
                  "key": "image_name",
                  "value": "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
                  "read_only": true
              },
              {
                  "key": "os_distro",
                  "value": "poplar-ubuntu",
                  "read_only": true
              },
              {
                  "key": "os_type",
                  "value": "linux",
                  "read_only": true
              },
              {
                  "key": "os_version",
                  "value": "20.04",
                  "read_only": true
              },
              {
                  "key": "poplar_sdk_version",
                  "value": "3.0.0",
                  "read_only": false
              },
              {
                  "key": "task_id",
                  "value": "e673bba0-fcef-44d9-904c-824546b608ec",
                  "read_only": true
              },
              {
                  "key": "vipu_version",
                  "value": "1.18.0",
                  "read_only": false
              }
          ]
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

const MetadataResponse = `
{
  "key": "cost-center",
  "value": "Atlanta",
  "read_only": false
}
`

const MetadataCreateRequest = `
{
"test1": "test1", 
"test2": "test2"
}
`

const InstanceConsoleResponse = `
{
    "remote_console": 
    {
        "url": "https://console-novnc-ed10.cloud.gcorelabs.com/vnc_auto.html?path=token%3Ddf5d4b4f-f78c-421f-9131-b6be2facf9bd",
        "type": "novnc",
        "protocol": "vnc"
    }
}
`

var (
	ip1                    = net.ParseIP("10.10.0.247")
	ip2                    = net.ParseIP("10.191.167.5")
	tm, _                  = time.Parse(gcorecloud.RFC3339MilliNoZ, "2023-09-28 15:25:06.115000")
	createdTime            = gcorecloud.JSONRFC3339MilliNoZ{Time: tm}
	volumeCreatedTime, _   = time.Parse(gcorecloud.RFC3339Z, "2023-09-28T15:23:04+0000")
	volumeUpdatedTime, _   = time.Parse(gcorecloud.RFC3339Z, "2023-09-28T15:24:34+0000")
	volumeAttachedTime, _  = time.Parse(gcorecloud.RFC3339Z, "2023-09-28T15:24:34+0000")
	instanceCreatedTime, _ = time.Parse(gcorecloud.RFC3339ZZ, "2023-09-28T15:24:32Z")
	taskID                 = "b34d8be3-73b2-402b-92c8-16e944d65f0c"
	creatorTaskID          = "e673bba0-fcef-44d9-904c-824546b608ec"

	AICluster1 = ai.AICluster{
		ClusterID:     "e673bba0-fcef-44d9-904c-824546b608ec",
		ClusterName:   "ivandts",
		ClusterStatus: "ACTIVE",
		TaskID:        &taskID,
		TaskStatus:    "FINISHED",
		CreatedAt:     createdTime,
		ImageID:       "06e62653-1f88-4d38-9aa6-62833e812b4f",
		ImageName:     "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
		Flavor:        "g2a-ai-fake-v1pod-8",
		Volumes: []volumes.Volume{
			{
				CreatedAt:  gcorecloud.JSONRFC3339Z{Time: volumeCreatedTime},
				UpdatedAt:  gcorecloud.JSONRFC3339Z{Time: volumeUpdatedTime},
				VolumeType: "standard",
				ID:         "459bf28d-df63-45d2-a462-6c216e571ddc",
				Name:       "ivandts_bootvolume",
				RegionName: "ED-10 Preprod",
				Status:     "in-use",
				Size:       20,
				Bootable:   true,
				ProjectID:  516070,
				RegionID:   7,
				Attachments: []volumes.Attachment{
					{
						ServerID:     "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
						AttachmentID: "a1f35e2b-afae-4caf-9f09-386c136cec45",
						AttachedAt:   gcorecloud.JSONRFC3339Z{Time: volumeAttachedTime},
						VolumeID:     "459bf28d-df63-45d2-a462-6c216e571ddc",
						Device:       "/dev/vda",
					},
				},
				Metadata: []metadata.Metadata{
					{
						Key:      "task_id",
						Value:    "e673bba0-fcef-44d9-904c-824546b608ec",
						ReadOnly: true,
					},
				},
				CreatorTaskID: creatorTaskID,
				VolumeImageMetadata: volumes.VolumeImageMetadata{
					ContainerFormat: "bare",
					MinRAM:          "0",
					DiskFormat:      "qcow2",
					ImageName:       "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
					ImageID:         "06e62653-1f88-4d38-9aa6-62833e812b4f",
					MinDisk:         "0",
					Checksum:        "dcb3767a59b4c1f0fbc09b439d8bc789",
					Size:            "5703401472",
				},
			},
		},
		SecurityGroups: []ai.PoplarInterfaceSecGrop{
			{
				PortID:         "d7136b4d-c5f3-4d3b-bd86-aeb01942cfc8",
				NetworkID:      "bf572176-2d95-4fe0-9de0-f54a5307fbe6",
				SecurityGroups: []gcorecloud.ItemName{{Name: "security-group-1"}},
			},
			{
				PortID:         "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142",
				NetworkID:      "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
				SecurityGroups: []gcorecloud.ItemName{{Name: "security-group-2"}},
			},
		},
		Interfaces: []ai.AIClusterInterface{
			{
				Type:      "any_subnet",
				NetworkID: "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
			},
		},
		UserData: "#cloud-config\nssh_pwauth: True\nusers:\n  - name: kolya\n    passwd: $6$rounds=4096$jB/jrhCWrbx65sHb$e5eLHfdJZ/IhiB06N0i/wPepo1fS3Y2o//D7C.jnw66mEqgPUWFuhGAOShC3lYF3eVGJOnEoWZ6N2fRCHj/4W.\n    lock-passwd: False\n    sudo:  ALL=(ALL:ALL) ALL\n",
		PoplarServer: []instances.Instance{
			{
				ID:               "a2ff6283-09f9-4c2a-a96f-0bedf7b3dd2d",
				Name:             "ivandts",
				CreatedAt:        gcorecloud.JSONRFC3339ZZ{Time: instanceCreatedTime},
				Status:           "ACTIVE",
				VMState:          "active",
				AvailabilityZone: "nova",
				Flavor: flavors.Flavor{
					FlavorID:   "g2a-ai-fake-v1pod-8",
					FlavorName: "g2a-ai-fake-v1pod-8",
					HardwareDescription: &flavors.HardwareDescription{
						CPU:     "1 vCPU",
						Network: "2x100G",
						RAM:     "2GB RAM",
						IPU:     "vPOD-8 (Classic)",
					},
					RAM:   2048,
					VCPUS: 1,
				},
				Metadata: map[string]interface{}{
					"task_id":            "e673bba0-fcef-44d9-904c-824546b608ec",
					"cluster_id":         "e673bba0-fcef-44d9-904c-824546b608ec",
					"vipu_version":       "1.18.0",
					"poplar_sdk_version": "3.0.0",
					"os_distro":          "poplar-ubuntu",
					"os_type":            "linux",
					"os_version":         "20.04",
					"image_name":         "ubuntu-20.04-x64-poplar-ironic-1.18.0-3.0.0",
					"image_id":           "06e62653-1f88-4d38-9aa6-62833e812b4f",
				},
				Volumes: []instances.InstanceVolume{
					{
						ID:                  "459bf28d-df63-45d2-a462-6c216e571ddc",
						DeleteOnTermination: false,
					},
				},
				Addresses: map[string][]instances.InstanceAddress{
					"qa-alex-network": {
						{
							Type:    "fixed",
							Address: ip1,
						},
					},
					"ipu-cluster-rdma-network-e673bba0-fcef-44d9-904c-824546b608ec": {
						{
							Type:    "fixed",
							Address: ip2,
						},
					},
				},
				SecurityGroups: []gcorecloud.ItemName{
					{
						Name: "default",
					},
					{
						Name: "ivandts FE",
					},
				},
				CreatorTaskID: &creatorTaskID,
				ProjectID:     516070,
				RegionID:      7,
				Region:        "ED-10 Preprod",
			},
		},
		ProjectID: 516070,
		RegionID:  7,
		Region:    "ED-10 Preprod",
	}

	PortID                         = "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142"
	PortMac, _                     = gcorecloud.ParseMacString("fa:16:3e:f5:f2:6b")
	PortIP1                        = net.ParseIP("10.10.0.247")
	PortNetworkUpdatedAt, _        = time.Parse(gcorecloud.RFC3339Z, "2023-09-21T06:24:34+0000")
	PortNetworkCreatedAt, _        = time.Parse(gcorecloud.RFC3339Z, "2023-09-21T06:24:13+0000")
	PortNetworkSubnet1CreatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2023-09-21T06:24:34+0000")
	PortNetworkSubnet1UpdatedAt, _ = time.Parse(gcorecloud.RFC3339Z, "2023-09-21T06:24:34+0000")
	PortNetworkSubnet1Cidr, _      = gcorecloud.ParseCIDRString("10.10.0.0/24")
	SubnetCreatorTaskID            = "58cb0400-13d9-4539-8e7c-bd5e66edde2c"
	NetworkDetailsCreatorTask      = "5f4dd40a-158b-49f2-b1c3-8bf764318ab1"
	SecurityGroup1                 = gcorecloud.ItemIDName{
		ID:   "77ae0765-f262-493a-ba32-d9892436ddd0",
		Name: "ivandts FE",
	}
	AIClusterPort1 = ai.AIClusterPort{
		ID:             "f3dcadf8-a4a5-4e5a-af7e-4c5902cd4142",
		Name:           "port for instance ivandts",
		SecurityGroups: ExpectedSecurityGroupsSlice,
	}
	AIClusterInterface1 = ai.Interface{
		PortID:              PortID,
		MacAddress:          *PortMac,
		PortSecurityEnabled: true,
		NetworkID:           "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
		IPAssignments: []instances.PortIP{
			{
				IPAddress: PortIP1,
				SubnetID:  "8a5d4b01-4d80-4c7e-ba88-96162e3781a4",
			},
		},
		NetworkDetails: instances.NetworkDetail{
			Mtu:           1500,
			UpdatedAt:     &gcorecloud.JSONRFC3339Z{Time: PortNetworkUpdatedAt},
			CreatedAt:     gcorecloud.JSONRFC3339Z{Time: PortNetworkCreatedAt},
			ID:            "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
			External:      false,
			Default:       false,
			Shared:        false,
			Name:          "qa-alex-network",
			CreatorTaskID: &NetworkDetailsCreatorTask,
			Subnets: []instances.Subnet{
				{
					ID:            "8a5d4b01-4d80-4c7e-ba88-96162e3781a4",
					Name:          "qa-alex-subnet",
					IPVersion:     gcorecloud.IPv4,
					EnableDHCP:    true,
					Cidr:          *PortNetworkSubnet1Cidr,
					CreatedAt:     gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet1CreatedAt},
					UpdatedAt:     &gcorecloud.JSONRFC3339Z{Time: PortNetworkSubnet1UpdatedAt},
					NetworkID:     "518ba531-496b-4676-8ea4-68e2ed3b2e4b",
					CreatorTaskID: &SubnetCreatorTaskID,
				},
			},
		},
		FloatingIPDetails: []instances.FloatingIP{},
	}
	ExpectedAIClusterSlice           = []ai.AICluster{AICluster1}
	ExpectedAIClusterInterfacesSlice = []ai.Interface{AIClusterInterface1}
	ExpectedSecurityGroupsSlice      = []gcorecloud.ItemIDName{SecurityGroup1}
	ExpectedPortsSlice               = []ai.AIClusterPort{AIClusterPort1}

	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
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
	Console              = ai.RemoteConsole{
		URL:      "https://console-novnc-ed10.cloud.gcorelabs.com/vnc_auto.html?path=token%3Ddf5d4b4f-f78c-421f-9131-b6be2facf9bd",
		Type:     "novnc",
		Protocol: "vnc",
	}
)
