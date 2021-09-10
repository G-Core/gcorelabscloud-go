package apptemplates

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// List retrieves list of app templates
func List(c *gcorecloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(c, rootURL(c), func(r pagination.PageResult) pagination.Page {
		return AppTemplatePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of app templates
func ListAll(c *gcorecloud.ServiceClient) ([]AppTemplate, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractAppTemplates(page)
}

// Get retrieves a specific app template based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := resourceURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
