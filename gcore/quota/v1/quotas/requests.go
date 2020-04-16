package quotas

import (
	"fmt"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// Get retrieves a specific quota based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// ReplaceOptsBuilder allows extensions to add additional parameters to the Replace request.
type ReplaceOptsBuilder interface {
	ToQuotaReplaceMap() (map[string]interface{}, error)
}

// ReplaceOpts represents options used to create or replace quotas.
type ReplaceOpts struct {
	Quota
}

// Validate
func (opts Quota) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToQuotaReplaceMap builds a request body from ReplaceOpts.
func (opts ReplaceOpts) ToQuotaReplaceMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}
	m, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	if len(opts.ToRequestMap()) != len(m) {
		return nil, fmt.Errorf("all Quota fields should be set")
	}
	return m, nil
}

// Replace accepts a ReplaceOpts struct and creates a new quota using the values provided.
func Replace(c *gcorecloud.ServiceClient, id int, opts ReplaceOptsBuilder) (r UpdateResult) {
	b, err := opts.ToQuotaReplaceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(replaceURL(c, id), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToQuotaUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a quota.
type UpdateOpts struct {
	Quota
}

// ToQuotaUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToQuotaUpdateMap() (map[string]interface{}, error) {
	err := gcorecloud.TranslateValidationError(opts.Validate())
	if err != nil {
		return nil, err
	}
	m := opts.ToRequestMap()
	if len(opts.ToRequestMap()) == 0 {
		return nil, fmt.Errorf("at least one of quota fields should be set")
	}
	return m, nil
}

// Update accepts a UpdateOpts struct and updates an existing quota using the values provided.
func Update(c *gcorecloud.ServiceClient, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, id), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// OwnQuota returns quota for current user
func OwnQuota(c *gcorecloud.ServiceClient) (r GetResult) {
	_, r.Err = c.Get(quotaURL(c), &r.Body, nil)
	return
}
