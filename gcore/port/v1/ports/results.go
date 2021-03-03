package ports

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a security group resource.
func (r commonResult) Extract() (*instances.Interface, error) {
	var s instances.Interface
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// UpdateResult represents the result of a enable or disable operation. Call its Extract
// method to interpret it as a Interface.
type UpdateResult struct {
	commonResult
}
