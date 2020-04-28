package projects

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/types"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a project resource.
func (r commonResult) Extract() (*Project, error) {
	var s Project
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Project.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Project.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Project.
type UpdateResult struct {
	commonResult
}

// Project represents a project structure.
type Project struct {
	ID          int                       `json:"id"`
	ClientID    int                       `json:"client_id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	State       types.ProjectState        `json:"state"`
	TaskID      *string                   `json:"task_id"`
	CreatedAt   gcorecloud.JSONRFC3339NoZ `json:"created_at"`
}

// ProjectPage is the page returned by a pager when traversing over a
// collection of projects.
type ProjectPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of projects has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ProjectPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ProjectPage struct is empty.
func (r ProjectPage) IsEmpty() (bool, error) {
	is, err := ExtractProjects(r)
	return len(is) == 0, err
}

// ExtractProject accepts a Page struct, specifically a ProjectPage struct,
// and extracts the elements into a slice of Project structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractProjects(r pagination.Page) ([]Project, error) {
	var s []Project
	err := ExtractProjectsInto(r, &s)
	return s, err
}

func ExtractProjectsInto(r pagination.Page, v interface{}) error {
	return r.(ProjectPage).Result.ExtractIntoSlicePtr(v, "results")
}
