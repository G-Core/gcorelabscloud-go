package limits

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"net/http"
)

// Limit represents a limit structure.
type GlobalLimits struct {
	KeypairCountLimit int `json:"keypair_count_limit" validate:"gte=-1"`
	ProjectCountLimit int `json:"project_count_limit" validate:"gte=-1"`
}
type RegionalLimits struct {
	BaremetalBasicCountLimit          int `json:"baremetal_basic_count_limit" validate:"gte=-1"`
	BaremetalHFCountLimit             int `json:"baremetal_hf_count_limit" validate:"gte=-1"`
	BaremetalInfrastructureCountLimit int `json:"baremetal_infrastructure_count_limit" validate:"gte=-1"`
	BaremetalNetworkCountLimit        int `json:"baremetal_network_count_limit" validate:"gte=-1"`
	BaremetalStorageCountLimit        int `json:"baremetal_storage_count_limit" validate:"gte=-1"`
	ClusterCountLimit                 int `json:"cluster_count_limit" validate:"gte=-1"`
	CPUCountLimit                     int `json:"cpu_count_limit" validate:"gte=-1"`
	ExternalIPCountLimit              int `json:"external_ip_count_limit" validate:"gte=-1"`
	FirewallCountLimit                int `json:"firewall_count_limit" validate:"gte=-1"`
	FloatingCountLimit                int `json:"floating_count_limit" validate:"gte=-1"`
	GPUCountLimit                     int `json:"gpu_count_limit" validate:"gte=-1"`
	ImageCountLimit                   int `json:"image_count_limit" validate:"gte=-1"`
	ImageSizeLimit                    int `json:"image_size_limit" validate:"gte=-1"`
	LoadbalancerCountLimit            int `json:"loadbalancer_count_limit" validate:"gte=-1"`
	NetworkCountLimit                 int `json:"network_count_limit" validate:"gte=-1"`
	RAMLimit                          int `json:"ram_limit" validate:"gte=-1"`
	RouterCountLimit                  int `json:"router_count_limit" validate:"gte=-1"`
	SecretCountLimit                  int `json:"secret_count_limit" validate:"gte=-1"`
	ServergroupCountLimit             int `json:"servergroup_count_limit" validate:"gte=-1"`
	SharedVMCountLimit                int `json:"shared_vm_count_limit" validate:"gte=-1"`
	SnapshotScheduleCountLimit        int `json:"snapshot_schedule_count_limit" validate:"gte=-1"`
	SubnetCountLimit                  int `json:"subnet_count_limit" validate:"gte=-1"`
	VMCountLimit                      int `json:"vm_count_limit" validate:"gte=-1"`
	VolumeCountLimit                  int `json:"volume_count_limit" validate:"gte=-1"`
	VolumeSizeLimit                   int `json:"volume_size_limit" validate:"gte=-1"`
	VolumeSnapshotsCountLimit         int `json:"volume_snapshots_count_limit" validate:"gte=-1"`
	VolumeSnapshotsSizeLimit          int `json:"volume_snapshots_size_limit" validate:"gte=-1"`
}
type Limit struct {
	GlobalLimits   GlobalLimits     `json:"global_limits"`
	RegionalLimits []RegionalLimits `json:"regional_limits"`
}

// Get retrieves a specific quota based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// Delete deleted limit request
func Delete(c *gcorecloud.ServiceClient, id int) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, id), &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
