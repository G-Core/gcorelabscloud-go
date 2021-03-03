package ports

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func enablePortSecurityURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "enable_port_security")
}

func disablePortSecurityURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "disable_port_security")
}
