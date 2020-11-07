package l7policies

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a L7Policy.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a l7 policy resource.
func (r commonResult) Extract() (*L7Policy, error) {
	var s L7Policy
	err := r.ExtractInto(&s)
	return &s, err
}

// Extract is a function that accepts a result and extracts a rule policy resource.
func (r commonResult) ExtractRule() (*Rule, error) {
	var s Rule
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// L7Policy represents a policy structure.
type L7Policy struct {
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	Action             string                   `json:"action"`
	ListenerID         string                   `json:"listener_id"`
	RedirectPoolID     string                   `json:"redirect_pool_id"`
	Position           int32                    `json:"position"`
	ProjectID          int32                    `json:"project_id"`
	OperatingStatus    string                   `json:"operating_status"`
	ProvisioningStatus string                   `json:"provisioning_status"`
	CreatedAt          gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt          *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	Rules              []Rule                   `json:"rules"`
}

// Rule represents layer 7 load balancing rule.
type Rule struct {
	ID                 string                   `json:"id"`
	CompareType        string                   `json:"compare_type"`
	Invert             bool                     `json:"invert"`
	Key                string                   `json:"key"`
	OperatingStatus    string                   `json:"operating_status"`
	ProvisioningStatus string                   `json:"provisioning_status"`
	CreatedAt          gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt          *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	Type               string                   `json:"type"`
	Value              string                   `json:"value"`
}

// L7PolicyPage is the page returned by a pager when traversing over a
// collection of l7polices.
type L7PolicyPage struct {
	pagination.LinkedPageBase
}

func (r L7PolicyPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a L7PolicyPage struct is empty.
func (r L7PolicyPage) IsEmpty() (bool, error) {
	is, err := ExtractL7Polices(r)
	return len(is) == 0, err
}

func ExtractL7Polices(r pagination.Page) ([]L7Policy, error) {
	var s []L7Policy
	err := ExtractL7PolicesInto(r, &s)
	return s, err
}

func ExtractL7PolicesInto(r pagination.Page, v interface{}) error {
	return r.(L7PolicyPage).Result.ExtractIntoSlicePtr(v, "results")
}

type L7PolicyTaskResult struct {
	L7Polices []string `json:"l7polices"`
}

func ExtractL7PolicyIDFromTask(task *tasks.Task) (string, error) {
	var result L7PolicyTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode l7policy information in task structure: %w", err)
	}
	if len(result.L7Polices) == 0 {
		return "", fmt.Errorf("cannot decode l7policy information in task structure: %w", err)
	}
	return result.L7Polices[0], nil
}

// L7PolicyPage is the page returned by a pager when traversing over a
// collection of l7polices.
type RulePage struct {
	pagination.LinkedPageBase
}

func (r RulePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RulePage struct is empty.
func (r RulePage) IsEmpty() (bool, error) {
	is, err := ExtractRule(r)
	return len(is) == 0, err
}

func ExtractRule(r pagination.Page) ([]Rule, error) {
	var s []Rule
	err := ExtractRuleInto(r, &s)
	return s, err
}

func ExtractRuleInto(r pagination.Page, v interface{}) error {
	return r.(RulePage).Result.ExtractIntoSlicePtr(v, "results")
}

type RuleTaskResult struct {
	L7Rules []string `json:"l7rules"`
}

func ExtractRuleIDFromTask(task *tasks.Task) (string, error) {
	var result RuleTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode l7rule information in task structure: %w", err)
	}
	if len(result.L7Rules) == 0 {
		return "", fmt.Errorf("cannot decode l7rule information in task structure: %w", err)
	}
	return result.L7Rules[0], nil
}
