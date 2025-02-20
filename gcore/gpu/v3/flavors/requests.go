package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	IncludePrices *bool `q:"include_prices"`
	Disabled      *bool `q:"disabled"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List retrieves list of GPU flavors
func List(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := FlavorsURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
