package users

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("users")
}
