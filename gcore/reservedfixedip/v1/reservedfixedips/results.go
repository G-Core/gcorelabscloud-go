package reservedfixedips

import (
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a reserved fixed ip resource.
func (r commonResult) Extract() (*ReservedFixedIP, error) {
	var s ReservedFixedIP
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type SliceResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a slice of device resource.
func (r SliceResult) Extract() ([]Device, error) {
	var s []Device
	err := r.ExtractInto(&s)
	return s, err
}

func (r SliceResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a ReservedFixedIP.
type GetResult struct {
	commonResult
}

//IPReservation represents the information of an ip reservation
type IPReservation struct {
	Status       string  `json:"status"`
	ResourceType *string `json:"resource_type"`
	ResourceID   *string `json:"resource_id"`
}

//AllowedAddressPairs represents the information of an allowed address pairs
// for any particular ReservedFixedIP
type AllowedAddressPairs struct {
	IPAddress  net.IP `json:"ip_address"`
	MacAddress string `json:"mac_address"`
}

//ReservedFixedIP represents a ReservedFixedIP structure.
type ReservedFixedIP struct {
	PortID              string                  `json:"port_id"`
	Name                string                  `json:"name"`
	CreatedAt           gcorecloud.JSONRFC3339Z `json:"created_at"`
	UpdatedAt           gcorecloud.JSONRFC3339Z `json:"updated_at"`
	Status              string                  `json:"status"`
	FixedIPAddress      net.IP                  `json:"fixed_ip_address"`
	SubnetID            string                  `json:"subnet_id"`
	CreatorTaskID       string                  `json:"creator_task_id"`
	TaskID              *string                 `json:"task_id"`
	IsExternal          bool                    `json:"is_external"`
	IsVip               bool                    `json:"is_vip"`
	Reservation         IPReservation           `json:"reservation"`
	Region              string                  `json:"region"`
	RegionID            int                     `json:"region_id"`
	ProjectID           int                     `json:"project_id"`
	AllowedAddressPairs []AllowedAddressPairs   `json:"allowed_address_pairs"`
	NetworkID           string                  `json:"network_id"`
	Network             networks.Network        `json:"network"`
}

// ReservedFixedIPPage is the page returned by a pager when traversing over a
// collection of ReservedFixedIPs.
type ReservedFixedIPPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of ReservedFixedIP has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ReservedFixedIPPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ReservedFixedIPPage struct is empty.
func (r ReservedFixedIPPage) IsEmpty() (bool, error) {
	is, err := ExtractReservedFixedIPs(r)
	return len(is) == 0, err
}

// ExtractReservedFixedIPs accepts a Page struct, specifically a ReservedFixedIPPage struct,
// and extracts the elements into a slice of ReservedFixedIP structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractReservedFixedIPs(r pagination.Page) ([]ReservedFixedIP, error) {
	var s []ReservedFixedIP
	err := ExtractReservedFixedIPInto(r, &s)
	return s, err
}

func ExtractReservedFixedIPInto(r pagination.Page, v interface{}) error {
	return r.(ReservedFixedIPPage).Result.ExtractIntoSlicePtr(v, "results")
}

// IPAssignment represents a IPAssignment structure.
type IPAssignment struct {
	IPAddress net.IP         `json:"ip_address"`
	SubnetID  string         `json:"subnet_id"`
	Subnet    subnets.Subnet `json:"subnet"`
}

// Device represents a Device structure.
type Device struct {
	PortID        string           `json:"port_id"`
	IPAssignments []IPAssignment   `json:"ip_assignments"`
	InstanceID    string           `json:"instance_id"`
	InstanceName  string           `json:"instance_name"`
	Network       networks.Network `json:"network"`
}

// DevicePage is the page returned by a pager when traversing over a
// collection of ReservedFixedIPs.
type DevicePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of Device has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r DevicePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a DevicePage struct is empty.
func (r DevicePage) IsEmpty() (bool, error) {
	is, err := ExtractDevices(r)
	return len(is) == 0, err
}

// ExtractDevices accepts a Page struct, specifically a DevicePage struct,
// and extracts the elements into a slice of Device structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractDevices(r pagination.Page) ([]Device, error) {
	var s []Device
	err := ExtractDeviceInto(r, &s)
	return s, err
}

func ExtractDeviceInto(r pagination.Page, v interface{}) error {
	return r.(DevicePage).Result.ExtractIntoSlicePtr(v, "results")
}
