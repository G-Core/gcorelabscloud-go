package keypairs

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToKeyPairListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	UserID    string `q:"user_id"`
	ProjectID int    `q:"project_id"`
}

// ToKeyPairListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToKeyPairListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToKeyPairListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return KeyPairPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]KeyPair, error) {
	pages, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractKeyPairs(pages)
}

// Get retrieves a specific keypair based on its name or ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToKeyPairCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a keypair.
type CreateOpts struct {
	Name      string `json:"sshkey_name" required:"true"`
	PublicKey string `json:"public_key,omitempty" required:"true"`
	ProjectID int    `json:"project_id" required:"true"`
}

// ToKeyPairCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToKeyPairCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new keypair using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeyPairCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// Delete accepts a unique ID and deletes the keypair associated with it.
func Delete(c *gcorecloud.ServiceClient, keypairID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, keypairID), nil)
	return
}
