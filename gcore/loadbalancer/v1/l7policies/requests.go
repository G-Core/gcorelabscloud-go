package l7policies

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type Action string
type RuleType string
type CompareType string

const (
	ActionRedirectToPool Action = "REDIRECT_TO_POOL"
	ActionRedirectToURL  Action = "REDIRECT_TO_URL"
	ActionReject         Action = "REJECT"

	TypeCookie   RuleType = "COOKIE"
	TypeFileType RuleType = "FILE_TYPE"
	TypeHeader   RuleType = "HEADER"
	TypeHostName RuleType = "HOST_NAME"
	TypePath     RuleType = "PATH"

	CompareTypeContains  CompareType = "CONTAINS"
	CompareTypeEndWith   CompareType = "ENDS_WITH"
	CompareTypeEqual     CompareType = "EQUAL_TO"
	CompareTypeRegex     CompareType = "REGEX"
	CompareTypeStartWith CompareType = "STARTS_WITH"
)

// CreateOpts represents options used to create a l7 policy.
type CreateOpts struct {
	Name           string `json:"name,omitempty"`
	ListenerID     string `json:"listener_id" required:"true"`
	Action         Action `json:"action" required:"true"`
	Position       int32  `json:"position,omitempty"`
	Description    string `json:"description,omitempty"`
	RedirectPoolID string `json:"redirect_pool_id,omitempty"`
}

// CreateOpts represents options used to create a rule for policy.
type CreateRuleOpts struct {
	CompareType CompareType `json:"compare_type" required:"true"`
	Invert      bool        `json:"invert,omitempty"`
	Key         string      `json:"key,omitempty"`
	Type        RuleType    `json:"type" required:"true"`
	Value       string      `json:"value" required:"true"`
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

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return L7PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Delete(c *gcorecloud.ServiceClient, policyID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, policyID), &r.Body, nil)
	return
}

func Create(c *gcorecloud.ServiceClient, opts CreateOpts) (r tasks.Result) {
	b, err := opts.ToL7PolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

func CreateRule(c *gcorecloud.ServiceClient, policyID string, opts CreateRuleOpts) (r tasks.Result) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rulescreateURL(c, policyID), b, &r.Body, nil)
	return
}

func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

func GetRule(c *gcorecloud.ServiceClient, plid string, rlid string) (r GetResult) {
	url := rulesgetURL(c, plid, rlid)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
