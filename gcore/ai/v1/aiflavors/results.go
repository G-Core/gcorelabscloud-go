package aiflavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	"github.com/shopspring/decimal"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a instance resource.
func (r commonResult) Extract() (*AIFlavor, error) {
	var s AIFlavor
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// AIFlavorPage is the page returned by a pager when traversing over a
// collection of instances.
type AIFlavorPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of flavors has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AIFlavorPage) NextPageURL() (string, error) {
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
func (r AIFlavorPage) IsEmpty() (bool, error) {
	is, err := ExtractAIFlavors(r)
	return len(is) == 0, err
}

// ExtractFlavor accepts a Page struct, specifically a FlavorPage struct,
// and extracts the elements into a slice of Flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIFlavors(r pagination.Page) ([]AIFlavor, error) {
	var s []AIFlavor
	err := ExtractAIFlavorsInto(r, &s)
	return s, err
}

func ExtractAIFlavorsInto(r pagination.Page, v interface{}) error {
	return r.(AIFlavorPage).Result.ExtractIntoSlicePtr(v, "results")
}


type HardwareDescription struct {
	CPU         string `json:"cpu,omitempty"`
	Disk        string `json:"disk,omitempty"`
	Network     string `json:"network,omitempty"`
	RAM         string `json:"ram,omitempty"`
	Ephemeral   string `json:"ephemeral,omitempty"`
	GPU         string `json:"gpu,omitempty"`
	IPU         string `json:"ipu,omitempty"`
	PoplarCount string `json:"poplar_count,omitempty"`
	SGXEPCSize  string `json:"sgx_epc_size,omitempty"`
}

// Flavor represents a flavor structure.
type AIFlavor struct {
	FlavorID            string               `json:"flavor_id"`
	FlavorName          string               `json:"flavor_name"`
	Disabled            bool                 `json:"disabled"`
	ResourceClass       string               `json:"resource_class"`
	PriceStatus         *string              `json:"price_status,omitempty"`
	CurrencyCode        *gcorecloud.Currency `json:"currency_code,omitempty"`
	PricePerHour        *decimal.Decimal     `json:"price_per_hour,omitempty"`
	PricePerMonth       *decimal.Decimal     `json:"price_per_month,omitempty"`
	HardwareDescription *HardwareDescription `json:"hardware_description,omitempty"`
	RAM                 *int                  `json:"ram,omitempty"`
	VCPUS               *int                  `json:"vcpus,omitempty"`
	Capacity            *int                  `json:"capacity,omitempty"`
}

