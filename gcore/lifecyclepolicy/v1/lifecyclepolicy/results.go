package lifecyclepolicy

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type commonResult struct {
	gcorecloud.Result
}

// GetResult represents the result of a Get operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type GetResult struct {
	commonResult
}

// ListResult represents the result of a ListAll operation.
// Call its Extract method to interpret it as a list of LifecyclePolicy.
type ListResult struct {
	commonResult
}

// CreateResult represents the result of a Create operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a Update operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type UpdateResult struct {
	commonResult
}

// AddVolumesResult represents the result of a AddVolumes operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type AddVolumesResult struct {
	commonResult
}

// RemoveVolumesResult represents the result of a RemoveVolumes operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type RemoveVolumesResult struct {
	commonResult
}

// AddSchedulesResult represents the result of a AddSchedules operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type AddSchedulesResult struct {
	commonResult
}

// RemoveSchedulesResult represents the result of a RemoveSchedules operation.
// Call its Extract method to interpret it as a LifecyclePolicy.
type RemoveSchedulesResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a lifecycle policy resource.
func (r commonResult) Extract() (*LifecyclePolicy, error) {
	var rawPolicy rawLifecyclePolicy
	if err := r.Result.ExtractIntoStructPtr(&rawPolicy, ""); err != nil {
		return nil, err
	}
	return rawPolicy.cook()
}

// Extract is a function that accepts a result and extracts a slice of lifecycle policy resources.
func (r ListResult) Extract() ([]LifecyclePolicy, error) {
	var rawPolicies []rawLifecyclePolicy
	err := r.Result.ExtractIntoSlicePtr(&rawPolicies, "results")
	policies := make([]LifecyclePolicy, len(rawPolicies))
	for i, rawPolicy := range rawPolicies {
		p, err := rawPolicy.cook()
		if err != nil {
			return nil, err
		}
		policies[i] = *p
	}
	return policies, err
}
