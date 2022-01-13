package limits

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/pagination"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const Sentinel = -1

// Limit represents a limit structure.
type Limit struct {
	//{
	KeypairCountLimit int `json:"keypair_count_limit" validate:"gte=-1"`
	ProjectCountLimit int `json:"project_count_limit" validate:"gte=-1"`
	//}

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

func NewLimit() Limit {
	return Limit{
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
	}
}

type CreateOpts struct {
	Description     string `json:"description" required:"true" validate:"required"`
	RequestedQuotas Limit  `json:"requested_quotas" required:"true" validate:"required"`
}

func NewCreateOpts(description string) CreateOpts {
	return CreateOpts{
		Description:     description,
		RequestedQuotas: NewLimit(),
	}
}

func (opts Limit) ToRequestMap() map[string]interface{} {
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

// Get retrieves a specific quota based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToLimitCreateMap() (map[string]interface{}, error)
}

// ToLimitCreateMap builds a request body from ReplaceOpts.
func (opts CreateOpts) ToLimitCreateMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}
	rm := opts.RequestedQuotas.ToRequestMap()
	if len(rm) == 0 {
		return nil, fmt.Errorf("at least one of quota fields should be set")
	}
	m, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	m["requested_quotas"] = rm
	return m, nil
}

func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// Create accepts a ReplaceOpts struct and creates a new quota using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLimitCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, listURL(c), func(r pagination.PageResult) pagination.Page {
		return LimitResultPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListAll(c *gcorecloud.ServiceClient) ([]LimitResponse, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractLimitResults(page)
}

// Delete deleted limit request
func Delete(c *gcorecloud.ServiceClient, id int) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, id), nil)
	return
}
