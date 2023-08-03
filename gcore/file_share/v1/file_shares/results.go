package file_shares

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// Extract is a function that accepts a result and extracts a file share resource.
func (r commonResult) Extract() (*FileShare, error) {
	var s FileShare
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a FileShare.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a FileShare.
type UpdateResult struct {
	commonResult
}

// FileShare represents a file share structure.
type FileShare struct {
	Name             string                          `json:"name"`
	ID               string                          `json:"id"`
	Protocol         string                          `json:"protocol"`
	Status           string                          `json:"status"`
	Size             int                             `json:"size"`
	VolumeType       string                          `json:"volume_type"`
	CreatedAt        *gcorecloud.JSONRFC3339MilliNoZ `json:"created_at"`
	ShareNetworkName string                          `json:"share_network_name"`
	NetworkName      string                          `json:"network_name"`
	SubnetName       string                          `json:"subnet_name"`
	ConnectionPoint  string                          `json:"connection_point"`
	TaskID           *string                         `json:"task_id"`
	CreatorTaskID    *string                         `json:"creator_task_id"`
	ProjectID        int                             `json:"project_id"`
	RegionID         int                             `json:"region_id"`
	Region           string                          `json:"region"`
	Metadata         map[string]interface{}          `json:"metadata"`
}

// FileSharePage is the page returned by a pager when traversing over a
// collection of file shares.
type FileSharePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of file shares has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r FileSharePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FileSharePage struct is empty.
func (r FileSharePage) IsEmpty() (bool, error) {
	is, err := ExtractFileShares(r)
	return len(is) == 0, err
}

// ExtractFileShare accepts a Page struct, specifically a FileSharePage struct,
// and extracts the elements into a slice of FileShare structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFileShares(r pagination.Page) ([]FileShare, error) {
	var s []FileShare
	err := ExtractFileSharesInto(r, &s)
	return s, err
}

func ExtractFileSharesInto(r pagination.Page, v interface{}) error {
	return r.(FileSharePage).Result.ExtractIntoSlicePtr(v, "results")
}

type FileShareTaskResult struct {
	FileShares []string `mapstructure:"file_shares"`
}

func ExtractFileShareIDFromTask(task *tasks.Task) (string, error) {
	var result FileShareTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode file share information in task structure: %w", err)
	}
	if len(result.FileShares) == 0 {
		return "", fmt.Errorf("cannot decode file share information in task structure: %w", err)
	}
	return result.FileShares[0], nil
}

type commonAccessRuleResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a file share access rule resource.
func (r commonAccessRuleResult) Extract() (*AccessRule, error) {
	var s AccessRule
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonAccessRuleResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a FileShare.
type GetAccessRuleResult struct {
	commonAccessRuleResult
}

type CreateAccessRuleResult struct {
	commonAccessRuleResult
}

// FileShare represents a file share structure.
type AccessRule struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	AccessTo    string `json:"access_to"`
	AccessLevel string `json:"access_level"`
}

// AccessRulePage is the page returned by a pager when traversing over a
// collection of file share access rules.
type AccessRulePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of file share access rules has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AccessRulePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AccessRulePage struct is empty.
func (r AccessRulePage) IsEmpty() (bool, error) {
	is, err := ExtractAccessRule(r)
	return len(is) == 0, err
}

// ExtractAccessRule accepts a Page struct, specifically a AccessRulePage struct,
// and extracts the elements into a slice of AccessRule structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAccessRule(r pagination.Page) ([]AccessRule, error) {
	var s []AccessRule
	err := ExtractAccessRuleInto(r, &s)
	return s, err
}

func ExtractAccessRuleInto(r pagination.Page, v interface{}) error {
	return r.(AccessRulePage).Result.ExtractIntoSlicePtr(v, "results")
}

// MetadataPage is the page returned by a pager when traversing over a
// collection of instance metadata objects.
type MetadataPage struct {
	pagination.LinkedPageBase
}

// MetadataResult represents the result of a get operation
type MetadataResult struct {
	commonResult
}

type Metadata struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

func ExtractMetadataInto(r pagination.Page, v interface{}) error {
	return r.(MetadataPage).Result.ExtractIntoSlicePtr(v, "results")
}

// ExtractMetadata accepts a Page struct, specifically a MetadataPage struct,
// and extracts the elements into a slice of securitygroups metadata structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMetadata(r pagination.Page) ([]Metadata, error) {
	var s []Metadata
	err := ExtractMetadataInto(r, &s)
	return s, err
}

// MetadataActionResult represents the result of a create, delete or update operation(no content)
type MetadataActionResult struct {
	gcorecloud.ErrResult
}
