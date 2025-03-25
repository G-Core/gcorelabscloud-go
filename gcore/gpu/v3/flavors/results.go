package flavors

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	"github.com/shopspring/decimal"
)

type commonResult struct {
	gcorecloud.Result
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// PriceDisplayStatus represents the status of price display
type PriceDisplayStatus string

const (
	PriceStatusShow  PriceDisplayStatus = "show"
	PriceStatusHide  PriceDisplayStatus = "hide"
	PriceStatusError PriceDisplayStatus = "error"
)

// Price represents a flavor price structure
type Price struct {
	CurrencyCode  *string             `json:"currency_code"`
	PricePerHour  *decimal.Decimal    `json:"price_per_hour"`
	PricePerMonth *decimal.Decimal    `json:"price_per_month"`
	PriceStatus   *PriceDisplayStatus `json:"price_status"`
}

// GPUInfo represents GPU related information
type GPUInfo struct {
	Model  string `json:"model"`
	Count  int    `json:"count"`
	Memory int    `json:"memory"`
}

// FormatGPU formats GPU information according to the formula
func FormatGPU(model string, count int, memory int) string {
	return fmt.Sprintf("NVIDIA %s-%dGPU (%dGB)", model, count, memory)
}

// HardwareProperties represents GPU hardware properties
type HardwareProperties struct {
	GPUModel        *string `json:"gpu_model"`
	GPUManufacturer *string `json:"gpu_manufacturer"`
	GPUCount        *int    `json:"gpu_count"`
}

// VMFlavor represents a virtual GPU flavor
type VMFlavor struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	VCPUs               int                    `json:"vcpus"`
	RAM                 int                    `json:"ram"`
	Price               *Price                 `json:"price"`
	Architecture        *string                `json:"architecture"`
	Disabled            bool                   `json:"disabled"`
	Capacity            int                    `json:"capacity"`
	GPU                 string                 `json:"gpu"`
	HardwareDescription map[string]interface{} `json:"hardware_description"`
	HardwareProperties  *HardwareProperties    `json:"hardware_properties"`
}

// BMFlavor represents a baremetal GPU flavor
type BMFlavor struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	CPU                 string                 `json:"cpu"`
	RAM                 int                    `json:"ram"`
	Disk                string                 `json:"disk"`
	Network             string                 `json:"network"`
	GPU                 string                 `json:"gpu"`
	Price               *Price                 `json:"price"`
	Architecture        *string                `json:"architecture"`
	Disabled            bool                   `json:"disabled"`
	Capacity            int                    `json:"capacity"`
	HardwareDescription map[string]interface{} `json:"hardware_description"`
	HardwareProperties  *HardwareProperties    `json:"hardware_properties"`
}

// FlavorPage is the page returned by a pager when traversing over a collection of flavors.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a FlavorPage struct is empty.
func (r FlavorPage) IsEmpty() (bool, error) {
	if r.Result.Err != nil {
		return true, r.Result.Err
	}

	var vmFlavors []VMFlavor
	vmErr := r.Result.ExtractIntoSlicePtr(&vmFlavors, "results")
	if vmErr == nil && len(vmFlavors) > 0 {
		return false, nil
	}

	var bmFlavors []BMFlavor
	bmErr := r.Result.ExtractIntoSlicePtr(&bmFlavors, "results")
	if bmErr == nil && len(bmFlavors) > 0 {
		return false, nil
	}

	if vmErr != nil && bmErr != nil {
		return true, vmErr
	}

	return true, nil
}

// ExtractVMFlavors extracts virtual machine flavors
func ExtractVMFlavors(r pagination.Page) ([]VMFlavor, error) {
	var s []VMFlavor
	err := r.(FlavorPage).Result.ExtractIntoSlicePtr(&s, "results")
	return s, err
}

// ExtractBMFlavors extracts baremetal flavors
func ExtractBMFlavors(r pagination.Page) ([]BMFlavor, error) {
	var s []BMFlavor
	err := r.(FlavorPage).Result.ExtractIntoSlicePtr(&s, "results")
	return s, err
}
