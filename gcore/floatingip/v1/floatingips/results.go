package floatingips

import (
	"fmt"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a security group resource.
func (r commonResult) Extract() (*instances.FloatingIP, error) {
	var s instances.FloatingIP
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a FloatingIP.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a assign or unassign operation. Call its Extract
// method to interpret it as a FloatingIP.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a FloatingIP.
type GetResult struct {
	commonResult
}

// FloatingIPDetail represents a floating IP with details.
type FloatingIPDetail struct {
	*instances.FloatingIP
	Instance instances.Instance `json:"instance"`
}

// FloatingIPPage is the page returned by a pager when traversing over a
// collection of security groups.
type FloatingIPPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of security groups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FloatingIPPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FloatingIPPage struct is empty.
func (r FloatingIPPage) IsEmpty() (bool, error) {
	is, err := ExtractFloatingIPs(r)
	return len(is) == 0, err
}

// ExtractFloatingIP accepts a Page struct, specifically a FloatingIPPage struct,
// and extracts the elements into a slice of FloatingIP structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFloatingIPs(r pagination.Page) ([]FloatingIPDetail, error) {
	var s []FloatingIPDetail
	err := ExtractFloatingIPsInto(r, &s)
	return s, err
}

func ExtractFloatingIPsInto(r pagination.Page, v interface{}) error {
	return r.(FloatingIPPage).Result.ExtractIntoSlicePtr(v, "results")
}

type FloatingIPTaskResult struct {
	FloatingIPs []string `json:"floatingips"`
}

func ExtractFloatingIPIDFromTask(task *tasks.Task) (string, error) {
	var result FloatingIPTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode floating IP in task structure: %w", err)
	}
	if len(result.FloatingIPs) == 0 {
		return "", fmt.Errorf("cannot decode floating IP in task structure: %w", err)
	}
	return result.FloatingIPs[0], nil
}
