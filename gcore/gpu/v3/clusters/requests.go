package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
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

// ListClustersOptsBuilder allows extensions to add additional parameters to the List request.
type ListClustersOptsBuilder interface {
	ToListClustersQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Limit  int `q:"limit" validate:"omitempty,gt=0"`
	Offset int `q:"offset" validate:"omitempty,gte=0"`
}

// ToListClustersQuery formats a ListOpts into a query string.
func (opts ListOpts) ToListClustersQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List retrieves list of GPU flavors
func List(client *gcorecloud.ServiceClient, opts ListClustersOptsBuilder) pagination.Pager {
	url := ClustersURL(client)
	if opts != nil {
		query, err := opts.ToListClustersQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// Delete removes a specific GPU cluster by its ID.
func Delete(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
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
