package secrets

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type InferenceSecret struct {
	Name string     `json:"name"`
	Type string     `json:"type"`
	Data SecretData `json:"data"`
}

type SecretData struct {
	AWSSecretKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

type commonResult struct {
	gcorecloud.Result
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// Extract is a function that accepts a result and extracts a registry credentials resource.
func (r commonResult) Extract() (*InferenceSecret, error) {
	var s InferenceSecret
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
func (r ListResult) Extract() ([]InferenceSecret, error) {
	var s []InferenceSecret
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
