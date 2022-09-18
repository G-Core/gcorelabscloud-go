package subnets

import (
	"fmt"
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a subnet resource.
func (r commonResult) Extract() (*Subnet, error) {
	var s Subnet
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Subnet.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Subnet.
type UpdateResult struct {
	commonResult
}

// Subnet represents a subnet structure.
type Subnet struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	IPVersion      int                     `json:"ip_version"`
	EnableDHCP     bool                    `json:"enable_dhcp"`
	CIDR           gcorecloud.CIDR         `json:"cidr"`
	CreatedAt      gcorecloud.JSONRFC3339Z `json:"created_at"`
	UpdatedAt      gcorecloud.JSONRFC3339Z `json:"updated_at"`
	NetworkID      string                  `json:"network_id"`
	TaskID         string                  `json:"task_id"`
	CreatorTaskID  string                  `json:"creator_task_id"`
	Region         string                  `json:"region"`
	ProjectID      int                     `json:"project_id"`
	RegionID       int                     `json:"region_id"`
	AvailableIps   int                     `json:"available_ips"`
	TotalIps       int                     `json:"total_ips"`
	HasRouter      bool                    `json:"has_router"`
	DNSNameservers []net.IP                `json:"dns_nameservers"`
	HostRoutes     []HostRoute             `json:"host_routes"`
	GatewayIP      net.IP                  `json:"gateway_ip"`
	Metadata       []metadata.Metadata     `json:"metadata"`
}

// SubnetPage is the page returned by a pager when traversing over a
// collection of subnets.
type SubnetPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of subnets has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SubnetPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SubnetPage struct is empty.
func (r SubnetPage) IsEmpty() (bool, error) {
	is, err := ExtractSubnets(r)
	return len(is) == 0, err
}

// ExtractSubnet accepts a Page struct, specifically a SubnetPage struct,
// and extracts the elements into a slice of Subnet structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSubnets(r pagination.Page) ([]Subnet, error) {
	var s []Subnet
	err := ExtractSubnetsInto(r, &s)
	return s, err
}

func ExtractSubnetsInto(r pagination.Page, v interface{}) error {
	return r.(SubnetPage).Result.ExtractIntoSlicePtr(v, "results")
}

type SubnetTaskResult struct {
	Subnets []string `json:"subnets"`
}

func ExtractSubnetIDFromTask(task *tasks.Task) (string, error) {
	var result SubnetTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode subnet information in task structure: %w", err)
	}
	if len(result.Subnets) == 0 {
		return "", fmt.Errorf("cannot decode subnet information in task structure: %w", err)
	}
	return result.Subnets[0], nil
}
