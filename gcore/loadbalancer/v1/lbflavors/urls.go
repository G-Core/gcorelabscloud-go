package lbflavors

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func listURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}
