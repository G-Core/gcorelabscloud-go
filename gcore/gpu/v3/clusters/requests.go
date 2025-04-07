package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"net/http"
)

// RenameClusterOptsBuilder allows extensions to add parameters to rename cluster options.
type RenameClusterOptsBuilder interface {
	ToRenameClusterActionMap() (map[string]interface{}, error)
}

// RenameClusterOpts specifies the parameters for the Rename method.
type RenameClusterOpts struct {
	Name string `json:"name" validate:"required"`
}

// Validate checks if the provided options are valid.
func (opts RenameClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToRenameClusterActionMap builds a request body from RenameInstanceOpts.
func (opts RenameClusterOpts) ToRenameClusterActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

type ServerCredentialsOpts struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	KeypairName string `json:"keypair_name,omitempty"`
}

type ServerSettingsOpts struct {
	Interfaces     []InterfaceOpts        `json:"interfaces"`
	SecurityGroups []string               `json:"security_groups,omitempty"`
	Volumes        []VolumeOpts           `json:"volumes"`
	UserData       *string                `json:"user_data,omitempty"`
	Credentials    *ServerCredentialsOpts `json:"credentials,omitempty"`
}

// VolumeOpts represents options used to create a volume.
type VolumeOpts struct {
	Source              VolumeSource      `json:"source" validate:"required,enum"`
	BootIndex           int               `json:"boot_index" validate:"required"`
	DeleteOnTermination bool              `json:"delete_on_termination,omitempty"`
	Name                string            `json:"name" validate:"required"`
	Size                int               `json:"size,omitempty" validate:"required"`
	ImageID             string            `json:"image_id,omitempty" validate:"rfe=Source:image,allowed_without=SnapshotID,omitempty,uuid4"`
	SnapshotID          string            `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,allowed_without=ImageID,omitempty,uuid4"`
	Tags                map[string]string `json:"tags,omitempty"`
	Type                VolumeType        `json:"type,omitempty" validate:"required,enum"`
}

type InterfaceOpts interface {
	implInterfaceOpts()
}

func (ExternalInterfaceOpts) implInterfaceOpts()  {}
func (SubnetInterfaceOpts) implInterfaceOpts()    {}
func (AnySubnetInterfaceOpts) implInterfaceOpts() {}

type ExternalInterfaceOpts struct {
	Name     *string      `json:"name"`
	Type     string       `json:"type" validate:"required"`
	IPFamily IPFamilyType `json:"ip_family,omitempty"`
}

type FloatingIPOpts struct {
	Source string `json:"source" validate:"required,enum"`
}

type SubnetInterfaceOpts struct {
	NetworkID  string          `json:"network_id" validate:"required"`
	Name       *string         `json:"name,omitempty"`
	Type       string          `json:"type" validate:"required"`
	SubnetID   string          `json:"subnet_id" validate:"required"`
	FloatingIP *FloatingIPOpts `json:"floating_ip,omitempty"`
}

type AnySubnetInterfaceOpts struct {
	NetworkID  string          `json:"network_id" validate:"required"`
	Name       *string         `json:"name,omitempty"`
	Type       string          `json:"type" validate:"required"`
	IPFamily   IPFamilyType    `json:"ip_family,omitempty"`
	IPAddress  *string         `json:"ip_address,omitempty"`
	FloatingIP *FloatingIPOpts `json:"floating_ip,omitempty"`
}

// CreateClusterOpts allows extensions to add parameters to create cluster options.
type CreateClusterOpts struct {
	Name            string             `json:"name" validate:"required"`
	Flavor          string             `json:"flavor" validate:"required"`
	Tags            map[string]string  `json:"tags,omitempty"`
	ServersCount    int                `json:"servers_count,omitempty"`
	ServersSettings ServerSettingsOpts `json:"servers_settings,omitempty"`
}

// Validate checks if the provided options are valid.
func (opts CreateClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

func (opts CreateClusterOpts) ToCreateClusterMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// CreateClusterOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateClusterOptsBuilder interface {
	ToCreateClusterMap() (map[string]interface{}, error)
}

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// Delete removes a specific GPU cluster by its ID.
func Delete(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}

// Rename changes the name of a GPU cluster.
func Rename(client *gcorecloud.ServiceClient, clusterID string, opts RenameClusterOptsBuilder) (r GetResult) {
	b, err := opts.ToRenameClusterActionMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterURL(client, clusterID)
	_, r.Err = client.Patch(url, b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}

func Create(client *gcorecloud.ServiceClient, opts CreateClusterOptsBuilder) (r tasks.Result) {
	b, err := opts.ToCreateClusterMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClustersURL(client)
	_, r.Err = client.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}
