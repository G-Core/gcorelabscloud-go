package limits

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// Limit represents a GlobalLimits structure.
type GlobalLimits struct {
	KeypairCountLimit int `json:"keypair_count_limit" validate:"gte=-1"`
	ProjectCountLimit int `json:"project_count_limit" validate:"gte=-1"`
}

func (g *GlobalLimits) Update(source interface{}) error {
	vGlobal := reflect.ValueOf(source).Elem()
	for i := 0; i < vGlobal.NumField(); i++ {
		field := vGlobal.Type().Field(i).Name
		value := vGlobal.Field(i).Int()
		r := reflect.ValueOf(g).Elem()
		tf := r.FieldByName(field)
		if tf.IsValid() && tf.CanSet() && tf.Kind() == reflect.Int {
			tf.SetInt(value)
		} else {
			return fmt.Errorf("cannot set global field %s", field)
		}
	}
	return nil
}

type RegionalLimits struct {
	RegionID                          int `json:"region_id" validate:"gte=0"`
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

func (r *RegionalLimits) Update(source interface{}) error {
	targetRegionalEl := reflect.ValueOf(&r).Elem().Elem()
	valueG := reflect.ValueOf(&source).Elem().Elem()
	for i := 0; i < valueG.NumField(); i++ {
		f := valueG.Type().Field(i).Name
		v := valueG.Field(i).Int()

		sf := targetRegionalEl.FieldByName(f)
		if sf.IsValid() && sf.CanSet() && sf.Kind() == reflect.Int {
			sf.SetInt(v)
		} else {
			return fmt.Errorf("cannot set regional field %s", f)
		}
	}
	return nil
}

type Limit struct {
	GlobalLimits   GlobalLimits     `json:"global_limits"`
	RegionalLimits []RegionalLimits `json:"regional_limits"`
}

const Sentinel = -1

func NewLimit() Limit {
	return Limit{
		GlobalLimits:   GlobalLimits{KeypairCountLimit: Sentinel, ProjectCountLimit: Sentinel},
		RegionalLimits: make([]RegionalLimits, 0),
	}
}

type CreateOpts struct {
	Description     string `json:"description" required:"true" validate:"required"`
	RequestedQuotas Limit  `json:"requested_limits" required:"true" validate:"required"`
}

func NewCreateOpts(description string) CreateOpts {
	return CreateOpts{
		Description:     description,
		RequestedQuotas: NewLimit(),
	}
}

// Create accepts a CreateOptsBuilder struct and creates a new quota using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLimitCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
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
	m["requested_limits"] = rm
	return m, nil
}

func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

func ToRequestMap(opts interface{}) map[string]interface{} {
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

func (opts Limit) ToRequestMap() map[string]interface{} {
	optsMap := make(map[string]interface{})
	optsMap["global_limits"] = ToRequestMap(opts.GlobalLimits)
	optsRegionalLimits := make([]map[string]interface{}, 0)
	for _, regionalItem := range opts.RegionalLimits {
		regionalMap := ToRequestMap(regionalItem)
		optsRegionalLimits = append(optsRegionalLimits, regionalMap)
	}
	optsMap["regional_limits"] = optsRegionalLimits
	return optsMap
}

// Get retrieves a specific quota based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
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
	_, r.Err = c.Delete(deleteURL(c, id), &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
