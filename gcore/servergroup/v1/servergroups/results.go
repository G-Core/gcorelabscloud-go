package servergroups

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a server group resource.
func (r commonResult) Extract() (*ServerGroup, error) {
	var s ServerGroup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a ServerGroup.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// ServerGroupInstance represent an instances in server group
type ServerGroupInstance struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`
}

// ServerGroup represents a server group.
type ServerGroup struct {
	ServerGroupID string                `json:"servergroup_id"`
	ProjectID     int                   `json:"project_id"`
	RegionID      int                   `json:"region_id"`
	Region        string                `json:"region"`
	Name          string                `json:"name"`
	Instances     []ServerGroupInstance `json:"instances"`
	Policy        ServerGroupPolicy     `json:"policy"`
}

// ServerGroupPage is the page returned by a pager when traversing over a
// collection of server groups.
type ServerGroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of server groups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ServerGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ServerGroup struct is empty.
func (r ServerGroupPage) IsEmpty() (bool, error) {
	is, err := ExtractServerGroups(r)
	return len(is) == 0, err
}

// ExtractServerGroups accepts a Page struct, specifically a ServerGroupPage struct,
// and extracts the elements into a slice of ServerGroup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractServerGroups(r pagination.Page) ([]ServerGroup, error) {
	var s []ServerGroup
	err := ExtractServerGroupsInto(r, &s)
	return s, err
}

func ExtractServerGroupsInto(r pagination.Page, v interface{}) error {
	return r.(ServerGroupPage).Result.ExtractIntoSlicePtr(v, "results")
}
