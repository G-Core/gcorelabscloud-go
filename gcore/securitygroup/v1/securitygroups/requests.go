package securitygroups

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/securitygroup/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SecurityGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific security group based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the request.
type CreateOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateRuleOptsBuilder allows extensions to add additional parameters to the request.
type CreateRuleOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateSecurityGroupRuleOpts represents options used to create a security group rule.
type CreateSecurityGroupRuleOpts struct {
	Direction       types.RuleDirection `json:"direction" required:"true"`
	EtherType       types.EtherType     `json:"ethertype,omitempty" required:"true"`
	Protocol        types.Protocol      `json:"protocol,omitempty" required:"true"`
	SecurityGroupID *string             `json:"security_group_id,omitempty"`
	RemoteGroupID   *string             `json:"remote_group_id,omitempty"`
	PortRangeMax    *int                `json:"port_range_max,omitempty"`
	PortRangeMin    *int                `json:"port_range_min,omitempty"`
	Description     *string             `json:"description,omitempty"`
	RemoteIPPrefix  *string             `json:"remote_ip_prefix,omitempty"`
}

// ToFloatingIPCreateMap builds a request body from CreateSecurityGroupRuleOpts.
func (opts CreateSecurityGroupRuleOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateSecurityGroupOpts represents options used to create a security group.
type CreateSecurityGroupOpts struct {
	Name               string                        `json:"name" required:"true"`
	Description        *string                       `json:"description,omitempty"`
	SecurityGroupRules []CreateSecurityGroupRuleOpts `json:"security_group_rules"`
}

// CreateOpts represents options used to create a security group.
type CreateOpts struct {
	SecurityGroup CreateSecurityGroupOpts `json:"security_group" required:"true"`
	Instances     []string                `json:"instances" required:"true"`
}

// ToFloatingIPCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new security group using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToFloatingIPUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a security group.
type UpdateOpts struct {
	Name string `json:"name" required:"true"`
}

// ToFloatingIPUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToFloatingIPUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing security group using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, securityGroupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFloatingIPUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, securityGroupID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the security group associated with it.
func Delete(c *gcorecloud.ServiceClient, securityGroupID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, securityGroupID), nil)
	return
}

// ListAll returns all SGs
func ListAll(c *gcorecloud.ServiceClient) ([]SecurityGroup, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroups(page)
}

// AddRule accepts a CreateSecurityGroupRuleOpts struct and add rule to existed group.
func AddRule(c *gcorecloud.ServiceClient, securityGroupID string, opts CreateRuleOptsBuilder) (r CreateRuleResult) {
	b, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(addRulesURL(c, securityGroupID), b, &r.Body, nil)
	return
}
