package extensions

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ExtensionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific extension based on its alias.
func Get(c *gcorecloud.ServiceClient, alias string) (r GetResult) {
	url := getURL(c, alias)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

func ListAll(c *gcorecloud.ServiceClient) ([]Extension, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractExtensions(page)
}
