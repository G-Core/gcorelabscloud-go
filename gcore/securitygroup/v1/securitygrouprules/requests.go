package securitygrouprules

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/securitygroup/v1/securitygroups"
)

// Replace accepts a CreateOpts struct and creates a new security group using the values provided.
func Replace(c *gcorecloud.ServiceClient, ruleID string, opts securitygroups.CreateRuleOptsBuilder) (r securitygroups.CreateRuleResult) {
	b, err := opts.ToSecurityGroupRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, ruleID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the security group associated with it.
func Delete(c *gcorecloud.ServiceClient, securityGroupID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, securityGroupID), nil)
	return
}
