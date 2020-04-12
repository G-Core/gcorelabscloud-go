package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/quota/v1/quotas"
)

const GetResponse = `
{
  "cpu_count_usage": 0,
  "cpu_count_limit": 1,
  "floating_count_usage": 0,
  "floating_count_limit": 1,
  "project_count_usage": 0,
  "project_count_limit": 1,
  "vm_count_usage": 0,
  "vm_count_limit": 1,
  "volume_size_usage": 0,
  "volume_size_limit": 1,
  "image_count_usage": 0,
  "image_count_limit": 1,
  "router_count_usage": 0,
  "router_count_limit": 1,
  "firewall_count_usage": 0,
  "firewall_count_limit": 1,
  "subnet_count_usage": 0,
  "subnet_count_limit": 1,
  "network_count_usage": 0,
  "network_count_limit": 1,
  "image_size_usage": 0,
  "image_size_limit": 1,
  "ram_usage": 0,
  "ram_limit": 1,
  "volume_snapshots_count_usage": 0,
  "volume_snapshots_count_limit": 1,
  "volume_snapshots_size_usage": 0,
  "volume_snapshots_size_limit": 1,
  "volume_count_usage": 0,
  "volume_count_limit": 1,
  "loadbalancer_count_usage": 0,
  "loadbalancer_count_limit": 1,
  "external_ip_count_usage": 0,
  "external_ip_count_limit": 1
}
`

const ReplaceRequest = GetResponse

const UpdateRequest = `
{
  "external_ip_count_usage": 4
}	
`

const ReplaceResponse = GetResponse
const UpdateResponse = GetResponse

var (
	Quota1 = quotas.Quota{
		ProjectCountLimit:         1,
		VMCountLimit:              1,
		CPUCountLimit:             1,
		RAMLimit:                  1,
		VolumeCountLimit:          1,
		VolumeSizeLimit:           1,
		VolumeSnapshotsCountLimit: 1,
		VolumeSnapshotsSizeLimit:  1,
		ImageCountLimit:           1,
		ImageSizeLimit:            1,
		NetworkCountLimit:         1,
		SubnetCountLimit:          1,
		FloatingCountLimit:        1,
		RouterCountLimit:          1,
		FirewallCountLimit:        1,
		LoadbalancerCountLimit:    1,
		ExternalIPCountLimit:      1,
		ProjectCountUsage:         0,
		VMCountUsage:              0,
		CPUCountUsage:             0,
		RAMUsage:                  0,
		VolumeCountUsage:          0,
		VolumeSizeUsage:           0,
		VolumeSnapshotsCountUsage: 0,
		VolumeSnapshotsSizeUsage:  0,
		ImageCountUsage:           0,
		ImageSizeUsage:            0,
		NetworkCountUsage:         0,
		SubnetCountUsage:          0,
		FloatingCountUsage:        0,
		RouterCountUsage:          0,
		FirewallCountUsage:        0,
		LoadbalancerCountUsage:    0,
		ExternalIPCountUsage:      0,
	}
)
