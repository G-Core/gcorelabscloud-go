package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListClustersOptsBuilder allows extensions to add additional parameters to the List request.
type ListClustersOptsBuilder interface {
	ToListClustersQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Limit  int `q:"limit" validate:"omitempty,gt=0"`
	Offset int `q:"offset" validate:"omitempty,gt=0"`
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
