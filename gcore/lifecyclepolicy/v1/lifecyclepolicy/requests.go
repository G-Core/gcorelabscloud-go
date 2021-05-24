package lifecyclepolicy

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Get retrieves a lifecycle policy with specified unique id.
// If present, opts are used to construct query parameters.
func Get(c *gcorecloud.ServiceClient, id int, opts GetOpts) (r GetResult) {
	url := getURL(c, id)
	query, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Get(url+query.String(), &r.Body, nil)
	return
}

// ListAll returns all lifecycle policies.
// If present, opts are used to construct query parameters.
func ListAll(c *gcorecloud.ServiceClient, opts ListOpts) (r ListResult) {
	url := listURL(c)
	query, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Get(url+query.String(), &r.Body, nil)
	return
}

// Delete deletes a lifecycle policy with specified unique id.
func Delete(c *gcorecloud.ServiceClient, id int) (err error) {
	url := deleteURL(c, id)
	_, err = c.Delete(url, nil)
	return
}

// Create creates a lifecycle policy.
// opts are used to construct request body.
func Create(c *gcorecloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	url := createURL(c)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// Update updates a lifecycle policy with specified unique id.
// opts are used to construct request body.
func Update(c *gcorecloud.ServiceClient, id int, opts UpdateOpts) (r UpdateResult) {
	url := updateURL(c, id)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(url, b, &r.Body, nil)
	return
}

// AddVolumes adds volumes to a lifecycle policy with specified unique id.
// opts are used to construct request body.
func AddVolumes(c *gcorecloud.ServiceClient, id int, opts AddVolumesOpts) (r AddVolumesResult) {
	url := addVolumesURL(c, id)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(url, b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{200}})
	return
}

// RemoveVolumes removes volumes from a lifecycle policy with specified unique id.
// opts are used to construct request body.
func RemoveVolumes(c *gcorecloud.ServiceClient, id int, opts RemoveVolumesOpts) (r RemoveVolumesResult) {
	url := removeVolumesURL(c, id)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(url, b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{200}})
	return
}

// AddSchedules adds schedules to lifecycle policy with specified unique id.
// opts are used to construct request body.
func AddSchedules(c *gcorecloud.ServiceClient, id int, opts AddSchedulesOpts) (r AddSchedulesResult) {
	url := addSchedulesURL(c, id)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// RemoveSchedules removes schedules from a lifecycle policy with specified unique id.
// opts are used to construct request body.
func RemoveSchedules(c *gcorecloud.ServiceClient, id int, opts RemoveSchedulesOpts) (r RemoveSchedulesResult) {
	url := removeSchedulesURL(c, id)
	b, err := ValidateAndBuildRequestBody(opts)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}
