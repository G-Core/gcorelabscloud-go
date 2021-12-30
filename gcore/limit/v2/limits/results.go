package limits

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// DeleteResult represents the result of an delete operation. Call its ExtractErr to get operation error.
type DeleteResult struct {
	gcorecloud.ErrResult
}
