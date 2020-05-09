package resources

import (
	"bytes"

	"github.com/G-Core/gcorelabscloud-go/gcore/heat/v1/stack/resources/types"

	"github.com/G-Core/gcorelabscloud-go/pagination"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Metadata retrieves metadata for heat resource
func Metadata(c *gcorecloud.ServiceClient, id, resource string) (r MetadataResult) {
	url := MetadataURL(c, id, resource)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Signal set heat resource status
func Signal(c *gcorecloud.ServiceClient, id, resource string, body []byte) (r SignalResult) {
	url := SignalURL(c, id, resource)
	_, r.Err = c.Post(url, nil, nil, &gcorecloud.RequestOpts{
		RawBody: bytes.NewReader(body),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	})
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToResourceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Type               string                    `q:"type"`
	Name               string                    `q:"name"`
	Status             types.StackResourceStatus `q:"status"`
	Action             types.StackResourceAction `q:"name"`
	LogicalResourceID  string                    `q:"id"`
	PhysicalResourceID string                    `q:"physical_resource_id"`
	NestedDepth        int                       `q:"nested_depth"`
	WithDetail         bool                      `q:"with_detail"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToResourceListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// MarkUnhealthyOpts contains the common options struct used in this package's
// MarkUnhealthy operations.
type MarkUnhealthyOpts struct {
	// A boolean indicating whether the target resource should be marked as unhealthy.
	MarkUnhealthy bool `json:"mark_unhealthy"`
	// The reason for the current stack resource state.
	ResourceStatusReason string `json:"resource_status_reason,omitempty"`
}

// MarkUnhealthyOptsBuilder is the interface options structs have to satisfy in order
// to be used in the MarkUnhealthy operation in this package
type MarkUnhealthyOptsBuilder interface {
	ToMarkUnhealthyMap() (map[string]interface{}, error)
}

// ToMarkUnhealthyMap validates that a template was supplied and calls
// the ToMarkUnhealthyMap private function.
func (opts MarkUnhealthyOpts) ToMarkUnhealthyMap() (map[string]interface{}, error) {
	b, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MarkUnhealthy marks the specified resource in the stack as unhealthy.
func MarkUnhealthy(c *gcorecloud.ServiceClient, stackID, resourceName string, opts MarkUnhealthyOptsBuilder) (r MarkUnhealthyResult) {
	b, err := opts.ToMarkUnhealthyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(markUnhealthyURL(c, stackID, resourceName), b, nil, nil)
	return
}

// List resources.
func List(c *gcorecloud.ServiceClient, stackID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c, stackID)
	if opts != nil {
		query, err := opts.ToResourceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ResourcePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific resource based on its unique ID.
func Get(c *gcorecloud.ServiceClient, stackID, resourceName string) (r GetResult) {
	url := getURL(c, stackID, resourceName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListAll is a convenience function that returns a all stack resources.
func ListAll(client *gcorecloud.ServiceClient, stackID string, opts ListOptsBuilder) ([]ResourceList, error) {
	pages, err := List(client, stackID, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractResources(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
