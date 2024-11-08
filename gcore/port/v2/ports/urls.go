package ports

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func assignAllowedAddressPairsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "allow_address_pairs")
}
