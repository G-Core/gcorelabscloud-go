package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
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

// ListClustersOptsBuilder allows extensions to add additional parameters to the List request.
type ListClustersOptsBuilder interface {
	ToListClustersQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Limit  int `q:"limit" validate:"omitempty,gt=0"`
	Offset int `q:"offset" validate:"omitempty,gte=0"`
}

// ToListClustersQuery formats a ListOpts into a query string.
func (opts ListOpts) ToListClustersQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List retrieves list of GPU flavors
func List(client *gcorecloud.ServiceClient, opts ListClustersOptsBuilder) pagination.Pager {
	url := ClustersURL(client)
	if opts != nil {
		query, err := opts.ToListClustersQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves all GPU clusters, using the provided ListClustersOptsBuilder to filter results.
func ListAll(client *gcorecloud.ServiceClient, opts ListClustersOptsBuilder) ([]Cluster, error) {
	allPages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusters(allPages)
}

type ServerCredentialsOpts struct {
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	SSHKeyName string `json:"ssh_key_name,omitempty"`
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
	Name     *string      `json:"name,omitempty"`
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
	ServersCount    int                `json:"servers_count" validate:"required"`
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

// ClusterActionsOpts allows extensions to add parameters to cluster actions.
type ClusterActionsOpts struct {
	Action       ClusterAction     `json:"action" required:"true" validate:"required,enum"`
	ServersCount int               `json:"servers_count,omitempty" validate:"rfe=Action:resize,gt=-1"`
	UpdateTags   map[string]string `json:"tags,omitempty" validate:"rfe=Action:update_tags"`
}

// Validate checks if the provided options are valid.
func (opts ClusterActionsOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToClusterActionsMap builds a request body from ClusterActionsOpts.
func (opts ClusterActionsOpts) ToClusterActionsMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// ClusterActionsOptsBuilder allows extensions to add additional parameters to cluster actions.
type ClusterActionsOptsBuilder interface {
	ToClusterActionsMap() (map[string]interface{}, error)
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

func ApplyAction(client *gcorecloud.ServiceClient, clusterID string, opts ClusterActionsOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterActionsMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterActionURL(client, clusterID)
	_, r.Err = client.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

func Resize(client *gcorecloud.ServiceClient, clusterID string, serversCount int) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action:       ResizeClusterAction,
		ServersCount: serversCount,
	}
	return ApplyAction(client, clusterID, opts)
}

func UpdateTags(client *gcorecloud.ServiceClient, clusterID string, tags map[string]string) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action:     UpdateTagsClusterAction,
		UpdateTags: tags,
	}
	return ApplyAction(client, clusterID, opts)
}

func SoftReboot(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action: SoftRebootClusterAction,
	}
	return ApplyAction(client, clusterID, opts)
}

func HardReboot(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action: HardRebootClusterAction,
	}
	return ApplyAction(client, clusterID, opts)
}

func Start(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action: StartClusterAction,
	}
	return ApplyAction(client, clusterID, opts)
}

func Stop(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	opts := ClusterActionsOpts{
		Action: StopClusterAction,
	}
	return ApplyAction(client, clusterID, opts)
}
