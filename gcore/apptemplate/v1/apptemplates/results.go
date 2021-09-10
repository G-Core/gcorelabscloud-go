package apptemplates

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a app template resource.
func (r commonResult) Extract() (*AppTemplate, error) {
	var a AppTemplate
	err := r.ExtractInto(&a)
	return &a, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a AppTemplate.
type GetResult struct {
	commonResult
}

// AppTemplatePage is the page returned by a pager when traversing over a
// collection of app templates.
type AppTemplatePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of app templates has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r AppTemplatePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AppTemplatePage struct is empty.
func (r AppTemplatePage) IsEmpty() (bool, error) {
	is, err := ExtractAppTemplates(r)
	return len(is) == 0, err
}

// ExtractAppTemplates accepts a Page struct, specifically a AppTemplatePage struct,
// and extracts the elements into a slice of AppTemplate structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractAppTemplates(r pagination.Page) ([]AppTemplate, error) {
	var s []AppTemplate
	err := ExtractAppTemplatesInto(r, &s)
	return s, err
}

func ExtractAppTemplatesInto(r pagination.Page, v interface{}) error {
	return r.(AppTemplatePage).Result.ExtractIntoSlicePtr(v, "results")
}

type AppTemplate struct {
	ID               string                   `json:"id"`
	OsName           string                   `json:"os_name"`
	Developer        string                   `json:"developer"`
	OsVersion        string                   `json:"os_version"`
	Category         string                   `json:"category"`
	Website          string                   `json:"website"`
	MinVCPUs         *int                     `json:"min_vcpus"`
	DisplayName      string                   `json:"display_name"`
	ImageName        string                   `json:"image_name"`
	Usage            string                   `json:"usage"`
	Description      string                   `json:"description"`
	ShortDescription string                   `json:"short_description"`
	RegionID         *int                     `json:"region_id"`
	MinRam           int                      `json:"min_ram"`
	AppConfig        []map[string]interface{} `json:"app_config"`
	Version          string                   `json:"version"`
	MinDisk          int                      `json:"min_disk"`
}
