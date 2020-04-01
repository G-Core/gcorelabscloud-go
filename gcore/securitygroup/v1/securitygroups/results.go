package securitygroups

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/securitygroup/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a security group resource.
func (r commonResult) Extract() (*SecurityGroup, error) {
	var s SecurityGroup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a SecurityGroup.
type CreateResult struct {
	commonResult
}

// CreateRuleResult represents the result of a create operation. Call its Extract
// method to interpret it as a SecurityGroupRule.
type CreateRuleResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a security group rule resource.
func (r CreateRuleResult) Extract() (*SecurityGroupRule, error) {
	var s SecurityGroupRule
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a SecurityGroup.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a SecurityGroup.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	gcorecloud.ErrResult
}

// SecurityGroup represents a security group.
type SecurityGroup struct {
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	ID                 string                   `json:"id"`
	CreatedAt          gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt          *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	RevisionNumber     int                      `json:"revision_number"`
	SecurityGroupRules []SecurityGroupRule      `json:"security_group_rules"`
	ProjectID          int                      `json:"project_id"`
	RegionID           int                      `json:"region_id"`
	Region             string                   `json:"region"`
}

// SecurityGroupRule represents a security group rule.
type SecurityGroupRule struct {
	ID              string                   `json:"id"`
	SecurityGroupID string                   `json:"security_group_id"`
	RemoteGroupID   *string                  `json:"remote_group_id"`
	Direction       types.RuleDirection      `json:"direction"`
	EtherType       *types.EtherType         `json:"ethertype"`
	Protocol        *types.Protocol          `json:"protocol"`
	PortRangeMax    *int                     `json:"port_range_max"`
	PortRangeMin    *int                     `json:"port_range_min"`
	Description     *string                  `json:"description"`
	RemoteIPPrefix  *string                  `json:"remote_ip_prefix"`
	CreatedAt       gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt       *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	RevisionNumber  int                      `json:"revision_number"`
}

// SecurityGroupPage is the page returned by a pager when traversing over a
// collection of security groups.
type SecurityGroupPage struct {
	pagination.LinkedPageBase
}

// SecurityGroupInstancesPage is the page returned by a pager when traversing over a
// collection of security group instances.
type SecurityGroupInstancesPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of security groups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SecurityGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// NextPageURL is invoked when a paginated collection of security group instances has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SecurityGroupInstancesPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SecurityGroupPage struct is empty.
func (r SecurityGroupPage) IsEmpty() (bool, error) {
	is, err := ExtractSecurityGroups(r)
	return len(is) == 0, err
}

// IsEmpty checks whether a SecurityGroupInstancesPage struct is empty.
func (r SecurityGroupInstancesPage) IsEmpty() (bool, error) {
	is, err := ExtractSecurityGroupInstances(r)
	return len(is) == 0, err
}

// ExtractSecurityGroup accepts a Page struct, specifically a SecurityGroupPage struct,
// and extracts the elements into a slice of SecurityGroup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSecurityGroups(r pagination.Page) ([]SecurityGroup, error) {
	var s []SecurityGroup
	err := ExtractSecurityGroupsInto(r, &s)
	return s, err
}

// ExtractSecurityGroupInstances accepts a Page struct, specifically a SecurityGroupInstancesPage struct,
// and extracts the elements into a slice of Instance structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSecurityGroupInstances(r pagination.Page) ([]instances.Instance, error) {
	var s []instances.Instance
	err := ExtractSecurityGroupInstancesInto(r, &s)
	return s, err
}

func ExtractSecurityGroupsInto(r pagination.Page, v interface{}) error {
	return r.(SecurityGroupPage).Result.ExtractIntoSlicePtr(v, "results")
}

func ExtractSecurityGroupInstancesInto(r pagination.Page, v interface{}) error {
	return r.(SecurityGroupInstancesPage).Result.ExtractIntoSlicePtr(v, "results")
}
