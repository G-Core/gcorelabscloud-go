package reservedfixedips

import (
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToReservedFixedIPListQuery() (string, error)
}

// ListOpts allows the filtering and sorting List API response.
type ListOpts struct {
	ExternalOnly  bool   `q:"external_only"`
	InternalOnly  bool   `q:"internal_only"`
	AvailableOnly bool   `q:"available_only"`
	VipOnly       bool   `q:"vip_only"`
	DeviceID      string `q:"device_id"`
}

// ToReservedFixedIPListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReservedFixedIPListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ToReservedFixedIPCreateMap allows extensions to add additional parameters to the Create request
type CreateOptsBuilder interface {
	ToReservedFixedIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a reserved fixed ip.
type CreateOpts struct {
	Type      ReservedFixedIPType `json:"type" required:"true" validate:"required,enum"`
	IPFamily  IPFamilyType        `json:"ip_family,omitempty" validate:"omitempty,enum"`
	NetworkID string              `json:"network_id,omitempty" validate:"rfe=Type:ip_address;any_subnet,omitempty,uuid4"`
	SubnetID  string              `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	IPAddress net.IP              `json:"ip_address,omitempty" validate:"rfe=Type:ip_address,omitempty"`
	IsVip     bool                `json:"is_vip"`
}

func (opts *CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToReservedFixedIPCreateMap builds a request body.
func (opts CreateOpts) ToReservedFixedIPCreateMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToReservedFixedIPSwitchVIPMap allows extensions to add additional parameters to the SwitchVIP request
type SwitchVIPOptsBuilder interface {
	ToReservedFixedIPSwitchVIPMap() (map[string]interface{}, error)
}

// SwitchVIPOpts represents options used to switch vip status.
type SwitchVIPOpts struct {
	IsVip bool `json:"is_vip"`
}

// ToReservedFixedIPSwitchVIPMap builds a request body.
func (opts SwitchVIPOpts) ToReservedFixedIPSwitchVIPMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToPortsToShareVIPOptsMap allows extensions to add additional parameters to the AddPortsToShareVIP request
type PortsToShareVIPOptsBuilder interface {
	ToPortsToShareVIPOptsMap() (map[string]interface{}, error)
}

// PortsToShareVIPOpts represents options used to add ports to share vip.
type PortsToShareVIPOpts struct {
	PortIDs []string `json:"port_ids" validate:"required,dive,uuid4"`
}

func (opts *PortsToShareVIPOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToAddPortsToShareVIPOptsMap builds a request body.
func (opts PortsToShareVIPOpts) ToPortsToShareVIPOptsMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// List retrieves list of reserved fixed ips.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToReservedFixedIPListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ReservedFixedIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll returns all reserved fixed ips.
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]ReservedFixedIP, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractReservedFixedIPs(page)
}

// Get retrieves a specific reserved fixed ip based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// Create accepts a CreateOpts struct and creates a new reserved fixed ip using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToReservedFixedIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the reserved fixed ip associated with it.
func Delete(c *gcorecloud.ServiceClient, id string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, id), &r.Body, nil)
	return
}

// ListConnectedDevice accepts a unique ID and retrieves list of connected devices associated with it.
func ListConnectedDevice(c *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := connectedDeviceListURL(c, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return DevicePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllConnectedDevice returns all connected device.
func ListAllConnectedDevice(c *gcorecloud.ServiceClient, id string) ([]Device, error) {
	page, err := ListConnectedDevice(c, id).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractDevices(page)
}

// ListAvailableDevice accepts a unique ID and retrieves list of available devices associated with it.
func ListAvailableDevice(c *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := availableDeviceListURL(c, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return DevicePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllAvailableDevice returns all available device.
func ListAllAvailableDevice(c *gcorecloud.ServiceClient, id string) ([]Device, error) {
	page, err := ListAvailableDevice(c, id).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractDevices(page)
}

// SwitchVIP accepts a SwitchVIPOpts struct and switch vip status using the values provided.
func SwitchVIP(c *gcorecloud.ServiceClient, id string, opts SwitchVIPOptsBuilder) (r GetResult) {
	b, err := opts.ToReservedFixedIPSwitchVIPMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(switchVIPURL(c, id), b, &r.Body, nil)
	return
}

// AddPortsToShareVIP accepts a unique ID, PortsToShareVIPOpts and add ports to share vip using the value provided.
func AddPortsToShareVIP(c *gcorecloud.ServiceClient, id string, opts PortsToShareVIPOptsBuilder) (r SliceResult) {
	b, err := opts.ToPortsToShareVIPOptsMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Patch(portsToShareVIPURL(c, id), b, &r.Body, nil)
	return
}

// ReplacePortsToShareVIP accepts a unique ID, PortsToShareVIPOpts and replace ports to share vip using the value provided.
func ReplacePortsToShareVIP(c *gcorecloud.ServiceClient, id string, opts PortsToShareVIPOptsBuilder) (r SliceResult) {
	b, err := opts.ToPortsToShareVIPOptsMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(portsToShareVIPURL(c, id), b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{200}})
	return
}
