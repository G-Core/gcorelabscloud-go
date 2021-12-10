package instances

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

func renameInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func resourceActionDetailsURL(c *gcorecloud.ServiceClient, id string, action string, actionID string) string {
	return c.ServiceURL(id, action, actionID)
}

func interfacesListURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "interfaces")
}

func securityGroupsListURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "securitygroups")
}

func addSecurityGroupsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "addsecuritygroup")
}

func deleteSecurityGroupsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "delsecuritygroup")
}

func attachInterfaceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "attach_interface")
}

func detachInterfaceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "detach_interface")
}

func startInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "start")
}

func stopInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "stop")
}

func powerCycleInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "powercycle")
}

func rebootInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "reboot")
}

func suspendInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "suspend")
}

func resumeInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resume")
}

func changeFlavorInstanceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "changeflavor")
}

func metadataURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metadata")
}

func metadataDetailsURL(c *gcorecloud.ServiceClient, id string, actionID string) string {
	return resourceActionDetailsURL(c, id, "metadata", actionID)
}

func listAvailableFlavorsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "available_flavors")
}

func listInstanceMetricsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metrics")
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getSpiceConsoleURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "get_spice_console")
}

func getInstanceConsoleURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "get_console")
}

func listInstanceLocationURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("instances", "search")
}
