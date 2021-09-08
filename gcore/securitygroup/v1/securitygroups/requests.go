package securitygroups

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/G-Core/gcorelabscloud-go/pagination"
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
	ToSecurityGroupCreateMap() (map[string]interface{}, error)
}

// CreateRuleOptsBuilder allows extensions to add additional parameters to the request.
type CreateRuleOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
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

// ToRuleCreateMap builds a request body from CreateSecurityGroupRuleOpts.
func (opts CreateSecurityGroupRuleOpts) ToRuleCreateMap() (map[string]interface{}, error) {
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
	Instances     []string                `json:"instances,omitempty"`
}

// ToSecurityGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSecurityGroupCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new security group using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecurityGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToSecurityGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a security group.
type UpdateOpts struct {
	Name string `json:"name" required:"true"`
}

// ToSecurityGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSecurityGroupUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing security group using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, securityGroupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSecurityGroupUpdateMap()
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
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(addRulesURL(c, securityGroupID), b, &r.Body, nil)
	return
}

// ListInstances returns page of instances for SG
func ListInstances(c *gcorecloud.ServiceClient, securityGroupID string) pagination.Pager {
	url := listInstancesURL(c, securityGroupID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SecurityGroupInstancesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllInstances returns all instances for SG
func ListAllInstances(c *gcorecloud.ServiceClient, securityGroupID string) ([]instances.Instance, error) {
	page, err := ListInstances(c, securityGroupID).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroupInstances(page)
}

// IDFromName is a convenience function that returns a security group's ID,
// given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	sgs, err := ListAll(client)
	if err != nil {
		return "", err
	}

	for _, s := range sgs {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "security group"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "security group"}
	}
}

// DeepCopyOptsBuilder allows extensions to add additional parameters to the request.
type DeepCopyOptsBuilder interface {
	ToDeepCopyMap() (map[string]interface{}, error)
}

// DeepCopyOpts represents options used to deep copy a security group.
type DeepCopyOpts struct {
	Name string `json:"name" required:"true"`
}

// ToDeepCopyMap builds a request body from DeepCopyOpts.
func (opts DeepCopyOpts) ToDeepCopyMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// DeepCopy accepts a DeepCopyOpts struct and create a deep copy of security group.
func DeepCopy(c *gcorecloud.ServiceClient, securityGroupID string, opts DeepCopyOptsBuilder) (r DeepCopyResult) {
	b, err := opts.ToDeepCopyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(deepCopyURL(c, securityGroupID), b, nil, nil)
	return
}
