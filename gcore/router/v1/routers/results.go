package routers

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a router resource.
func (r commonResult) Extract() (*Router, error) {
	var s Router
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a router.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a router.
type UpdateResult struct {
	commonResult
}

// ExternalFixedIP is the IP address and subnet ID of the external gateway of a
// router.
type ExtFixedIPs struct {
	IPAddress string `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

// GatewayInfo represents the information of an external gateway for any
// particular network router.
type ExtGatewayInfo struct {
	EnableSNat       bool          `json:"enable_snat"`
	ExternalFixedIPs []ExtFixedIPs `json:"external_fixed_ips"`
	NetworkID        string        `json:"network_id"`
}

// Router represents a router structure.
type Router struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Status              string                  `json:"status"`
	ExternalGatewayInfo ExtGatewayInfo          `json:"external_gateway_info"`
	Routes              []subnets.HostRoute     `json:"routes"`
	Interfaces          []instances.Interface   `json:"interfaces"`
	TaskID              string                  `json:"task_id"`
	CreatorTaskID       string                  `json:"creator_task_id"`
	ProjectID           int                     `json:"project_id"`
	RegionID            int                     `json:"region_id"`
	CreatedAt           gcorecloud.JSONRFC3339Z `json:"created_at"`
	UpdatedAt           gcorecloud.JSONRFC3339Z `json:"updated_at"`
}

// RouterPage is the page returned by a pager when traversing over a
// collection of routers.
type RouterPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r RouterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RouterPage struct is empty.
func (r RouterPage) IsEmpty() (bool, error) {
	is, err := ExtractRouters(r)
	return len(is) == 0, err
}

// ExtractRouter accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRouters(r pagination.Page) ([]Router, error) {
	var s []Router
	err := ExtractRoutersInto(r, &s)
	return s, err
}

func ExtractRoutersInto(r pagination.Page, v interface{}) error {
	return r.(RouterPage).Result.ExtractIntoSlicePtr(v, "results")
}

type RouterTaskResult struct {
	Routers []string `json:"routers"`
}

func ExtractRouterIDFromTask(task *tasks.Task) (string, error) {
	var result RouterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode router information in task structure: %w", err)
	}
	if len(result.Routers) == 0 {
		return "", fmt.Errorf("cannot decode router information in task structure: %w", err)
	}
	return result.Routers[0], nil
}
