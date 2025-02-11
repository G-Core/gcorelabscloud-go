package secrets

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// CreateInferenceSecretOptsBuilder allows extensions to add additional parameters to the request.
type CreateInferenceSecretOptsBuilder interface {
	ToInferenceSecretCreateMap() (map[string]interface{}, error)
}

// CreateInferenceSecretOpts represents options used to create a function.
type CreateInferenceSecretOpts struct {
	Name string           `json:"name"`
	Type string           `json:"type"`
	Data CreateSecretData `json:"data"`
}

type CreateSecretData struct {
	AWSSecretKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

// ToInferenceSecretCreateMap builds a request body from CreateInferenceSecretOpts.
func (opts CreateInferenceSecretOpts) ToInferenceSecretCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create registry credential.
func Create(c *gcorecloud.ServiceClient, opts CreateInferenceSecretOptsBuilder) (r GetResult) {
	url := createURL(c)
	b, err := opts.ToInferenceSecretCreateMap()
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
func ListAll(c *gcorecloud.ServiceClient) ([]InferenceSecret, error) {
	var r ListResult
	_, r.Err = c.Get(listURL(c), &r.Body, nil)
	return r.Extract()
}

// Delete accepts a unique name and deletes the registry credential associated with it.
func Delete(c *gcorecloud.ServiceClient, name string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, name), nil)
	return
}

// UpdateInferenceSecretOptsBuilder allows extensions to add additional parameters to the request.
type UpdateInferenceSecretOptsBuilder interface {
	ToInferenceSecretUpdateMap() (map[string]interface{}, error)
}

// UpdateInferenceSecretOpts represents options used to update a function.
type UpdateInferenceSecretOpts struct {
	Type string           `json:"type"`
	Data CreateSecretData `json:"data"`
}

// ToInferenceSecretUpdateMap builds a request body from UpdateInferenceSecretOpts.
func (opts UpdateInferenceSecretOpts) ToInferenceSecretUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update existing registry credential.
func Update(c *gcorecloud.ServiceClient, name string, opts UpdateInferenceSecretOptsBuilder) (r GetResult) {
	url := updateURL(c, name)
	b, err := opts.ToInferenceSecretUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(url, b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}
