package l7policies

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateOpts represents options used to create a l7 policy.
type CreateOpts struct {
	Name             string   `json:"name,omitempty"`
	ListenerID       string   `json:"listener_id" required:"true" validate:"required,uuid"`
	Action           Action   `json:"action" required:"true" validate:"required,enum"`
	Position         int32    `json:"position,omitempty" validate:"gt=-1,omitempty"`
	RedirectHTTPCode int      `json:"redirect_http_code,omitempty"`
	RedirectPoolID   string   `json:"redirect_pool_id,omitempty" validate:"rfe=Action:REDIRECT_TO_POOL"`
	RedirectPrefix   string   `json:"redirect_prefix,omitempty" validate:"rfe=Action:REDIRECT_PREFIX"`
	RedirectURL      string   `json:"redirect_url,omitempty" validate:"rfe=Action:REDIRECT_TO_URL"`
	Tags             []string `json:"tags,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToL7PolicyCreateMap() (map[string]interface{}, error)
}

// ReplaceOpts represents options used to replace a l7 policy.
type ReplaceOpts struct {
	Name             string   `json:"name,omitempty"`
	Action           Action   `json:"action" required:"true" validate:"required,enum"`
	Position         int32    `json:"position,omitempty" validate:"gt=-1,omitempty"`
	RedirectHTTPCode int      `json:"redirect_http_code,omitempty"`
	RedirectPoolID   string   `json:"redirect_pool_id,omitempty" validate:"rfe=Action:REDIRECT_TO_POOL"`
	RedirectPrefix   string   `json:"redirect_prefix,omitempty" validate:"rfe=Action:REDIRECT_PREFIX"`
	RedirectURL      string   `json:"redirect_url,omitempty" validate:"rfe=Action:REDIRECT_TO_URL"`
	Tags             []string `json:"tags,omitempty"`
}

// ReplaceOptsBuilder allows extensions to add additional parameters to the Replace request.
type ReplaceOptsBuilder interface {
	ToL7PolicyReplaceMap() (map[string]interface{}, error)
}

// ToL7PolicyReplaceMap builds a request body from ReplaceOpts.
func (opts ReplaceOpts) ToL7PolicyReplaceMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateRuleOpts represents options used to create a rule for policy.
type CreateRuleOpts struct {
	CompareType CompareType `json:"compare_type" required:"true" validate:"required,enum"`
	Invert      bool        `json:"invert"`
	Key         string      `json:"key,omitempty"`
	Type        RuleType    `json:"type" required:"true" validate:"required,enum"`
	Value       string      `json:"value" required:"true"`
	Tags        []string    `json:"tags,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateRuleOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
}

// ToL7PolicyCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToL7PolicyCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToRuleCreateMap builds a request body from CreateRuleOpts.
func (opts CreateRuleOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// List retrieves list of policies.
func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return L7PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of policies.
func ListAll(c *gcorecloud.ServiceClient) ([]L7Policy, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractL7Polices(page)
}

// Delete accepts a policy id and delete existing policy.
func Delete(c *gcorecloud.ServiceClient, policyID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, policyID), &r.Body, nil)
	return
}

// Create accepts a CreateOpts struct and creates a new policy using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToL7PolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Replace accepts a ReplaceOpts struct and policy id and replaced an existing policy using the values provided.
func Replace(c *gcorecloud.ServiceClient, policyID string, opts ReplaceOptsBuilder) (r tasks.Result) {
	b, err := opts.ToL7PolicyReplaceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(replaceURL(c, policyID), b, &r.Body, nil)
	return
}

// Get retrieves a specific policy based on its unique ID.
func Get(c *gcorecloud.ServiceClient, policyID string) (r GetResult) {
	url := getURL(c, policyID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// GetRule retrieves a specific policy based on its policy id and rule unique ID.
func GetRule(c *gcorecloud.ServiceClient, plid string, rlid string) (r GetRuleResult) {
	url := rulesGetURL(c, plid, rlid)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateRule accepts a policy id and CreateRuleOpts struct and creates a new rule using the values provided.
func CreateRule(c *gcorecloud.ServiceClient, policyID string, opts CreateRuleOptsBuilder) (r tasks.Result) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rulesCreateURL(c, policyID), b, &r.Body, nil)
	return
}

// ReplaceRule accepts a CreateRuleOpts struct, rule id and policy id and replaced an existing rule using the values provided.
func ReplaceRule(c *gcorecloud.ServiceClient, policyID, ruleID string, opts CreateRuleOptsBuilder) (r tasks.Result) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(rulesReplaceURL(c, policyID, ruleID), b, &r.Body, nil)
	return
}

// ListRule accept a policy id and retrieves list of rules.
func ListRule(c *gcorecloud.ServiceClient, policyID string) pagination.Pager {
	url := rulesListURL(c, policyID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllRule accept a policy id and retrieves list of rules.
func ListAllRule(c *gcorecloud.ServiceClient, policyID string) ([]L7Rule, error) {
	page, err := ListRule(c, policyID).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractL7Rules(page)
}

// DeleteRule accepts a policy id, rule id and delete existing rule.
func DeleteRule(c *gcorecloud.ServiceClient, policyID, ruleID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(rulesDeleteURL(c, policyID, ruleID), &r.Body, nil)
	return
}
