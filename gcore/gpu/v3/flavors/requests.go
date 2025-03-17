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
	// IncludePrices true to include prices in the response, defaults to false
	IncludePrices *bool `q:"include_prices"`
	// HideDisabled true to hide disabled flavors, defaults to false
	HideDisabled *bool `q:"hide_disabled"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListVirtual retrieves list of virtual GPU flavors
func ListVirtual(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	// Build complete URL for virtual flavors - client already has the path & project/region set
	url := client.ServiceURL(flavorsPath)

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

// ListBaremetal retrieves list of baremetal GPU flavors
func ListBaremetal(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	// Build complete URL for baremetal flavors - client already has the path & project/region set
	url := client.ServiceURL(flavorsPath)

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
