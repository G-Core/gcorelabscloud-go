package quotas

import (
	"reflect"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const Sentinel = -1

type commonResult struct {
	gcorecloud.Result
}

// Quota represents a quota structure.
type Quota struct {
	ProjectCountLimit         int `json:"project_count_limit" validate:"gte=-1"`
	VMCountLimit              int `json:"vm_count_limit" validate:"gte=-1"`
	CPUCountLimit             int `json:"cpu_count_limit" validate:"gte=-1"`
	RAMLimit                  int `json:"ram_limit" validate:"gte=-1"`
	VolumeCountLimit          int `json:"volume_count_limit" validate:"gte=-1"`
	VolumeSizeLimit           int `json:"volume_size_limit" validate:"gte=-1"`
	VolumeSnapshotsCountLimit int `json:"volume_snapshots_count_limit" validate:"gte=-1"`
	VolumeSnapshotsSizeLimit  int `json:"volume_snapshots_size_limit" validate:"gte=-1"`
	ImageCountLimit           int `json:"image_count_limit" validate:"gte=-1"`
	ImageSizeLimit            int `json:"image_size_limit" validate:"gte=-1"`
	NetworkCountLimit         int `json:"network_count_limit" validate:"gte=-1"`
	SubnetCountLimit          int `json:"subnet_count_limit" validate:"gte=-1"`
	FloatingCountLimit        int `json:"floating_count_limit" validate:"gte=-1"`
	RouterCountLimit          int `json:"router_count_limit" validate:"gte=-1"`
	FirewallCountLimit        int `json:"firewall_count_limit" validate:"gte=-1"`
	LoadbalancerCountLimit    int `json:"loadbalancer_count_limit" validate:"gte=-1"`
	ExternalIPCountLimit      int `json:"external_ip_count_limit" validate:"gte=-1"`
	ClusterCountLimit         int `json:"cluster_count_limit" validate:"gte=-1"`
	ProjectCountUsage         int `json:"project_count_usage" validate:"gte=-1"`
	VMCountUsage              int `json:"vm_count_usage" validate:"gte=-1"`
	CPUCountUsage             int `json:"cpu_count_usage" validate:"gte=-1"`
	RAMUsage                  int `json:"ram_usage" validate:"gte=-1"`
	VolumeCountUsage          int `json:"volume_count_usage" validate:"gte=-1"`
	VolumeSizeUsage           int `json:"volume_size_usage" validate:"gte=-1"`
	VolumeSnapshotsCountUsage int `json:"volume_snapshots_count_usage" validate:"gte=-1"`
	VolumeSnapshotsSizeUsage  int `json:"volume_snapshots_size_usage" validate:"gte=-1"`
	ImageCountUsage           int `json:"image_count_usage" validate:"gte=-1"`
	ImageSizeUsage            int `json:"image_size_usage" validate:"gte=-1"`
	NetworkCountUsage         int `json:"network_count_usage" validate:"gte=-1"`
	SubnetCountUsage          int `json:"subnet_count_usage" validate:"gte=-1"`
	FloatingCountUsage        int `json:"floating_count_usage" validate:"gte=-1"`
	RouterCountUsage          int `json:"router_count_usage" validate:"gte=-1"`
	FirewallCountUsage        int `json:"firewall_count_usage" validate:"gte=-1"`
	LoadbalancerCountUsage    int `json:"loadbalancer_count_usage" validate:"gte=-1"`
	ExternalIPCountUsage      int `json:"external_ip_count_usage" validate:"gte=-1"`
	ClusterCountUsage         int `json:"cluster_count_usage" validate:"gte=-1"`
}

func NewQuota() Quota {
	return Quota{
		ProjectCountLimit:         Sentinel,
		VMCountLimit:              Sentinel,
		CPUCountLimit:             Sentinel,
		RAMLimit:                  Sentinel,
		VolumeCountLimit:          Sentinel,
		VolumeSizeLimit:           Sentinel,
		VolumeSnapshotsCountLimit: Sentinel,
		VolumeSnapshotsSizeLimit:  Sentinel,
		ImageCountLimit:           Sentinel,
		ImageSizeLimit:            Sentinel,
		NetworkCountLimit:         Sentinel,
		SubnetCountLimit:          Sentinel,
		FloatingCountLimit:        Sentinel,
		RouterCountLimit:          Sentinel,
		FirewallCountLimit:        Sentinel,
		LoadbalancerCountLimit:    Sentinel,
		ExternalIPCountLimit:      Sentinel,
		ClusterCountLimit:         Sentinel,
		ProjectCountUsage:         Sentinel,
		VMCountUsage:              Sentinel,
		CPUCountUsage:             Sentinel,
		RAMUsage:                  Sentinel,
		VolumeCountUsage:          Sentinel,
		VolumeSizeUsage:           Sentinel,
		VolumeSnapshotsCountUsage: Sentinel,
		VolumeSnapshotsSizeUsage:  Sentinel,
		ImageCountUsage:           Sentinel,
		ImageSizeUsage:            Sentinel,
		NetworkCountUsage:         Sentinel,
		SubnetCountUsage:          Sentinel,
		FloatingCountUsage:        Sentinel,
		RouterCountUsage:          Sentinel,
		FirewallCountUsage:        Sentinel,
		LoadbalancerCountUsage:    Sentinel,
		ExternalIPCountUsage:      Sentinel,
		ClusterCountUsage:         Sentinel,
	}
}

func (opts Quota) ToRequestMap() map[string]interface{} {
	optsValue := reflect.ValueOf(opts)
	optsType := reflect.TypeOf(opts)
	optsMap := make(map[string]interface{})
	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		jsonTag := f.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		value := int(v.Int())
		if value != Sentinel {
			optsMap[name] = value
		}
	}
	return optsMap
}

// Extract is a function that accepts a result and extracts a quota resource.
func (r commonResult) Extract() (*Quota, error) {
	var s Quota
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Quota.
type UpdateResult struct {
	commonResult
}
