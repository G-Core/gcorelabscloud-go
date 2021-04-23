package servergroups

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the request.
type CreateOptsBuilder interface {
	ToServerGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a security group.
type CreateOpts struct {
	Name   string            `json:"name" required:"true"`
	Policy ServerGroupPolicy `json:"policy" required:"true" validate:"enum"`
}

// ToServerGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToServerGroupCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new server group using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r GetResult) {
	b, err := opts.ToServerGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ServerGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific server group based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the server group associated with it.
func Delete(c *gcorecloud.ServiceClient, securityGroupID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, securityGroupID), nil)
	return
}

// ListAll returns all SGs
func ListAll(c *gcorecloud.ServiceClient) ([]ServerGroup, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractServerGroups(page)
}
