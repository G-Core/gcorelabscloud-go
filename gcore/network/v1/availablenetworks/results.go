package availablenetworks

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// Network represents a network structure.
type Network struct {
	Name      string                   `json:"name"`
	ID        string                   `json:"id"`
	Subnets   []subnets.Subnet         `json:"subnets"`
	MTU       int                      `json:"mtu"`
	Type      string                   `json:"type"`
	CreatedAt gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	External  bool                     `json:"external"`
	Default   bool                     `json:"default"`
	Shared    bool                     `json:"shared"`
	TaskID    *string                  `json:"task_id"`
	ProjectID int                      `json:"project_id"`
	RegionID  int                      `json:"region_id"`
	Region    string                   `json:"region"`
	Metadata  []metadata.Metadata      `json:"metadata"`
}

// NetworkPage is the page returned by a pager when traversing over a
// collection of networks.
type NetworkPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of networks has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r NetworkPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r NetworkPage) IsEmpty() (bool, error) {
	is, err := ExtractNetworks(r)
	return len(is) == 0, err
}

// ExtractNetwork accepts a Page struct, specifically a NetworkPage struct,
// and extracts the elements into a slice of Network structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNetworks(r pagination.Page) ([]Network, error) {
	var s []Network
	err := ExtractNetworksInto(r, &s)
	return s, err
}

func ExtractNetworksInto(r pagination.Page, v interface{}) error {
	return r.(NetworkPage).Result.ExtractIntoSlicePtr(v, "results")
}
