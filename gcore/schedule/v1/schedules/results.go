package schedules

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a schedule resource.
func (r commonResult) Extract() (*lifecyclepolicy.RawSchedule, error) {
	var s lifecyclepolicy.RawSchedule
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a schedule.
type GetResult struct {
	commonResult
}
