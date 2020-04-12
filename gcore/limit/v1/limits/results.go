package limits

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

type LimitResponse struct {
	ID        int                       `json:"id"`
	ClientID  int                       `json:"client_id"`
	Limits    string                    `json:"limits"`
	Status    types.LimitRequestStatus  `json:"status"`
	CreatedAt gcorecloud.JSONRFC3339NoZ `json:"created_at"`
}

// Extract is a function that accepts a result and extracts a limit response resource.
func (r commonResult) Extract() (*LimitResponse, error) {
	var s LimitResponse
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a LimitResponse.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LimitResponse.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a LimitResponse.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of an delete operation. Call its ExtractErr to get operation error.
type DeleteResult struct {
	gcorecloud.ErrResult
}

// LimitResultPage is the page returned by a pager when traversing over a
// collection of limit requests.
type LimitResultPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of limit requests has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r LimitResultPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a LimitResultPage struct is empty.
func (r LimitResultPage) IsEmpty() (bool, error) {
	is, err := ExtractLimitResults(r)
	return len(is) == 0, err
}

// ExtractLimitResult accepts a Page struct, specifically a LimitResultPage struct,
// and extracts the elements into a slice of LimitResponse structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractLimitResults(r pagination.Page) ([]LimitResponse, error) {
	var s []LimitResponse
	err := ExtractLimitResultsInto(r, &s)
	return s, err
}

func ExtractLimitResultsInto(r pagination.Page, v interface{}) error {
	return r.(LimitResultPage).Result.ExtractIntoSlicePtr(v, "results")
}
