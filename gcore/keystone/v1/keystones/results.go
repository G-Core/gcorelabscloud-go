package keystones

import (
	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/keystone/v1/types"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Keystone represents a keystone structure.
type Keystone struct {
	ID                        int                       `json:"id"`
	URL                       gcorecloud.URL            `json:"url"`
	State                     types.KeystoneState       `json:"state"`
	KeystoneFederatedDomainID string                    `json:"keystone_federated_domain_id"`
	CreatedOn                 gcorecloud.JSONRFC3339NoZ `json:"created_on"`
	AdminPassword             string                    `json:"admin_password"`
}

// Extract is a function that accepts a result and extracts a keystone resource.
func (r commonResult) Extract() (*Keystone, error) {
	var s Keystone
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Keystone.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Keystone.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Keystone.
type UpdateResult struct {
	commonResult
}

// KeystonePage is the page returned by a pager when traversing over a
// collection of keystones.
type KeystonePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of keystones has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r KeystonePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a KeystonePage struct is empty.
func (r KeystonePage) IsEmpty() (bool, error) {
	is, err := ExtractKeystones(r)
	return len(is) == 0, err
}

// ExtractKeystone accepts a Page struct, specifically a KeystonePage struct,
// and extracts the elements into a slice of Keystone structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractKeystones(r pagination.Page) ([]Keystone, error) {
	var s []Keystone
	err := ExtractKeystonesInto(r, &s)
	return s, err
}

func ExtractKeystonesInto(r pagination.Page, v interface{}) error {
	return r.(KeystonePage).Result.ExtractIntoSlicePtr(v, "results")
}
