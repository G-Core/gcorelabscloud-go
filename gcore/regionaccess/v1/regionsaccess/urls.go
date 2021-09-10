package regionsaccess

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("reseller_region")
}

func resourceURL(c *gcorecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("reseller_region", strconv.Itoa(id))
}
