package volumes

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// List retrieves list of volumes for GPU virtual cluster
func List(c *gcorecloud.ServiceClient, clusterID string) pagination.Pager {
	url := listURL(c, c.ProjectID, c.RegionID, clusterID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return commonResult{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}
