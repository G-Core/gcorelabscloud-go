package ddos

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type FieldType string

const (
	StringField = FieldType("str")
	IntField    = FieldType("int")
	BoolField   = FieldType("bool")
)

type accessStatusResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a DDoS Protection access status resource
func (s accessStatusResult) Extract() (*AccessStatus, error) {
	var st AccessStatus
	err := s.ExtractInto(&st)

	return &st, err
}

func (s accessStatusResult) ExtractInto(v interface{}) error {
	return s.Result.ExtractIntoStructPtr(v, "")
}

// AccessStatus represents DDoS Protection service access status
type AccessStatus struct {
	HTTPCode     int    `json:"http_code"`
	IsAccessible bool   `json:"is_accessible"`
	Message      string `json:"message"`
}

// GetAccessStatusResult represents the result of a get operation. Call its Extract
// method to interpret it as a DDoS Protection access status
type GetAccessStatusResult struct {
	accessStatusResult
}

type regionCoverageResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts an information if provided region can be covered by the Advanced DDoS protection features
func (r regionCoverageResult) Extract() (*RegionCoverage, error) {
	var rc RegionCoverage
	err := r.ExtractInto(&rc)

	return &rc, err
}

func (r regionCoverageResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// RegionCoverage represents an information about coverage of provided region by the Advanced DDoS protection features
type RegionCoverage struct {
	IsCovered bool `json:"is_covered"`
}

// CheckRegionCoverageResult represents the result of a get operation. Call its Extract
// method to interpret it as a region coverage by the DDoS Protection features.
type CheckRegionCoverageResult struct {
	regionCoverageResult
}

// ProfileTemplate represents DDoS protection profile template
type ProfileTemplate struct {
	Fields      []TemplateField `json:"fields"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ID          int             `json:"id"`
}

// TemplateField represents additional fields for protection profile template
type TemplateField struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Required    bool      `json:"required"`
	FieldType   FieldType `json:"field_type"`
	Default     *string   `json:"default,omitempty"`
}

// Profile represents active client DDoS protection profile
type Profile struct {
	ID              int            `json:"id"`
	Options         Options        `json:"options"`
	IPAddress       string         `json:"ip_address"`
	Fields          []ProfileField `json:"fields"`
	ProfileTemplate int            `json:"profile_template"`
}

// Options represent options of active client DDoS protection profile
type Options struct {
	Price  string `json:"price"`
	BGP    bool   `json:"bgp"`
	Active bool   `json:"active"`
}

// ProfileField represent fields of active client DDoS protection profile
type ProfileField struct {
	ID          int       `json:"id,omitempty"`
	Value       string    `json:"value,omitempty" required:"true" validate:"required,max=500"`
	Description string    `json:"description,omitempty"`
	Name        string    `json:"name,omitempty"`
	Required    bool      `json:"required,omitempty"`
	FieldType   FieldType `json:"field_type,omitempty"`
	Default     *string   `json:"default,omitempty"`
	BaseField   int       `json:"base_field,omitempty" required:"true" validate:"required"`
}

// ProfileTemplatesPage is the page returned by a pager when traversing over a
// collection of profile templates.
type ProfileTemplatesPage struct {
	pagination.LinkedPageBase
}

// ProfilesPage is the page returned by a pager when traversing over a
// collection of profile templates.
type ProfilesPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of profile templates has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ProfileTemplatesPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of profiles has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ProfilesPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ProfileTemplatesPage struct is empty.
func (r ProfileTemplatesPage) IsEmpty() (bool, error) {
	is, err := ExtractProfileTemplates(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a ProfilesPage struct is empty.
func (r ProfilesPage) IsEmpty() (bool, error) {
	is, err := ExtractProfiles(r)
	return len(is) == 0, err
}

// ExtractProfileTemplates accepts a Page struct, specifically a ProfileTemplatesPage struct,
// and extracts the elements into a slice of instance structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractProfileTemplates(r pagination.Page) ([]ProfileTemplate, error) {
	var s []ProfileTemplate
	err := ExtractProfileTemplatesInto(r, &s)
	return s, err
}

// ExtractProfiles accepts a Page struct, specifically a ProfilesPage struct,
// and extracts the elements into a slice of instance structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractProfiles(r pagination.Page) ([]Profile, error) {
	var s []Profile
	err := ExtractProfilesInto(r, &s)
	return s, err
}

func ExtractProfileTemplatesInto(r pagination.Page, v interface{}) error {
	return r.(ProfileTemplatesPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractProfilesInto(r pagination.Page, v interface{}) error {
	return r.(ProfilesPage).Result.ExtractIntoSlicePtr(v, "results")
}

type ProfileTaskResult struct {
	Profiles []string `json:"ddos_profiles"`
}

func ExtractProfileIDFromTask(task *tasks.Task) (string, error) {
	var result ProfileTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode DDoS protection profile ID in task structure: %w", err)
	}

	if len(result.Profiles) == 0 {
		return "", fmt.Errorf("cannot decode DDoS protection profile ID in task structure: %w", err)
	}

	return result.Profiles[0], nil
}
