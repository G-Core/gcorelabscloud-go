package testing

import (
	"fmt"

	"github.com/G-Core/gcorelabscloud-go/gcore/quota/v2/quotas"
)

const GlobalResponse = `
{
  "keypair_count_limit": 100,
  "keypair_count_usage": 0,
  "project_count_limit": 2,
  "project_count_usage": 1
}
`

const RegionalResponse = `
{
  "region_id": 1,
  "baremetal_basic_count_limit": 0,
  "baremetal_basic_count_usage": 0,
  "baremetal_hf_count_limit": 0,
  "baremetal_hf_count_usage": 0,
  "baremetal_infrastructure_count_limit": 0,
  "baremetal_infrastructure_count_usage": 0,
  "baremetal_network_count_limit": 2,
  "baremetal_network_count_usage": 0,
  "baremetal_storage_count_limit": 0,
  "baremetal_storage_count_usage": 0,
  "cpu_count_limit": 2,
  "cpu_count_usage": 0,
  "external_ip_count_limit": 2,
  "external_ip_count_usage": 0,
  "firewall_count_limit": 2,
  "firewall_count_usage": 0,
  "floating_count_limit": 0,
  "floating_count_usage": 0,
  "gpu_count_limit": 0,
  "gpu_count_usage": 0,
  "image_count_limit": 2,
  "image_count_usage": 0,
  "image_size_limit": 26843545600,
  "image_size_usage": 0,
  "loadbalancer_count_limit": 0,
  "loadbalancer_count_usage": 0,
  "network_count_limit": 0,
  "network_count_usage": 0,
  "ram_limit": 4096,
  "ram_usage": 0,
  "router_count_limit": 0,
  "router_count_usage": 0,
  "secret_count_limit": 10,
  "secret_count_usage": 0,
  "servergroup_count_limit": 0,
  "servergroup_count_usage": 0,
  "shared_vm_count_limit": 2,
  "shared_vm_count_usage": 0,
  "snapshot_schedule_count_limit": 1,
  "snapshot_schedule_count_usage": 0,
  "subnet_count_limit": 0,
  "subnet_count_usage": 0,
  "vm_count_limit": 0,
  "vm_count_usage": 0,
  "volume_count_limit": 2,
  "volume_count_usage": 0,
  "volume_size_limit": 25,
  "volume_size_usage": 0,
  "volume_snapshots_count_limit": 2,
  "volume_snapshots_count_usage": 0,
  "volume_snapshots_size_limit": 50,
  "volume_snapshots_size_usage": 0
}
`

var CombinedResponse = fmt.Sprintf(`
{
  "global_quotas": %s,
  "regional_quotas": [
    %s
  ]
}
`, GlobalResponse, RegionalResponse)

var (
	CombinedQuota1 = quotas.CombinedQuota{
		GlobalQuotas: quotas.Quota{
			"keypair_count_limit": 100,
			"keypair_count_usage": 0,
			"project_count_limit": 2,
			"project_count_usage": 1,
		},
		RegionalQuotas: []quotas.Quota{{
			"region_id":                            1,
			"baremetal_basic_count_limit":          0,
			"baremetal_basic_count_usage":          0,
			"baremetal_hf_count_limit":             0,
			"baremetal_hf_count_usage":             0,
			"baremetal_infrastructure_count_limit": 0,
			"baremetal_infrastructure_count_usage": 0,
			"baremetal_network_count_limit":        2,
			"baremetal_network_count_usage":        0,
			"baremetal_storage_count_limit":        0,
			"baremetal_storage_count_usage":        0,
			"cpu_count_limit":                      2,
			"cpu_count_usage":                      0,
			"external_ip_count_limit":              2,
			"external_ip_count_usage":              0,
			"firewall_count_limit":                 2,
			"firewall_count_usage":                 0,
			"floating_count_limit":                 0,
			"floating_count_usage":                 0,
			"gpu_count_limit":                      0,
			"gpu_count_usage":                      0,
			"image_count_limit":                    2,
			"image_count_usage":                    0,
			"image_size_limit":                     26843545600,
			"image_size_usage":                     0,
			"loadbalancer_count_limit":             0,
			"loadbalancer_count_usage":             0,
			"network_count_limit":                  0,
			"network_count_usage":                  0,
			"ram_limit":                            4096,
			"ram_usage":                            0,
			"router_count_limit":                   0,
			"router_count_usage":                   0,
			"secret_count_limit":                   10,
			"secret_count_usage":                   0,
			"servergroup_count_limit":              0,
			"servergroup_count_usage":              0,
			"shared_vm_count_limit":                2,
			"shared_vm_count_usage":                0,
			"snapshot_schedule_count_limit":        1,
			"snapshot_schedule_count_usage":        0,
			"subnet_count_limit":                   0,
			"subnet_count_usage":                   0,
			"vm_count_limit":                       0,
			"vm_count_usage":                       0,
			"volume_count_limit":                   2,
			"volume_count_usage":                   0,
			"volume_size_limit":                    25,
			"volume_size_usage":                    0,
			"volume_snapshots_count_limit":         2,
			"volume_snapshots_count_usage":         0,
			"volume_snapshots_size_limit":          50,
			"volume_snapshots_size_usage":          0,
		}},
	}
	clientID = 3
	regionID = 1
)
