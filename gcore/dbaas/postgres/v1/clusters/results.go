package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a PostgreSQL cluster resource.
func (r commonResult) Extract() (*PostgresSQLCluster, error) {
	var s PostgresSQLCluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a PG Cluster.
type GetResult struct {
	commonResult
}

// ClusterPage is the page returned by a pager when traversing over a collection of PG clusters.
type ClusterPage struct {
	pagination.OffsetPageBase
}

// ExtractClusters accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of PG cluster structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]PostgresSQLClusterShort, error) {
	var s []PostgresSQLClusterShort
	err := ExtractClustersInto(r, &s)
	return s, err
}

func ExtractClustersInto(r pagination.Page, v interface{}) error {
	return r.(ClusterPage).Result.ExtractIntoSlicePtr(v, "results")
}

type DatabaseOverview struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Size  int    `json:"size"`
}

type Flavor struct {
	CPU       int `json:"cpu"`
	MemoryGiB int `json:"memory_gib"`
}

type HighAvailability struct {
	ReplicationMode HighAvailabilityReplicationMode `json:"replication_mode"`
}

type Network struct {
	ACL              []string `json:"acl"`
	ConnectionString string   `json:"connection_string"`
	Host             string   `json:"host"`
	NetworkType      string   `json:"network_type"`
}

type Pooler struct {
	Mode PoolerMode `json:"mode"`
	Type string     `json:"type"`
}

type PGServerConfiguration struct {
	PGConf  string  `json:"pg_conf"`
	Version string  `json:"version"`
	Pooler  *Pooler `json:"pooler"`
}

type PGStorageConfiguration struct {
	SizeGiB int    `json:"size_gib"`
	Type    string `json:"type"`
}

type PgUserOverview struct {
	IsSecretRevealed bool            `json:"is_secret_revealed"`
	Name             string          `json:"name"`
	RoleAttributes   []RoleAttribute `json:"role_attributes"`
}

type PostgresSQLCluster struct {
	ClusterName           string                  `json:"cluster_name"`
	CreatedAt             gcorecloud.JSONRFC3339Z `json:"created_at"`
	Databases             []DatabaseOverview      `json:"databases"`
	Flavor                Flavor                  `json:"flavor"`
	HighAvailability      *HighAvailability       `json:"high_availability"`
	Network               Network                 `json:"network"`
	PGServerConfiguration PGServerConfiguration   `json:"pg_server_configuration"`
	Status                ClusterStatus           `json:"status"`
	Storage               PGStorageConfiguration  `json:"storage"`
	Users                 []PgUserOverview        `json:"users"`
}

type PostgresSQLClusterShort struct {
	ClusterName string                  `json:"cluster_name"`
	CreatedAt   gcorecloud.JSONRFC3339Z `json:"created_at"`
	Status      ClusterStatus           `json:"status"`
	Version     string                  `json:"version"`
}
