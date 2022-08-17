package bmcapacity

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

// GetAvailableNodes retrieves available baremetal nodes
func GetAvailableNodes(c *gcorecloud.ServiceClient) (r GetAvailableNodesResult) {
	url := getAvailableNodesURL(c)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
