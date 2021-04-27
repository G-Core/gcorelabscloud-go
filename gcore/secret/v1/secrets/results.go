package secrets

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a secret resource.
func (r commonResult) Extract() (*Secret, error) {
	var s Secret
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Secret.
type GetResult struct {
	commonResult
}

// Secret represents a secret.
type Secret struct {
	ID           string                       `json:"id"`
	Name         string                       `json:"name"`
	Status       string                       `json:"status"`
	Algorithm    string                       `json:"algorithm"`
	BitLength    int                          `json:"bit_length"`
	ContentTypes map[string]string            `json:"content_types"`
	Mode         string                       `json:"mode"`
	Type         SecretType                   `json:"secret_type"`
	CreatedAt    gcorecloud.JSONRFC3339ZColon `json:"created"`
	UpdatedAt    gcorecloud.JSONRFC3339ZColon `json:"updated"`
	Expiration   gcorecloud.JSONRFC3339ZColon `json:"expiration"`
}

// SecretPage is the page returned by a pager when traversing over a
// collection of secrets.
type SecretPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of secrets has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SecretPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a Secret struct is empty.
func (r SecretPage) IsEmpty() (bool, error) {
	is, err := ExtractSecrets(r)
	return len(is) == 0, err
}

// ExtractSecrets accepts a Page struct, specifically a SecretPage struct,
// and extracts the elements into a slice of Secret structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSecrets(r pagination.Page) ([]Secret, error) {
	var s []Secret
	err := ExtractSecretsInto(r, &s)
	return s, err
}

func ExtractSecretsInto(r pagination.Page, v interface{}) error {
	return r.(SecretPage).Result.ExtractIntoSlicePtr(v, "results")
}

type SecretTaskResult struct {
	Secrets []string `json:"secrets"`
}

func ExtractSecretIDFromTask(task *tasks.Task) (string, error) {
	var result SecretTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode secret information in task structure: %w", err)
	}

	if len(result.Secrets) == 0 {
		return "", fmt.Errorf("empty secrets in task structure: %w", err)
	}
	return result.Secrets[0], nil
}
