package regionsaccess

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a region resource.
func (r commonResult) Extract() (*RegionAccess, error) {
	var s RegionAccess
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of an update operation. Call its Extract
// method to interpret it as a RegionAccess.
type CreateResult struct {
	commonResult
}

// DeleteResult represent a result of a deletion operation.
type DeleteResult struct {
	gcorecloud.ErrResult
}

// RegionAccessPage is the page returned by a pager when traversing over a
// collection of regions access.
type RegionAccessPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of regions has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r RegionAccessPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RegionPage struct is empty.
func (r RegionAccessPage) IsEmpty() (bool, error) {
	is, err := ExtractRegionsAccess(r)
	return len(is) == 0, err
}

// ExtractRegionsAccess accepts a Page struct, specifically a RegionAccessPage struct,
// and extracts the elements into a slice of RegionAccess structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRegionsAccess(r pagination.Page) ([]RegionAccess, error) {
	var s []RegionAccess
	err := ExtractRegionsAccessInto(r, &s)
	return s, err
}

func ExtractRegionsAccessInto(r pagination.Page, v interface{}) error {
	return r.(RegionAccessPage).Result.ExtractIntoSlicePtr(v, "results")
}

type RegionAccess struct {
	ID                   int   `json:"id"`
	AccessAllEdgeRegions bool  `json:"access_all_edge_regions"`
	ClientID             *int  `json:"client_id"`
	RegionIDs            []int `json:"region_ids"`
	ResellerID           *int  `json:"reseller_id"`
}
