package pools

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters
// to the Create request.
type CreateOptsBuilder interface {
	ToClusterPoolCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster Pool.
type CreateOpts struct {
	Name               string                         `json:"name" required:"true" validate:"required"`
	FlavorID           string                         `json:"flavor_id" required:"true" validate:"required"`
	MinNodeCount       int                            `json:"min_node_count" required:"true" validate:"required,gt=0,ltefield=MaxNodeCount"`
	MaxNodeCount       int                            `json:"max_node_count,omitempty" validate:"omitempty,gt=0,gtefield=MinNodeCount"`
	BootVolumeSize     int                            `json:"boot_volume_size,omitempty" validate:"omitempty,gt=0"`
	BootVolumeType     volumes.VolumeType             `json:"boot_volume_type,omitempty" validate:"omitempty,enum"`
	AutoHealingEnabled bool                           `json:"auto_healing_enabled,omitempty"`
	ServerGroupPolicy  servergroups.ServerGroupPolicy `json:"servergroup_policy,omitempty" validate:"omitempty,enum"`
	IsPublicIPv4       bool                           `json:"is_public_ipv4,omitempty"`
	Labels             map[string]string              `json:"labels,omitempty"`
	Taints             map[string]string              `json:"taints,omitempty"`
	CrioConfig         map[string]string              `json:"crio_config,omitempty" validate:"omitempty"`
	KubeletConfig      map[string]string              `json:"kubelet_config,omitempty" validate:"omitempty"`
}

// Validate CreateOpts
func (opts CreateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// ToClusterPoolCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterPoolCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateOptsBuilder allows extensions to add additional parameters
// to the Update request.
type UpdateOptsBuilder interface {
	ToClusterPoolUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a pool.
type UpdateOpts struct {
	AutoHealingEnabled *bool              `json:"auto_healing_enabled,omitempty"`
	MinNodeCount       int                `json:"min_node_count,omitempty" validate:"omitempty,gt=0,lte=20,ltefield=MaxNodeCount"`
	MaxNodeCount       int                `json:"max_node_count,omitempty" validate:"omitempty,gt=0,lte=20,gtefield=MinNodeCount"`
	Labels             *map[string]string `json:"labels,omitempty"`
	Taints             *map[string]string `json:"taints,omitempty"`
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

// ResizeOptsBuilder allows extensions to add additional parameters to the Resize request.
type ResizeOptsBuilder interface {
	ToClusterPoolResizeMap() (map[string]interface{}, error)
}

// ResizeOpts represents options used to update a cluster.
type ResizeOpts struct {
	NodeCount int `json:"node_count" required:"true" validate:"required,gt=0,lte=20"`
}

// Validate ResizeOpts
func (opts ResizeOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// ToClusterPoolResizeMap builds a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterPoolResizeMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// List returns a Pager which allows you to iterate over a collection of cluster pools.
func List(c *gcorecloud.ServiceClient, clusterName string) pagination.Pager {
	url := listURL(c, clusterName)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all cluster pools.
func ListAll(client *gcorecloud.ServiceClient, clusterName string) ([]ClusterPool, error) {
	pages, err := List(client, clusterName).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractClusterPools(pages)
}

// Create accepts a CreateOpts struct and creates a new cluster pool
// using the values provided.
func Create(c *gcorecloud.ServiceClient, clusterName string, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c, clusterName), b, &r.Body, nil)
	return
}

// Get retrieves a specific cluster pool based on its name.
func Get(c *gcorecloud.ServiceClient, clusterName, poolName string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, clusterName, poolName), &r.Body, nil)
	return
}

// Update accepts an UpdateOpts struct and updates an existing cluster pool
// using the values provided.
func Update(c *gcorecloud.ServiceClient, clusterName, poolName string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, clusterName, poolName), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a pool name and deletes the cluster pool associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterName, poolName string) (r tasks.Result) {
	url := deleteURL(c, clusterName, poolName)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, nil)
	return
}

// Resize accepts a ResizeOpts struct and resizes an existing cluster
// using the values provided.
func Resize(c *gcorecloud.ServiceClient, clusterName, poolName string, opts ResizeOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterPoolResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(resizeURL(c, clusterName, poolName), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// ListInstances returns a Pager which allows you to iterate over a collection of pool instances.
func ListInstances(c *gcorecloud.ServiceClient, clusterName, poolName string) pagination.Pager {
	url := instancesURL(c, clusterName, poolName)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// List returns all pool instances.
func ListInstancesAll(c *gcorecloud.ServiceClient, clusterName, poolName string) ([]instances.Instance, error) {
	page, err := ListInstances(c, clusterName, poolName).AllPages()
	if err != nil {
		return nil, err
	}
	return instances.ExtractInstances(page)
}
