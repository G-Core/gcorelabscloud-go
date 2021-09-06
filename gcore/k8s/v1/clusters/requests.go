package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// Versions request.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the clusters attributes you want to see returned. SortKey allows you to sort
// by a particular clusters attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Marker  string `q:"marker"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
	Detail  bool   `q:"detail"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// clusters. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Versions returns available cluster versions.
func Versions(c *gcorecloud.ServiceClient) pagination.Pager {
	url := versionsURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VersionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Instances returns a Pager which allows you to iterate over a collection of cluster instances.
func Instances(c *gcorecloud.ServiceClient, clusterID string) pagination.Pager {
	url := instancesURL(c, clusterID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster.
type CreateOpts struct {
	Name                      string             `json:"name" required:"true"`
	FixedNetwork              string             `json:"fixed_network" required:"true" validate:"required,uuid4"`
	FixedSubnet               string             `json:"fixed_subnet" required:"true" validate:"required,uuid4"`
	KeyPair                   string             `json:"keypair,omitempty"`
	PodsIPPool                *gcorecloud.CIDR   `json:"pods_ip_pool,omitempty"`
	ServicesIPPool            *gcorecloud.CIDR   `json:"services_ip_pool,omitempty"`
	AutoHealingEnabled        bool               `json:"auto_healing_enabled"`
	MasterLBFloatingIPEnabled bool               `json:"master_lb_floating_ip_enabled,omitempty"`
	Version                   string             `json:"version,omitempty" validate:"omitempty,sem"`
	Pools                     []pools.CreateOpts `json:"pools" required:"true" validate:"required,min=1,dive"`
}

// ToClusterCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// ResizeOptsBuilder allows extensions to add additional parameters to the Resize request.
type ResizeOptsBuilder interface {
	ToClusterResizeMap() (map[string]interface{}, error)
}

// ResizeOpts represents options used to update a cluster.
type ResizeOpts struct {
	NodeCount     int      `json:"node_count" required:"true" validate:"required,gt=0"`
	NodesToRemove []string `json:"nodes_to_remove,omitempty" validate:"omitempty,dive,uuid4"`
}

// ToClusterResizeMap builds a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterResizeMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Resize accepts a ResizeOpts struct and updates an existing cluster using the values provided.
func Resize(c *gcorecloud.ServiceClient, clusterID, poolID string, opts ResizeOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(resizeURL(c, clusterID, poolID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// UpgradeOptsBuilder allows extensions to add additional parameters to the Upgrade request.
type UpgradeOptsBuilder interface {
	ToClusterUpgradeMap() (map[string]interface{}, error)
}

// UpgradeOpts represents options used to upgrade a cluster.
type UpgradeOpts struct {
	Pool    string `json:"pool,omitempty" validate:"omitempty,uuid4"`
	Version string `json:"version" required:"true" validate:"required,sem"`
}

// ToClusterUpgradeMap builds a request body from UpgradeOpts.
func (opts UpgradeOpts) ToClusterUpgradeMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Upgrade accepts a UpgradeOpts struct and upgrades an existing cluster using the values provided.
func Upgrade(c *gcorecloud.ServiceClient, clusterID string, opts UpgradeOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterUpgradeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(upgradeURL(c, clusterID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the cluster associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, clusterID), &r.Body, nil)
	return
}

// GetConfig accepts a unique ID and get cluster k8s config.
func GetConfig(c *gcorecloud.ServiceClient, clusterID string) (r ConfigResult) {
	url := configURL(c, clusterID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListALL returns all magnum clusters.
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]ClusterListWithPool, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusters(page)
}

// InstancesAll returns all cluster instances.
func InstancesAll(c *gcorecloud.ServiceClient, clusterID string) ([]instances.Instance, error) {
	page, err := Instances(c, clusterID).AllPages()
	if err != nil {
		return nil, err
	}
	return instances.ExtractInstances(page)
}

// VersionsAll returns all cluster versions.
func VersionsAll(c *gcorecloud.ServiceClient) ([]string, error) {
	page, err := Versions(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractVersions(page)
}

// ClusterSignCertificateOptsBuilder allows extensions to add additional parameters to the SignCertificate request.
type ClusterSignCertificateOptsBuilder interface {
	ToClusterSignCertificateMap() (map[string]interface{}, error)
}

// ClusterSignCertificateOpts represents a options to sign cluster certificate.
type ClusterSignCertificateOpts struct {
	CSR string `json:"csr" validate:"required"`
}

// ToClusterSignCertificateMap builds a request body from ClusterSignCertificateOpts.
func (opts ClusterSignCertificateOpts) ToClusterSignCertificateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// SignCertificate accepts a unique ID and sign cluster certificate.
func SignCertificate(c *gcorecloud.ServiceClient, clusterID string, opts ClusterSignCertificateOptsBuilder) (r ClusterCertificateSignResult) {
	url := certificatesURL(c, clusterID)
	b, err := opts.ToClusterSignCertificateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Certificate accepts a unique ID and return cluster CA.
func Certificate(c *gcorecloud.ServiceClient, clusterID string) (r ClusterCertificateCAResult) {
	url := certificatesURL(c, clusterID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
