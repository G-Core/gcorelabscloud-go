package lbflavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// Flavor represent loadbalancer's flavor structure
type Flavor struct {
	PricePerMonth       float64             `json:"price_per_month,omitempty"`
	ResourceClass       string              `json:"resource_class,omitempty"`
	VCPUs               int                 `json:"vcpus,omitempty"`
	FlavorName          string              `json:"flavor_name,omitempty"`
	HardwareDescription HardwareDescription `json:"hardware_description"`
	CurrencyCode        string              `json:"currency_code,omitempty"`
	PriceStatus         string              `json:"price_status,omitempty"`
	PricePerHour        float64             `json:"price_per_hour,omitempty"`
	RAM                 int                 `json:"ram,omitempty"`
	FlavorID            string              `json:"flavor_id,omitempty"`
}

// HardwareDescription represent flavor's hardware description structure
type HardwareDescription struct {
	SGXEPCSize string `json:"sgx_epc_size"`
	CPU        string `json:"cpu"`
	Disk       string `json:"disk"`
	Network    string `json:"network"`
	GPU        string `json:"gpu"`
	RAM        string `json:"ram"`
}

// FlavorPage is the page returned by a pager when traversing over a
// collection of loadbalancer flavors.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of loadbalancer flavors has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FlavorPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FlavorPage struct is empty.
func (r FlavorPage) IsEmpty() (bool, error) {
	is, err := ExtractFlavors(r)
	return len(is) == 0, err
}

// ExtractFlavors accepts a Page struct, specifically a FlavorPage struct,
// and extracts the elements into a slice of Flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s []Flavor
	err := ExtractFlavorsInto(r, &s)
	return s, err
}

func ExtractFlavorsInto(r pagination.Page, v interface{}) error {
	return r.(FlavorPage).Result.ExtractIntoSlicePtr(v, "results")
}
