package availablenetworks

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all networks.
func ListAll(client *gcorecloud.ServiceClient) ([]Network, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractNetworks(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
