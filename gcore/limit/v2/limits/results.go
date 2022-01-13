package limits

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v1/types"
)

type commonResult struct {
	gcorecloud.Result
}

type LimitResponse struct {
	ID              int                       `json:"id"`
	ClientID        int                       `json:"client_id"`
	RequestedLimits Limit                     `json:"requested_limits"`
	Status          types.LimitRequestStatus  `json:"status"`
	CreatedAt       gcorecloud.JSONRFC3339NoZ `json:"created_at"`
}

// Extract is a function that accepts a result and extracts a limit response resource.
func (r commonResult) Extract() (*LimitResponse, error) {
	var s LimitResponse
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LimitResponse.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of an delete operation. Call its ExtractErr to get operation error.
type DeleteResult struct {
	gcorecloud.ErrResult
}
