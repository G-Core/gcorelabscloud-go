package bmcapacity

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

type availableNodesResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts amount of available baremetal nodes
func (r availableNodesResult) Extract() (*AvailableNodes, error) {
	var c AvailableNodes
	err := r.ExtractInto(&c)

	return &c, err
}

func (r availableNodesResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// AvailableNodes represents available baremetal nodes
type AvailableNodes struct {
	Capacity map[string]int `json:"capacity"`
}

// GetAvailableNodesResult represents the result of a get operation. Call its Extract
// method to interpret it as a AvailableNodes.
type GetAvailableNodesResult struct {
	availableNodesResult
}
