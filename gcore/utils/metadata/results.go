package metadata

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

type Metadata struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

type MetadataPage struct {
	pagination.LinkedPageBase
}

// MetadataResult represents the result of a get operation
type MetadataResult struct {
	commonResult
}

func ExtractMetadataInto(r pagination.Page, v interface{}) error {
	return r.(MetadataPage).Result.ExtractIntoSlicePtr(v, "results")
}

// ExtractMetadata accepts a Page struct, specifically a MetadataPage struct,
// and extracts the elements into a slice of securitygroups metadata structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMetadata(r pagination.Page) ([]Metadata, error) {
	var s []Metadata
	err := ExtractMetadataInto(r, &s)
	return s, err
}

// MetadataActionResult represents the result of a creation, delete or update operation(no content)
type MetadataActionResult struct {
	gcorecloud.ErrResult
}

func (r MetadataResult) Extract() (*Metadata, error) {
	var s Metadata
	err := r.ExtractInto(&s)
	return &s, err
}

// IsEmpty checks whether a MetadataPage struct is empty.
func (r MetadataPage) IsEmpty() (bool, error) {
	is, err := ExtractMetadata(r)
	return len(is) == 0, err
}
