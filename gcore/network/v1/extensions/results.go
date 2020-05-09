package extensions

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a network resource.
func (r commonResult) Extract() (*Extension, error) {
	var s Extension
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Extension.
type GetResult struct {
	commonResult
}

// Extension represents a neutron extension.
type Extension struct {
	Name        string                  `json:"name"`
	Alias       string                  `json:"alias"`
	Links       []string                `json:"links"`
	Description string                  `json:"description"`
	Updated     gcorecloud.JSONRFC3339Z `json:"updated"`
}

// ExtensionPage is the page returned by a pager when traversing over a
// collection of networks.
type ExtensionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of networks has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ExtensionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ExtensionPage struct is empty.
func (r ExtensionPage) IsEmpty() (bool, error) {
	is, err := ExtractExtensions(r)
	return len(is) == 0, err
}

// ExtractExtension accepts a Page struct, specifically a ExtensionPage struct,
// and extracts the elements into a slice of Extension structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractExtensions(r pagination.Page) ([]Extension, error) {
	var s []Extension
	err := ExtractExtensionsInto(r, &s)
	return s, err
}

func ExtractExtensionsInto(r pagination.Page, v interface{}) error {
	return r.(ExtensionPage).Result.ExtractIntoSlicePtr(v, "results")
}
