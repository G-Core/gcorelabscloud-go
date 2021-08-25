package lbflavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll returns all LB flavors
func ListAll(c *gcorecloud.ServiceClient) ([]Flavor, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractFlavors(page)
}
