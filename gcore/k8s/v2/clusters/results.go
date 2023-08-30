package clusters

import (
	"fmt"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
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

// Cluster represents a cluster structure.
type Cluster struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	KeyPair       string              `json:"keypair"`
	NodeCount     int                 `json:"node_count"`
	FlavorID      string              `json:"flavor_id"`
	Status        string              `json:"status"`
	Pools         []pools.ClusterPool `json:"pools"`
	Version       string              `json:"version"`
	IsPublic      bool                `json:"is_public"`
	FixedNetwork  string              `json:"fixed_network"`
	FixedSubnet   string              `json:"fixed_subnet"`
	CreatedAt     time.Time           `json:"created_at"`
	CreatorTaskID string              `json:"creator_task_id"`
	TaskID        string              `json:"task_id"`
	ProjectID     int                 `json:"project_id"`
	RegionID      int                 `json:"region_id"`
	Region        string              `json:"region"`
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

type ClusterTaskResult struct {
	K8sClusters []string `json:"k8s_clusters" mapstructure:"k8s_clusters"`
}

func ExtractClusterIDFromTask(task *tasks.Task) (string, error) {
	var result ClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	if len(result.K8sClusters) == 0 {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	return result.K8sClusters[0], nil
}
