package aiflavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)


type ListOptsBuilder interface {
	ToAIFlavorListQuery() (string, error)
}

type AIFlavorListOpts struct {
	Disabled        bool `q:"disabled"`
	IncludeCapacity bool `q:"include_capacity"`
	IncludePrices   bool `q:"include_prices"`
}

// ToAIFlavorListQuery formats a AIFlavorListOpts into a query string.
func (opts AIFlavorListOpts) ToAIFlavorListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}
		

// List retrieves list of AI flavors
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listAIFlavorsURL(c)
	if opts != nil {
		query, err := opts.ToAIFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AIFlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of all AI flavors
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]AIFlavor, error) {
	results, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractAIFlavors(results)
}