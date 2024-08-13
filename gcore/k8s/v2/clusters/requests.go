package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/quota/v2/quotas"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters
// to the CheckLimits request.

type CheckLimitsOptsBuilder interface {
	ToCheckLimitsMap() (map[string]interface{}, error)
}

type CheckLimitsOpts struct {
	Pools []CheckLimitsPoolOpts `json:"pools,omitempty"`
}

type CheckLimitsPoolOpts struct {
	Name              string                         `json:"name,omitempty" validate:"omitempty"`
	FlavorID          string                         `json:"flavor_id" required:"true" validate:"required"`
	MinNodeCount      int                            `json:"min_node_count,omitempty" validate:"omitempty,gt=0,ltefield=MaxNodeCount"`
	MaxNodeCount      int                            `json:"max_node_count,omitempty" validate:"omitempty,gt=0,gtefield=MinNodeCount"`
	NodeCount         int                            `json:"node_count,omitempty" validate:"omitempty"`
	BootVolumeSize    int                            `json:"boot_volume_size,omitempty" validate:"omitempty,gt=0"`
	ServerGroupPolicy servergroups.ServerGroupPolicy `json:"servergroup_policy,omitempty" validate:"omitempty,enum"`
}

// ToCheckLimitsMap builds a request body from CheckLimitsOpts.
func (opts CheckLimitsOpts) ToCheckLimitsMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateOptsBuilder allows extensions to add additional parameters
// to the Create request.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

type AuthenticationCreateOpts struct {
	OIDC *OIDCCreateOpts `json:"oidc,omitempty" validate:"omitempty"`
}

type OIDCCreateOpts struct {
	ClientID       string            `json:"client_id,omitempty" validate:"omitempty"`
	GroupsClaim    string            `json:"groups_claim,omitempty" validate:"omitempty"`
	GroupsPrefix   string            `json:"groups_prefix,omitempty" validate:"omitempty"`
	IssuerURL      string            `json:"issuer_url,omitempty" validate:"omitempty"`
	RequiredClaims map[string]string `json:"required_claims,omitempty" validate:"omitempty"`
	SigningAlgs    []string          `json:"signing_algs,omitempty" validate:"omitempty"`
	UsernameClaim  string            `json:"username_claim,omitempty" validate:"omitempty"`
	UsernamePrefix string            `json:"username_prefix,omitempty" validate:"omitempty"`
}

type CiliumCreateOpts struct {
	MaskSize                 int             `json:"mask_size,omitempty"`
	MaskSizeV6               int             `json:"mask_size_v6,omitempty"`
	Tunnel                   TunnelType      `json:"tunnel"`
	Encryption               bool            `json:"encryption"`
	LoadBalancerMode         LBModeType      `json:"lb_mode,omitempty"`
	LoadBalancerAcceleration bool            `json:"lb_acceleration"`
	RoutingMode              RoutingModeType `json:"routing_mode,omitempty"`
	HubbleRelay              bool            `json:"hubble_relay"`
	HubbleUI                 bool            `json:"hubble_ui"`
}

type CNICreateOpts struct {
	Provider CNIProvider       `json:"provider" required:"true" validate:"required"`
	Cilium   *CiliumCreateOpts `json:"cilium,omitempty"`
}

// CreateOpts represents options used to create a cluster.
type CreateOpts struct {
	Name             string                    `json:"name" required:"true" validate:"required,gt=0,lte=20"`
	Authentication   *AuthenticationCreateOpts `json:"authentication,omitempty" validate:"omitempty"`
	AutoscalerConfig map[string]string         `json:"autoscaler_config,omitempty" validate:"omitempty"`
	CNI              *CNICreateOpts            `json:"cni,omitempty" validate:"omitempty"`
	FixedNetwork     string                    `json:"fixed_network" validate:"omitempty,uuid4"`
	FixedSubnet      string                    `json:"fixed_subnet" validate:"omitempty,uuid4"`
	PodsIPPool       *gcorecloud.CIDR          `json:"pods_ip_pool,omitempty" validate:"omitempty,cidr"`
	ServicesIPPool   *gcorecloud.CIDR          `json:"services_ip_pool,omitempty" validate:"omitempty,cidr"`
	PodsIPV6Pool     *gcorecloud.CIDR          `json:"pods_ipv6_pool,omitempty" validate:"omitempty,cidr"`
	ServicesIPV6Pool *gcorecloud.CIDR          `json:"services_ipv6_pool,omitempty" validate:"omitempty,cidr"`
	KeyPair          string                    `json:"keypair" required:"true" validate:"required"`
	Version          string                    `json:"version" required:"true" validate:"required"`
	IsIPV6           bool                      `json:"is_ipv6,omitempty"`
	Pools            []pools.CreateOpts        `json:"pools" required:"true" validate:"required,min=1,dive"`
}

// ToClusterCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

type AuthenticationUpdateOpts = AuthenticationCreateOpts

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to upgrade a cluster.
type UpdateOpts struct {
	Authentication   *AuthenticationUpdateOpts `json:"authentication,omitempty" validate:"omitempty"`
	AutoscalerConfig map[string]string         `json:"autoscaler_config,omitempty" validate:"omitempty"`
	CNI              *CNICreateOpts            `json:"cni,omitempty" validate:"omitempty"`
}

// ToClusterUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpgradeOptsBuilder allows extensions to add additional parameters to the Upgrade request.
type UpgradeOptsBuilder interface {
	ToClusterUpgradeMap() (map[string]interface{}, error)
}

// UpgradeOpts represents options used to upgrade a cluster.
type UpgradeOpts struct {
	Version string `json:"version" required:"true"`
}

// ToClusterUpgradeMap builds a request body from UpgradeOpts.
func (opts UpgradeOpts) ToClusterUpgradeMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CheckLimits checks quota limits for the values provided and returns the diff for exceeded quota.
func CheckLimits(c *gcorecloud.ServiceClient, opts CheckLimitsOptsBuilder) (r quotas.CommonResult) {
	b, err := opts.ToCheckLimitsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(checkLimitsURL(c), b, &r.Body, nil)
	return
}

// List returns a Pager which allows you to iterate over a collection of clusters.
func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListALL is a convenience function that returns all clusters.
func ListAll(c *gcorecloud.ServiceClient) ([]Cluster, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusters(page)
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

// Get retrieves a specific cluster based on its name.
func Get(c *gcorecloud.ServiceClient, clusterName string) (r GetResult) {
	url := getURL(c, clusterName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Delete accepts cluster name and deletes the cluster associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterName string) (r tasks.Result) {
	url := deleteURL(c, clusterName)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, nil)
	return
}

// GetCertificate accepts cluster name and returns the cluster CA certificate.
func GetCertificate(c *gcorecloud.ServiceClient, clusterName string) (r CertificateResult) {
	url := certificatesURL(c, clusterName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// GetConfig accepts cluster name and returns the cluster kubeconfig.
func GetConfig(c *gcorecloud.ServiceClient, clusterName string) (r ConfigResult) {
	url := configURL(c, clusterName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListInstances returns a Pager which allows you to iterate over a collection of cluster instances.
func ListInstances(c *gcorecloud.ServiceClient, clusterID string) pagination.Pager {
	url := instancesURL(c, clusterID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListInstancesAll is a convenience function that returns all cluster instances.
func ListInstancesAll(c *gcorecloud.ServiceClient, clusterID string) ([]instances.Instance, error) {
	page, err := ListInstances(c, clusterID).AllPages()
	if err != nil {
		return nil, err
	}
	return instances.ExtractInstances(page)
}

// Update accepts a UpdateOpts struct and updates an existing cluster using the values provided.
func Update(c *gcorecloud.ServiceClient, clusterID string, opts UpdateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, clusterID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
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

// Versions returns a Pager which allows you to iterate over a collection of
// supported cluster versions.
func Versions(c *gcorecloud.ServiceClient) pagination.Pager {
	url := versionsURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VersionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// VersionsAll is a convenience function that returns all supported cluster versions.
func VersionsAll(c *gcorecloud.ServiceClient) ([]Version, error) {
	page, err := Versions(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractVersions(page)
}
