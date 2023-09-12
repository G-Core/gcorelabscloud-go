package aiimages

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a instance resource.
func (r commonResult) Extract() (*AIImage, error) {
	var s AIImage
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// AIFlavorPage is the page returned by a pager when traversing over a
// collection of instances.
type AIImagePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of flavors has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AIImagePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a FlavorPage struct is empty.
func (r AIImagePage) IsEmpty() (bool, error) {
	is, err := ExtractAIImages(r)
	return len(is) == 0, err
}

// ExtractFlavor accepts a Page struct, specifically a FlavorPage struct,
// and extracts the elements into a slice of Flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAIImages(r pagination.Page) ([]AIImage, error) {
	var s []AIImage
	err := ExtractAIImagesInto(r, &s)
	return s, err
}

func ExtractAIImagesInto(r pagination.Page, v interface{}) error {
	return r.(AIImagePage).Result.ExtractIntoSlicePtr(v, "results")
}

type AIImage struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description,omitempty"`
	Status        string                   `json:"status"`
	Visibility    string                   `json:"visibility"`
	MinDisk       int                      `json:"min_disk"`
	MinRAM        int                      `json:"min_ram"`
	OsDistro      string                   `json:"os_distro"`
	OsType        string                   `json:"os_type"`
	OsVersion     string                   `json:"os_version"`
	DisplayOrder  int                      `json:"display_order,omitempty"`
	CreatedAt     gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt     *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	SshKey        string                   `json:"ssh_key,omitempty"`
	Size          int                      `json:"size"`
	CreatorTaskID *string                  `json:"creator_task_id,omitempty"`
	TaskID        *string                  `json:"task_id"`
	Region        string                   `json:"region"`
	RegionID      int                      `json:"region_id"`
	ProjectID     int                      `json:"project_id"`
	DiskFormat    string                   `json:"disk_format"`
	IsBaremetal   bool                     `json:"is_baremetal,omitempty"`
	HwFirmareType string                   `json:"hw_firmware_type,omitempty"`
	HwMachineType string                   `json:"hw_machine_type,omitempty"`
	Architecture  string                   `json:"architecture,omitempty"`
	Metadata      []metadata.Metadata      `json:"metadata_detailed"`
}
