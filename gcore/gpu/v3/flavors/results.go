package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// Flavor represents a GPU flavor
type Flavor struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	VCPUs               int               `json:"vcpus"`
	RAM                 int               `json:"ram"`
	Price               map[string]any    `json:"price"`
	Architecture        *string           `json:"architecture"`
	Disabled            bool              `json:"disabled"`
	Capacity            int               `json:"capacity"`
	HardwareDescription map[string]string `json:"hardware_description"`
}

// FlavorPage is the page returned by a pager when traversing over a collection of flavors.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a FlavorPage struct is empty.
func (r FlavorPage) IsEmpty() (bool, error) {
	flavors, err := ExtractFlavors(r)
	return len(flavors) == 0, err
}

// ExtractFlavors accepts a Page struct, specifically a FlavorPage struct,
// and extracts the elements into a slice of Flavor structs.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s []Flavor
	err := ExtractFlavorsInto(r, &s)
	return s, err
}

// ExtractFlavorsInto similar to ExtractInto but operates on a List of Flavors
func ExtractFlavorsInto(r pagination.Page, v interface{}) error {
	return r.(FlavorPage).Result.ExtractIntoSlicePtr(v, "results")
}

// Extract will get the Flavor object from the commonResult
func (r commonResult) Extract() (*Flavor, error) {
	var s Flavor
	err := r.ExtractInto(&s)
	return &s, err
}
