package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"net/http"
)

// RenameClusterOptsBuilder allows extensions to add parameters to rename cluster options.
type RenameClusterOptsBuilder interface {
	ToRenameClusterActionMap() (map[string]interface{}, error)
}

// RenameClusterOpts specifies the parameters for the Rename method.
type RenameClusterOpts struct {
	Name string `json:"name" required:"true" validate:"required"`
}

// Validate checks if the provided options are valid.
func (opts RenameClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToRenameClusterActionMap builds a request body from RenameInstanceOpts.
func (opts RenameClusterOpts) ToRenameClusterActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// Rename changes the name of a GPU cluster.
func Rename(client *gcorecloud.ServiceClient, clusterID string, opts RenameClusterOptsBuilder) (r GetResult) {
	b, err := opts.ToRenameClusterActionMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterURL(client, clusterID)
	_, r.Err = client.Patch(url, b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}
