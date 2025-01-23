package credentials

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// CreateRegistryCredentialOptsBuilder allows extensions to add additional parameters to the request.
type CreateRegistryCredentialOptsBuilder interface {
	ToRegistryCredentialCreateMap() (map[string]interface{}, error)
}

// CreateRegistryCredentialOpts represents options used to create a function.
type CreateRegistryCredentialOpts struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RegistryURL string `json:"registry_url"`
}

// ToRegistryCredentialCreateMap builds a request body from CreateRegistryCredentialOpts.
func (opts CreateRegistryCredentialOpts) ToRegistryCredentialCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create registry credential.
func Create(c *gcorecloud.ServiceClient, opts CreateRegistryCredentialOptsBuilder) (r GetResult) {
	url := createURL(c)
	b, err := opts.ToRegistryCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// Get registry credential.
func Get(c *gcorecloud.ServiceClient, name string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, name), &r.Body, nil)
	return
}

// ListAll registry credentials.
func ListAll(c *gcorecloud.ServiceClient) ([]RegistryCredentials, error) {
	var r ListResult
	_, r.Err = c.Get(listURL(c), &r.Body, nil)
	return r.Extract()
}

// Delete accepts a unique name and deletes the registry credential associated with it.
func Delete(c *gcorecloud.ServiceClient, name string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, name), nil)
	return
}

// UpdateRegistryCredentialOptsBuilder allows extensions to add additional parameters to the request.
type UpdateRegistryCredentialOptsBuilder interface {
	ToRegistryCredentialUpdateMap() (map[string]interface{}, error)
}

// UpdateRegistryCredentialOpts represents options used to update a function.
type UpdateRegistryCredentialOpts struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	RegistryURL string `json:"registry_url"`
}

// ToRegistryCredentialUpdateMap builds a request body from UpdateRegistryCredentialOpts.
func (opts UpdateRegistryCredentialOpts) ToRegistryCredentialUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update existing registry credential.
func Update(c *gcorecloud.ServiceClient, name string, opts UpdateRegistryCredentialOptsBuilder) (r GetResult) {
	url := updateURL(c, name)
	b, err := opts.ToRegistryCredentialUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(url, b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}
