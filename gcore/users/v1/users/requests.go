package users

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// UserAssignmentOptsBuilder allows extensions to add additional parameters to the user assignment request.
type UserAssignmentOptsBuilder interface {
	ToUserAssignmentMap() (map[string]interface{}, error)
}

// UserAssignmentOpts represents options used to assign role to user.
type UserAssignmentOpts struct {
	ClientID  *int   `json:"client_id"`
	ProjectID *int   `json:"project_id"`
	UserID    int    `json:"user_id" required:"true"`
	Role      string `json:"role" required:"true"`
}

// ToUserAssignmentMap builds a request body from UserAssignmentOpts.
func (opts UserAssignmentOpts) ToUserAssignmentMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// AssignUser accepts a UserAssignmentOpts struct and assigns role to user using the values provided.
func AssignUser(c *gcorecloud.ServiceClient, opts UserAssignmentOptsBuilder) (r UserAssignmentResult) {
	b, err := opts.ToUserAssignmentMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(assignUserRoleURL(c), b, &r.Body, nil)
	return
}

// CreateUserOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateUserOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateUserOpts represents options used to create an user.
type CreateUserOpts struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

// ToUserCreateMap builds a request body from CreateUserOpts.
func (opts CreateUserOpts) ToUserCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateUser accepts a CreateUserOpts struct and creates a new user using the values provided.
func CreateUser(c *gcorecloud.ServiceClient, opts CreateUserOptsBuilder) (r CreateUserResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createUserURL(c), b, &r.Body, nil)
	return
}

// CreateApiTokenOptsBuilder allows extensions to add additional parameters to the CreateApiToken request.
type CreateApiTokenOptsBuilder interface {
	ToApiTokenCreateMap() (map[string]interface{}, error)
}

// CreateApiTokenOpts represents options used to create an api token.
type CreateApiTokenOpts struct {
	Email            string `json:"email" required:"true"`
	Password         string `json:"password" required:"true"`
	TokenName        string `json:"token_name" required:"true"`
	TokenDescription string `json:"token_description" required:"true"`
}

// ToApiTokenCreateMap builds a request body from CreateApiTokenOpts.
func (opts CreateApiTokenOpts) ToApiTokenCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateApiToken accepts a CreateApiTokenOpts struct and creates a new api token using the values provided.
func CreateApiToken(c *gcorecloud.ServiceClient, opts CreateApiTokenOptsBuilder) (r CreateApiTokenResult) {
	b, err := opts.ToApiTokenCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createApiTokenURL(c), b, &r.Body, nil)
	return
}
