package networks

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToNetworkListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific network based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNetworkCreateMap() (map[string]interface{}, error)
}

type ListOptsBuilder interface {
	ToNetworkListQuery() (string, error)
}

// CreateOpts represents options used to create a network.
type CreateOpts struct {
	Name         string            `json:"name" required:"true" validate:"required"`
	Mtu          int               `json:"mtu,omitempty" validate:"omitempty,lte=1500"`
	CreateRouter bool              `json:"create_router"`
	Type         string            `json:"type,omitempty" validate:"omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ToNetworkCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToNetworkCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToNetworkListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNetworkListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Create accepts a CreateOpts struct and creates a new network using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToNetworkCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToNetworkUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a network.
type UpdateOpts struct {
	Name string `json:"name" required:"true" validate:"required"`
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToNetworkUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToNetworkUpdateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts UpdateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Update accepts a UpdateOpts struct and updates an existing network using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, networkID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNetworkUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, networkID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the network associated with it.
func Delete(c *gcorecloud.ServiceClient, networkID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, networkID), &r.Body, nil)
	return
}

// ListAll is a convenience function that returns all networks.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Network, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractNetworks(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListAllInstancePort retrieves a specific list of instance ports by network_id.
func ListAllInstancePort(c *gcorecloud.ServiceClient, networkID string) (r GetInstancePortResult) {
	url := listInstancePortsURL(c, networkID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// IDFromName is a convenience function that returns a network ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractNetworks(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "networks"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "networks"}
	}
}
