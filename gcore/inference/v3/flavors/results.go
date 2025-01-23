package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type Flavor struct {
	Name                 string  `json:"name"`
	Cpu                  float64 `json:"cpu"`
	Memory               float64 `json:"memory"`
	Gpu                  int     `json:"gpu"`
	GpuModel             string  `json:"gpu_model"`
	GpuMemory            float64 `json:"gpu_memory"`
	IsGpuShared          bool    `json:"is_gpu_shared"`
	GpuComputeCapability string  `json:"gpu_compute_capability"`
}

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts an inference flavor resource.
func (r commonResult) Extract() (*Flavor, error) {
	var s Flavor
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type ListResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts an inference flavors resource.
func (r ListResult) Extract() ([]Flavor, error) {
	var s []Flavor
	err := r.ExtractInto(&s)
	return s, err
}

func (r ListResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a InferenceDeployment.
type GetResult struct {
	commonResult
}
