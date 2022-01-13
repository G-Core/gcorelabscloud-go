package limits

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("limits_request", strconv.Itoa(id))
}

func getURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
