package clusters

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CertificateResult represents the result of a get certificate operation.
// Call its Extract method to interpret it as a Certificate.
type CertificateResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster CA certificate.
func (r CertificateResult) Extract() (*Certificate, error) {
	var c Certificate
	err := r.ExtractInto(&c)
	return &c, err
}

// ConfigResult represents the result of a get config operation. Call its Extract
// method to interpret it as a Config.
type ConfigResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster config.
func (r ConfigResult) Extract() (*Config, error) {
	var c Config
	err := r.ExtractInto(&c)
	return &c, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}

type Authentication struct {
	OIDC *OIDC `json:"oidc,omitempty"`
}

type OIDC struct {
	ClientID       string            `json:"client_id,omitempty"`
	GroupsClaim    string            `json:"groups_claim,omitempty"`
	GroupsPrefix   string            `json:"groups_prefix,omitempty"`
	IssuerURL      string            `json:"issuer_url,omitempty"`
	RequiredClaims map[string]string `json:"required_claims,omitempty"`
	SigningAlgs    []string          `json:"signing_algs,omitempty"`
	UsernameClaim  string            `json:"username_claim,omitempty"`
	UsernamePrefix string            `json:"username_prefix,omitempty"`
}

type Cilium struct {
	MaskSize                 int             `json:"mask_size,omitempty"`
	MaskSizeV6               int             `json:"mask_size_v6,omitempty"`
	Tunnel                   TunnelType      `json:"tunnel,omitempty"`
	Encryption               bool            `json:"encryption"`
	LoadBalancerMode         LBModeType      `json:"lb_mode,omitempty"`
	LoadBalancerAcceleration bool            `json:"lb_acceleration"`
	RoutingMode              RoutingModeType `json:"routing_mode,omitempty"`
	HubbleRelay              bool            `json:"hubble_relay"`
	HubbleUI                 bool            `json:"hubble_ui"`
}

type CNI struct {
	Provider CNIProvider `json:"provider"`
	Cilium   *Cilium     `json:"cilium,omitempty"`
}

type DDosProfile struct {
	Enabled             bool               `json:"enabled,omitempty"`
	Fields              []DDosProfileField `json:"fields,omitempty"`
	ProfileTemplate     *int               `json:"profile_template,omitempty"`
	ProfileTemplateName *string            `json:"profile_template_name,omitempty"`
}

// Cluster represents a cluster structure.
type Cluster struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	KeyPair          string              `json:"keypair"`
	NodeCount        int                 `json:"node_count"`
	FlavorID         string              `json:"flavor_id"`
	Status           string              `json:"status"`
	Pools            []pools.ClusterPool `json:"pools"`
	Version          string              `json:"version"`
	IsPublic         bool                `json:"is_public"`
	Authentication   *Authentication     `json:"authentication,omitempty"`
	AutoscalerConfig map[string]string   `json:"autoscaler_config,omitempty"`
	CNI              *CNI                `json:"cni,omitempty"`
	FixedNetwork     string              `json:"fixed_network"`
	FixedSubnet      string              `json:"fixed_subnet"`
	PodsIPPool       *gcorecloud.CIDR    `json:"pods_ip_pool,omitempty" validate:"omitempty,cidr"`
	ServicesIPPool   *gcorecloud.CIDR    `json:"services_ip_pool,omitempty" validate:"omitempty,cidr"`
	PodsIPV6Pool     *gcorecloud.CIDR    `json:"pods_ipv6_pool,omitempty" validate:"omitempty,cidr"`
	ServicesIPV6Pool *gcorecloud.CIDR    `json:"services_ipv6_pool,omitempty" validate:"omitempty,cidr"`
	IsIPV6           bool                `json:"is_ipv6,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	CreatorTaskID    string              `json:"creator_task_id"`
	DDosProfile      *DDosProfile        `json:"ddos_profile,omitempty"`
	TaskID           string              `json:"task_id"`
	ProjectID        int                 `json:"project_id"`
	RegionID         int                 `json:"region_id"`
	Region           string              `json:"region"`
}

// Certificate represents a cluster CA certificate.
type Certificate struct {
	Key         string `json:"key"`
	Certificate string `json:"certificate"`
}

// Config represents a kubeconfig structure.
type Config struct {
	Config string `json:"config"`
}

// Version represents a cluster version structure.
type Version struct {
	Version string `json:"version"`
}

// ClusterPage is the page returned by a pager when traversing over a collection of clusters.
type ClusterPage struct {
	pagination.LinkedPageBase
}

// VersionPage is the page returned by a pager when traversing over
// a collection of cluster versions.
type VersionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r VersionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ClusterPage struct is empty.
func (r ClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a VersionPage struct is empty.
func (r VersionPage) IsEmpty() (bool, error) {
	is, err := ExtractVersions(r)
	return len(is) == 0, err
}

// ExtractCluster accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of ClusterListWithPool structs.
// In other words, a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]Cluster, error) {
	var s []Cluster
	err := ExtractClustersInto(r, &s)
	return s, err
}

// ExtractVersions accepts a Page struct, specifically a VersionPage struct,
// and extracts the elements into a slice of Version structs..
// In other words, a generic collection is mapped into a relevant slice.
func ExtractVersions(r pagination.Page) ([]Version, error) {
	var s []Version
	err := ExtractVersionInto(r, &s)
	return s, err
}

func ExtractClustersInto(r pagination.Page, v interface{}) error {
	return r.(ClusterPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractVersionInto(r pagination.Page, v interface{}) error {
	return r.(VersionPage).Result.ExtractIntoSlicePtr(v, "results")
}
