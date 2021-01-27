package reservedfixedips

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func switchVIPURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func connectedDeviceListURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "connected_devices")
}

func availableDeviceListURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "available_devices")
}

func portsToShareVIPURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "connected_devices")
}
