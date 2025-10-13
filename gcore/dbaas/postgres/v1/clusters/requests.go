package clusters

import (
	"errors"
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
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

// DeleteOptsBuilder allows extensions to add parameters to delete cluster options.
type DeleteOptsBuilder interface {
	ToClusterDeleteQuery() (string, error)
}

// DeleteOpts specifies the parameters for the Delete method.
type DeleteOpts struct{}

func (d DeleteOpts) ToClusterDeleteQuery() (string, error) {
	return "", nil
}

type DatabaseOpts struct {
	Name  string `json:"name" validate:"required"`
	Owner string `json:"owner" validate:"required"`
}

type FlavorOpts struct {
	CPU       int `json:"cpu" validate:"required,gte=1"`
	MemoryGiB int `json:"memory_gib" validate:"required,gte=1"`
}

type HighAvailabilityOpts struct {
	ReplicationMode HighAvailabilityReplicationMode `json:"replication_mode" validate:"required,oneof=async sync"`
}

type NetworkOpts struct {
	ACL         []string `json:"acl" validate:"required"`
	NetworkType string   `json:"network_type" validate:"required,oneof=public"`
}

type PoolerOpts struct {
	Mode PoolerMode `json:"mode" validate:"required,oneof=session statement transaction"`
	Type string     `json:"type" validate:"oneof=pgbouncer"`
}

type PGServerConfigurationOpts struct {
	PGConf  string      `json:"pg_conf" validate:"required"`
	Version string      `json:"version" validate:"required"`
	Pooler  *PoolerOpts `json:"pooler,omitempty" validate:"omitempty"`
}

type PGServerConfigurationUpdateOpts struct {
	PGConf  string      `json:"pg_conf,omitempty" validate:"omitempty"`
	Version string      `json:"version,omitempty" validate:"omitempty"`
	Pooler  *PoolerOpts `json:"pooler,omitempty" validate:"omitempty"`
}

type PGStorageConfigurationOpts struct {
	SizeGiB int    `json:"size_gib" validate:"required,gte=1,lte=100"`
	Type    string `json:"type" validate:"required"`
}

type PGStorageConfigurationUpdateOpts struct {
	SizeGiB int `json:"size_gib" validate:"required,gte=1,lte=100"`
}

type PgUserOpts struct {
	Name           string          `json:"name" validate:"required"`
	RoleAttributes []RoleAttribute `json:"role_attributes" validate:"required,min=1,dive,oneof=BYPASSRLS CREATEROLE CREATEDB INHERIT LOGIN NOLOGIN"`
}

// CreateOpts specifies the parameters for the Create method.
type CreateOpts struct {
	ClusterName           string                     `json:"cluster_name" validate:"required"`
	Databases             []DatabaseOpts             `json:"databases" validate:"required,dive"`
	Flavor                FlavorOpts                 `json:"flavor" validate:"required"`
	HighAvailability      *HighAvailabilityOpts      `json:"high_availability" validate:"omitempty"`
	Network               NetworkOpts                `json:"network" validate:"required"`
	PGServerConfiguration PGServerConfigurationOpts  `json:"pg_server_configuration" validate:"required"`
	Storage               PGStorageConfigurationOpts `json:"storage" validate:"required"`
	Users                 []PgUserOpts               `json:"users" validate:"required,dive"`
}

// Validate checks if the provided options are valid.
func (opts CreateOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToCreateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Databases             []DatabaseOpts                    `json:"databases,omitempty" validate:"omitempty,dive"`
	Flavor                *FlavorOpts                       `json:"flavor,omitempty" validate:"omitempty"`
	HighAvailability      *HighAvailabilityOpts             `json:"high_availability,omitempty" validate:"omitempty"`
	Network               *NetworkOpts                      `json:"network,omitempty" validate:"omitempty"`
	PGServerConfiguration *PGServerConfigurationUpdateOpts  `json:"pg_server_configuration,omitempty" validate:"omitempty"`
	Storage               *PGStorageConfigurationUpdateOpts `json:"storage,omitempty" validate:"omitempty"`
	Users                 []PgUserOpts                      `json:"users,omitempty" validate:"omitempty,dive"`
}

// Validate checks if the provided options are valid.
func (opts UpdateOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToUpdateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	if len(mp) == 0 {
		return nil, errors.New("empty UpdateOpts")
	}
	return mp, nil
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

func Get(client *gcorecloud.ServiceClient, clusterName string) (r GetResult) {
	url := ClusterURL(client, clusterName)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// List retrieves PostgreSQL clusters.
func List(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := client.ServiceURL()
	if opts != nil {
		query, err := opts.ToListClustersQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.OffsetPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all PG clusters.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]PostgresSQLClusterShort, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractClusters(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Delete deletes a specific PostgreSQL cluster.
func Delete(client *gcorecloud.ServiceClient, clusterName string, opts DeleteOptsBuilder) (r tasks.Result) {
	url := ClusterURL(client, clusterName)
	if opts != nil {
		query, err := opts.ToClusterDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil)
	return
}

// Create creates a new PostgreSQL cluster.
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	url := client.ServiceURL()
	_, r.Err = client.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}

// Update updates a new PostgreSQL cluster.
func Update(client *gcorecloud.ServiceClient, clusterName string, opts UpdateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterURL(client, clusterName)
	_, r.Err = client.Patch(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}
