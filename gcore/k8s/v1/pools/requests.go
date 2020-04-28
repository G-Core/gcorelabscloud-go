package pools

import (
	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterPoolsListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the cluster pools attributes you want to see returned. SortKey allows you to sort
// by a particular cluster pools attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Marker  string `q:"marker"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
	Detail  bool   `q:"detail"`
}

// ToClusterPoolsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterPoolsListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// cluster pools. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c, clusterID)
	if opts != nil {
		query, err := opts.ToClusterPoolsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Instances returns a Pager which allows you to iterate over a collection of pool instances.
func Instances(c *gcorecloud.ServiceClient, clusterID string, id string) pagination.Pager {
	url := instancesURL(c, clusterID, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// Volumes returns a Pager which allows you to iterate over a collection of pool instances.
func Volumes(c *gcorecloud.ServiceClient, clusterID string, id string) pagination.Pager {
	url := volumesURL(c, clusterID, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return volumes.VolumePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster pool based on its unique ID.
func Get(c *gcorecloud.ServiceClient, clusterID string, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, clusterID, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToClusterPoolCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster Pool.
type CreateOpts struct {
	Name             string `json:"name" required:"true" validate:"required"`
	FlavorID         string `json:"flavor_id" required:"true" validate:"required"`
	NodeCount        int    `json:"node_count" required:"true" validate:"required,gt=0"`
	DockerVolumeSize int    `json:"docker_volume_size,omitempty" validate:"omitempty,required,gt=0"`
	MinNodeCount     int    `json:"min_node_count,omitempty" validate:"omitempty,required,gt=0,ltefield=NodeCount"`
	MaxNodeCount     int    `json:"max_node_count,omitempty" validate:"omitempty,required,gt=0,gtefield=MinNodeCount,gtefield=NodeCount"`
}

// ToClusterPoolCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterPoolCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate CreateOpts
func (opts CreateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Create accepts a CreateOpts struct and creates a new cluster Pool using the values
// provided. This operation does not actually require a request body, i.e. the CreateOpts struct argument can be empty.
func Create(c *gcorecloud.ServiceClient, clusterID string, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c, clusterID), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterPoolUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a pool.
type UpdateOpts struct {
	Name         string `json:"name,omitempty" validate:"required_without_all=MinNodeCount MaxNodeCount,omitempty"`
	MinNodeCount int    `json:"min_node_count,omitempty" validate:"required_without_all=Name MaxNodeCount,omitempty,gt=0"`
	MaxNodeCount int    `json:"max_node_count,omitempty" validate:"required_without_all=Name MixNodeCount,omitempty,gt=0,gtefield=MinNodeCount"`
}

// Validate UpdateOpts
func (opts UpdateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// ToClusterPoolUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToClusterPoolUpdateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing Pool using the values provided.
func Update(c *gcorecloud.ServiceClient, clusterID, poolID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, clusterID, poolID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the cluster Pool associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterID, poolID string) (r tasks.Result) {
	url := deleteURL(c, clusterID, poolID)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, nil)
	return
}

// ListAll is a convenience function that returns a all cluster pools.
func ListAll(client *gcorecloud.ServiceClient, clusterID string, opts ListOptsBuilder) ([]ClusterListPool, error) {
	pages, err := List(client, clusterID, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusterPools(pages)
}

// List returns all pool instances.
func InstancesAll(c *gcorecloud.ServiceClient, clusterID, id string) ([]instances.Instance, error) {
	page, err := Instances(c, clusterID, id).AllPages()
	if err != nil {
		return nil, err
	}
	return instances.ExtractInstances(page)
}

// List returns all pool volumes.
func VolumesAll(c *gcorecloud.ServiceClient, clusterID, id string) ([]volumes.Volume, error) {
	page, err := Volumes(c, clusterID, id).AllPages()
	if err != nil {
		return nil, err
	}
	return volumes.ExtractVolumes(page)
}
