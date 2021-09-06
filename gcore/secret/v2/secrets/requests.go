package secrets

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

// CreateOptsBuilder allows extensions to add additional parameters to the request.
type CreateOptsBuilder interface {
	ToSecretCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a secret.
type CreateOpts struct {
	Expiration *time.Time  `json:"-"`
	Name       string      `json:"name" required:"true"`
	Payload    PayloadOpts `json:"payload" required:"true"`
}

type PayloadOpts struct {
	CertificateChain string `json:"certificate_chain" required:"true"`
	Certificate      string `json:"certificate" required:"true"`
	PrivateKey       string `json:"private_key" required:"true"`
}

// ToSecretCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSecretCreateMap() (map[string]interface{}, error) {
	result, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Expiration != nil {
		result["expiration"] = opts.Expiration.Format(gcorecloud.RFC3339MilliNoZ)
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
