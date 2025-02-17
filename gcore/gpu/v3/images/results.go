package images

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type ImageHwFirmwareType string
type ImageOsType string
type SshKeyType string

const (
	HwFirmwareTypeBios ImageHwFirmwareType = "bios"
	HwFirmwareTypeUefi ImageHwFirmwareType = "uefi"

	OsTypeLinux   ImageOsType = "linux"
	OsTypeWindows ImageOsType = "windows"

	SshKeyAllow    SshKeyType = "allow"
	SshKeyDeny     SshKeyType = "deny"
	SshKeyRequired SshKeyType = "required"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a image resource.
func (r commonResult) Extract() (*Image, error) {
	var s Image
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Image.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Image.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Image.
type UpdateResult struct {
	commonResult
}

type Image struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Status         string                   `json:"status"`
	MinDisk        int                      `json:"min_disk"`
	MinRAM         int                      `json:"min_ram"`
	Size           int                      `json:"size"`
	Visibility     string                   `json:"visibility"`
	Architecture   string                   `json:"architecture"`
	CowFormat      bool                     `json:"cow_format"`
	HwFirmwareType *ImageHwFirmwareType     `json:"hw_firmware_type"`
	OsDistro       *string                  `json:"os_distro"`
	OsType         *ImageOsType             `json:"os_type"`
	OsVersion      *string                  `json:"os_version"`
	SshKey         SshKeyType               `json:"ssh_key"`
	URL            string                   `json:"url"`
	CreatedAt      gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt      *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	CreatorTaskID  *string                  `json:"creator_task_id"`
	TaskID         *string                  `json:"task_id"`
	Metadata       map[string]interface{}   `json:"metadata"`
}

// ImagePage is the page returned by a pager when traversing over a
// collection of images.
type ImagePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of images has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ImagePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ImagePage struct is empty.
func (r ImagePage) IsEmpty() (bool, error) {
	is, err := ExtractImages(r)
	return len(is) == 0, err
}

// ExtractImages accepts a Page struct, specifically a ImagePage struct,
// and extracts the elements into a slice of image structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractImages(r pagination.Page) ([]Image, error) {
	var s []Image
	err := ExtractImagesInto(r, &s)
	return s, err
}

func ExtractImagesInto(r pagination.Page, v interface{}) error {
	return r.(ImagePage).Result.ExtractIntoSlicePtr(v, "results")
}

type ImageTaskResult struct {
	Images []string `json:"images"`
}

func ExtractImageIDFromTask(task *tasks.Task) (string, error) {
	var result ImageTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode image information in task structure: %w", err)
	}
	if len(result.Images) == 0 {
		return "", fmt.Errorf("cannot decode image information in task structure: %w", err)
	}
	return result.Images[0], nil
}
