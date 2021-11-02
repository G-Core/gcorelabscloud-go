package apitokens

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/types"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToAPITokenListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	RoleID      types.RoleIDType `q:"role,omitempty" validate:"omitempty,enum"`
	IssuedBy    int              `q:"issued_by,omitempty"`
	NotIssuedBy int              `q:"not_issued_by,omitempty"`
	Deleted     bool             `q:"deleted,omitempty"`
}

// ToAPITokenListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAPITokenListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a convenience function that returns all api tokens
func List(c *gcorecloud.ServiceClient, clientID int, opts ListOptsBuilder) (r ListResult) {
	url := listURL(c, clientID)
	if opts != nil {
		query, err := opts.ToAPITokenListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Get retrieves a specific api token based on its unique ID.
func Get(c *gcorecloud.ServiceClient, clientID, tokenID int) (r GetResult) {
	_, r.Err = c.Get(getURL(c, clientID, tokenID), &r.Body, nil)
	return
}

// Delete a specific api token based on its unique ID.
func Delete(c *gcorecloud.ServiceClient, clientID, tokenID int) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, clientID, tokenID), nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToApiTokenCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create an api token.
type CreateOpts struct {
	Name        string                   `json:"name" required:"true" validate:"required"`
	Description string                   `json:"description" required:"true" validate:"required"`
	ExpDate     *gcorecloud.JSONRFC3339Z `json:"exp_date"`
	ClientUser  CreateClientUser         `json:"client_user" required:"true" validate:"required"`
}

type CreateClientUser struct {
	Role ClientRole `json:"role" validate:"dive"`
}

// ToApiTokenCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToApiTokenCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Create creates an APIToken.
func Create(client *gcorecloud.ServiceClient, clientID int, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToApiTokenCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client, clientID), b, &r.Body, nil) // nolint
	return
}
