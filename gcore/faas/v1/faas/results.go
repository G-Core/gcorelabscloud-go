package faas

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// Namespace represents FaaS namespace.
type Namespace struct {
	Name                  string                   `json:"name"`
	Description           string                   `json:"description"`
	Status                string                   `json:"status"`
	Envs                  map[string]string        `json:"envs"`
	Functions             []Function               `json:"functions"`
	FunctionsDeployStatus DeployStatus             `json:"functions_deploy_status"`
	CreatedAt             gcorecloud.JSONRFC3339ZZ `json:"created_at"`
}

type DeployStatus struct {
	Total int `json:"total"`
	Ready int `json:"ready"`
}

type NamespaceResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a NamespaceResult resource.
func (r NamespaceResult) Extract() (*Namespace, error) {
	var n Namespace
	err := r.ExtractInto(&n)
	return &n, err
}

func (r NamespaceResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// NamespacePage is the page returned by a pager when traversing over a
// collection of namespace.
type NamespacePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of namespaces has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r NamespacePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a Namespace struct is empty.
func (r NamespacePage) IsEmpty() (bool, error) {
	is, err := ExtractNamespaces(r)
	return len(is) == 0, err
}

// ExtractNamespaces accepts a Page struct, specifically a NamespacePage struct,
// and extracts the elements into a slice of Namespace structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNamespaces(r pagination.Page) ([]Namespace, error) {
	var s []Namespace
	err := ExtractNamespacesInto(r, &s)
	return s, err
}

func ExtractNamespacesInto(r pagination.Page, v interface{}) error {
	return r.(NamespacePage).Result.ExtractIntoSlicePtr(v, "results")
}

// Function represents FaaS function.
type Function struct {
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	BuildMessage string                   `json:"build_message"`
	BuildStatus  string                   `json:"build_status"`
	Status       string                   `json:"status"`
	DeployStatus DeployStatus             `json:"deploy_status"`
	Envs         map[string]string        `json:"envs"`
	Runtime      string                   `json:"runtime"`
	Timeout      int                      `json:"timeout"`
	Flavor       string                   `json:"flavor"`
	Autoscaling  FunctionAutoscaling      `json:"autoscaling"`
	CodeText     string                   `json:"code_text"`
	MainMethod   string                   `json:"main_method"`
	Endpoint     string                   `json:"endpoint"`
	CreatedAt    gcorecloud.JSONRFC3339ZZ `json:"created_at"`
}

type FunctionAutoscaling struct {
	MinInstances int `json:"min_instances"`
	MaxInstances int `json:"max_instances"`
}

type FunctionResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a FunctionResult resource.
func (r FunctionResult) Extract() (*Function, error) {
	var f Function
	err := r.ExtractInto(&f)
	return &f, err
}

func (r FunctionResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// FunctionPage is the page returned by a pager when traversing over a
// collection of functions.
type FunctionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of namespaces has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (f FunctionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := f.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a Function struct is empty.
func (f FunctionPage) IsEmpty() (bool, error) {
	is, err := ExtractFunctions(f)
	return len(is) == 0, err
}

// ExtractFunctions accepts a Page struct, specifically a FunctionPage struct,
// and extracts the elements into a slice of Function structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFunctions(p pagination.Page) ([]Function, error) {
	var f []Function
	err := ExtractFunctionsInto(p, &f)
	return f, err
}

func ExtractFunctionsInto(p pagination.Page, v interface{}) error {
	return p.(FunctionPage).Result.ExtractIntoSlicePtr(v, "results")
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}
