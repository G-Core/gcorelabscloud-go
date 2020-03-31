package availablefloatingips

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/floatingip/v1/floatingips"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return floatingips.FloatingIPPage{
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
		}
	})
}

// ListAll returns all floating IPs
func ListAll(c *gcorecloud.ServiceClient) ([]floatingips.FloatingIPDetail, error) {
	page, err := List(c).AllPages()
	if err != nil {
		return nil, err
	}
	return floatingips.ExtractFloatingIPs(page)
}
