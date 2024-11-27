package instances

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceActionURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "action")
}

func metadataItemURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "metadata_item")
}
