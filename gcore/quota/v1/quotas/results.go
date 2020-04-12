package quotas

import (
	"reflect"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

const Sentinel = -1

type commonResult struct {
	gcorecloud.Result
}

// Quota represents a quota structure.
type Quota struct {
	ProjectCountLimit         int `json:"project_count_limit"`
	VMCountLimit              int `json:"vm_count_limit"`
	CPUCountLimit             int `json:"cpu_count_limit"`
	RAMLimit                  int `json:"ram_limit"`
	VolumeCountLimit          int `json:"volume_count_limit"`
	VolumeSizeLimit           int `json:"volume_size_limit"`
	VolumeSnapshotsCountLimit int `json:"volume_snapshots_count_limit"`
	VolumeSnapshotsSizeLimit  int `json:"volume_snapshots_size_limit"`
	ImageCountLimit           int `json:"image_count_limit"`
	ImageSizeLimit            int `json:"image_size_limit"`
	NetworkCountLimit         int `json:"network_count_limit"`
	SubnetCountLimit          int `json:"subnet_count_limit"`
	FloatingCountLimit        int `json:"floating_count_limit"`
	RouterCountLimit          int `json:"router_count_limit"`
	FirewallCountLimit        int `json:"firewall_count_limit"`
	LoadbalancerCountLimit    int `json:"loadbalancer_count_limit"`
	ExternalIPCountLimit      int `json:"external_ip_count_limit"`
	ProjectCountUsage         int `json:"project_count_usage"`
	VMCountUsage              int `json:"vm_count_usage"`
	CPUCountUsage             int `json:"cpu_count_usage"`
	RAMUsage                  int `json:"ram_usage"`
	VolumeCountUsage          int `json:"volume_count_usage"`
	VolumeSizeUsage           int `json:"volume_size_usage"`
	VolumeSnapshotsCountUsage int `json:"volume_snapshots_count_usage"`
	VolumeSnapshotsSizeUsage  int `json:"volume_snapshots_size_usage"`
	ImageCountUsage           int `json:"image_count_usage"`
	ImageSizeUsage            int `json:"image_size_usage"`
	NetworkCountUsage         int `json:"network_count_usage"`
	SubnetCountUsage          int `json:"subnet_count_usage"`
	FloatingCountUsage        int `json:"floating_count_usage"`
	RouterCountUsage          int `json:"router_count_usage"`
	FirewallCountUsage        int `json:"firewall_count_usage"`
	LoadbalancerCountUsage    int `json:"loadbalancer_count_usage"`
	ExternalIPCountUsage      int `json:"external_ip_count_usage"`
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
