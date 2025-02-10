package credentials

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type RegistryCredentials struct {
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RegistryURL string `json:"registry_url"`
}

type commonResult struct {
	gcorecloud.Result
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// Extract is a function that accepts a result and extracts a registry credentials resource.
func (r commonResult) Extract() (*RegistryCredentials, error) {
	var s RegistryCredentials
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type ListResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a registry credentials resource.
func (r ListResult) Extract() ([]RegistryCredentials, error) {
	var s []RegistryCredentials
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
