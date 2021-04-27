package secrets

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the request.
type CreateOptsBuilder interface {
	ToSecretCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a secret.
type CreateOpts struct {
	Algorithm              *string    `json:"algorithm,omitempty"`
	BitLength              *int       `json:"bit_length,omitempty"`
	Expiration             *time.Time `json:"-"`
	Name                   string     `json:"name" required:"true"`
	Mode                   *string    `json:"mode,omitempty"`
	Type                   SecretType `json:"secret_type" required:"true" validate:"enum"`
	Payload                string     `json:"payload" required:"true"`
	PayloadContentEncoding string     `json:"payload_content_encoding" required:"true"`
	PayloadContentType     string     `json:"payload_content_type" required:"true"`
}

// ToSecretCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSecretCreateMap() (map[string]interface{}, error) {
	result, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Expiration != nil {
		result["expiration"] = opts.Expiration.Format(gcorecloud.RFC3339NoZ)
	}
	return result, nil
}

// Create accepts a CreateOpts struct and creates a new secret using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToSecretCreateMap()
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
		return SecretPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific secret based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the secret associated with it.
func Delete(c *gcorecloud.ServiceClient, securityGroupID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, securityGroupID), &r.Body, nil)
	return
}

// ListAll returns all secrets
func ListAll(c *gcorecloud.ServiceClient) ([]Secret, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractSecrets(page)
}
