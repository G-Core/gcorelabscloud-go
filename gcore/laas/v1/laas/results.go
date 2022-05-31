package laas

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// User represents regenerated laas credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a User resource.
func (r UserResult) Extract() (*User, error) {
	var u User
	err := r.ExtractInto(&u)
	return &u, err
}

func (r UserResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// Status represents laas status
type Status struct {
	Namespace     string `json:"namespace"`
	IsInitialized bool   `json:"is_initialized"`
}

type StatusResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a Status resource.
func (r StatusResult) Extract() (*Status, error) {
	var s Status
	err := r.ExtractInto(&s)
	return &s, err
}

func (r StatusResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// Hosts represents kafka/opensearch hosts url
type Hosts []string

type HostsResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a Hosts resource.
func (r HostsResult) Extract() (*Hosts, error) {
	var h Hosts
	err := r.ExtractInto(&h)
	return &h, err
}

func (r HostsResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "hosts")
}

// Topic represents kafka/opensearch topic
type Topic struct {
	Name string
}

type TopicResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a Topic resource.
func (r TopicResult) Extract() (*Topic, error) {
	var h Topic
	err := r.ExtractInto(&h)
	return &h, err
}

func (r TopicResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// TopicPage is the page returned by a pager when traversing over a
// collection of topics.
type TopicPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of topics has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r TopicPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a Topic struct is empty.
func (r TopicPage) IsEmpty() (bool, error) {
	is, err := ExtractTopics(r)
	return len(is) == 0, err
}

// ExtractTopics accepts a Page struct, specifically a TopicPage struct,
// and extracts the elements into a slice of Topic structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractTopics(r pagination.Page) ([]Topic, error) {
	var s []Topic
	err := ExtractTopicsInto(r, &s)
	return s, err
}

func ExtractTopicsInto(r pagination.Page, v interface{}) error {
	return r.(TopicPage).Result.ExtractIntoSlicePtr(v, "results")
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}
