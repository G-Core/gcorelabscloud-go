package regions

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/types"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

type ListOpts struct {
	ShowVolumeTypes bool `json:"show_volume_types,omitempty" validate:"omitempty"`
}

// ToInstanceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RegionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetOptsBuilder allows extensions to add additional parameters to the Get request.
type GetOptsBuilder interface {
	ToInstanceGetQuery() (string, error)
}

type GetOpts struct {
	ShowVolumeTypes bool `json:"show_volume_types,omitempty" validate:"omitempty"`
}

// ToInstanceGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToInstanceGetQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Get retrieves a specific region based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id int, opts GetOptsBuilder) (r GetResult) {
	url := getURL(c, id)
	if opts != nil {
		query, err := opts.ToInstanceGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToRegionCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a region.
type CreateOpts struct {
	DisplayName       string             `json:"display_name" required:"true" validate:"required"`
	KeystoneName      string             `json:"keystone_name" required:"true" validate:"required"`
	State             types.RegionState  `json:"state" required:"true" validate:"required,enum"`
	EndpointType      types.EndpointType `json:"endpoint_type,omitempty" validate:"omitempty,enum"`
	ExternalNetworkID string             `json:"external_network_id" required:"true" validate:"required,uuid4"`
	SpiceProxyURL     *gcorecloud.URL    `json:"spice_proxy_url,omitempty"`
	KeystoneID        int                `json:"keystone_id" required:"true" validate:"required"`
}

// ToRegionCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRegionCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// Create accepts a CreateOpts struct and creates a new region using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRegionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToRegionUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a region.
type UpdateOpts struct {
	DisplayName       string             `json:"display_name,omitempty" validate:"required_without_all=State EndpointType ExternalNetworkID SpiceProxyURL,omitempty"`
	State             types.RegionState  `json:"state,omitempty" validate:"required_without_all=DisplayName EndpointType ExternalNetworkID SpiceProxyURL,omitempty,enum"`
	EndpointType      types.EndpointType `json:"endpoint_type,omitempty" validate:"required_without_all=DisplayName State ExternalNetworkID SpiceProxyURL,omitempty,enum"`
	ExternalNetworkID string             `json:"external_network_id,omitempty" validate:"required_without_all=DisplayName State EndpointType SpiceProxyURL,omitempty,uuid4"`
	SpiceProxyURL     *gcorecloud.URL    `json:"spice_proxy_url,omitempty" validate:"required_without_all=DisplayName State EndpointType ExternalNetworkID,omitempty"`
}

// ToRegionUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToRegionUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts UpdateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// Update accepts a UpdateOpts struct and updates an existing region using the values provided.
func Update(c *gcorecloud.ServiceClient, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRegionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, id), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// ListAll is a convenience function that returns all regions.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Region, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractRegions(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
