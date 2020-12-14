package routers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// List returns a Pager which allows you to iterate over a collection of
// routers.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToRouterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RouterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ToRouterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRouterListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Interface represents a list of interfaces to attach to router immediately after creation.
type Interface struct {
	Type     types.InterfaceType `json:"type,omitempty" validate:"enum,required_with=SubnetID,omitempty"`
	SubnetID string              `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,required_with=Type,omitempty,uuid4"`
}

// GatewayInfo represents the information of an external gateway for any
// particular network router.
type GatewayInfo struct {
	Type       types.GatewayType `json:"type,omitempty" validate:"omitempty,enum"`
	EnableSNat *bool             `json:"enable_snat"`
	NetworkID  string            `json:"network_id,omitempty" validate:"rfe=Type:manual,omitempty,uuid4"`
}

// CreateOpts represents options used to create a router.
type CreateOpts struct {
	Name                string              `json:"name" required:"true"`
	ExternalGatewayInfo GatewayInfo         `json:"external_gateway_info,omitempty"`
	Interfaces          []Interface         `json:"interfaces,omitempty"`
	Routes              []subnets.HostRoute `json:"routes,omitempty"`
}

// UpdateOpts represents options used to update a router.
type UpdateOpts struct {
	Name                string              `json:"name,omitempty"`
	ExternalGatewayInfo GatewayInfo         `json:"external_gateway_info,omitempty"`
	Routes              []subnets.HostRoute `json:"routes"`
}

// ListOpts allows the filtering and sorting List API response.
type ListOpts struct {
	ID        string `q:"id"`
	Name      string `q:"name"`
	Status    string `q:"status"`
	ProjectID string `q:"project_id"`
	Limit     int    `q:"limit"`
}

type AttachOpts struct {
	SubnetID string `json:"subnet_id" required:"true"`
}

func (opts *Interface) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

func (opts *GatewayInfo) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToRouterCreateMap() (map[string]interface{}, error)
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToRouterUpdateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToRouterListQuery() (string, error)
}

// ToRouterCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRouterCreateMap() (map[string]interface{}, error) {
	for _, i := range opts.Interfaces {
		err := i.Validate()
		if err != nil {
			return nil, err
		}
	}
	err := gcorecloud.TranslateValidationError(opts.ExternalGatewayInfo.Validate())
	if err != nil {
		return nil, err
	}
	body, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	gw := body["external_gateway_info"].(map[string]interface{})

	if len(gw) == 1 && gw["enable_snat"] == nil {
		delete(body, "external_gateway_info")
	}

	return body, err
}

// ToRouterUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToRouterUpdateMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.ExternalGatewayInfo.Validate())
	if err != nil {
		return nil, err
	}
	body, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	gw := body["external_gateway_info"].(map[string]interface{})

	if len(gw) == 1 && gw["enable_snat"] == nil {
		delete(body, "external_gateway_info")
	}

	if body["routes"] == nil {
		delete(body, "routes")
	}

	return body, err
}

// Get retrieves a specific router based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Create accepts a CreateOpts struct and creates a new router using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToRouterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Update accepts a UpdateOpts struct and updates an existing router using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, routerID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRouterUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, routerID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the router associated with it.
func Delete(c *gcorecloud.ServiceClient, routerID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, routerID), &r.Body, nil)
	return
}

// ListAll returns all routers.
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Router, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractRouters(page)
}

// Attach subnet to router.
func Attach(c *gcorecloud.ServiceClient, routerID string, subnetID string) (r GetResult) {
	attachOpts := AttachOpts{SubnetID: subnetID}
	body, err := gcorecloud.BuildRequestBody(attachOpts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(attachURL(c, routerID), body, &r.Body, nil)
	return
}

// Detach subnet to router.
func Detach(c *gcorecloud.ServiceClient, routerID string, subnetID string) (r GetResult) {
	attachOpts := AttachOpts{SubnetID: subnetID}
	body, err := gcorecloud.BuildRequestBody(attachOpts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(detachURL(c, routerID), body, &r.Body, nil)
	return
}
