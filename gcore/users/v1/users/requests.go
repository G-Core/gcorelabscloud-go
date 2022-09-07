package users

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create an user.
type CreateOpts struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

// ToUserCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new user using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}
