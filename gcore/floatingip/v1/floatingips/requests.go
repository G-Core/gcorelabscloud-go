package floatingips

import (
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToFloatingIPListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FloatingIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type ListOptsBuilder interface {
	ToFloatingIPListQuery() (string, error)
}

// ToFloatingIPListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFloatingIPListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Get retrieves a specific floating ip based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder adds additional parameters to the request.
type CreateOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a floating ip.
type CreateOpts struct {
	PortID         string            `json:"port_id,omitempty"`
	FixedIPAddress net.IP            `json:"fixed_ip_address,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ToFloatingIPCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new floating ip using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToFloatingIPUpdateMap() (map[string]interface{}, error)
}

// Delete accepts a unique ID and deletes the floating ip associated with it.
func Delete(c *gcorecloud.ServiceClient, floatingID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, floatingID), &r.Body, nil)
	return
}

// ListAll is a convenience function that returns all floating IPs.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]FloatingIPDetail, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractFloatingIPs(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// Assign accepts a CreateOpts struct and assign floating IP.
func Assign(c *gcorecloud.ServiceClient, floatingIPID string, opts CreateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(assignURL(c, floatingIPID), b, &r.Body, nil)
	return
}

func UnAssign(c *gcorecloud.ServiceClient, floatingIPID string) (r UpdateResult) {
	_, r.Err = c.Post(unAssignURL(c, floatingIPID), nil, &r.Body, nil)
	return
}
