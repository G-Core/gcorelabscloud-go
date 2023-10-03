package ai

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
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

func resourceAIInstanceActionURL(c *gcorecloud.ServiceClient, instance_id string, action string) string {
	return c.ServiceURL(instance_id, action)
}

func interfacesListURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "interfaces")
}

func portsListURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "ports")
}

func attachAIInstanceInterfaceURL(c *gcorecloud.ServiceClient, instance_id string) string {
	return resourceAIInstanceActionURL(c, instance_id, "attach_interface")
}

func detachAIInstanceInterfaceURL(c *gcorecloud.ServiceClient, instance_id string) string {
	return resourceAIInstanceActionURL(c, instance_id, "detach_interface")
}

func addSecurityGroupsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "addsecuritygroup")
}

func deleteSecurityGroupsURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "delsecuritygroup")
}

func powerCycleAIURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "powercycle")
}

func powerCycleAIInstanceURL(c *gcorecloud.ServiceClient, instance_id string) string {
	return resourceAIInstanceActionURL(c, instance_id, "powercycle")
}

func rebootAIURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "reboot")
}

func rebootAIInstanceURL(c *gcorecloud.ServiceClient, instance_id string) string {
	return resourceAIInstanceActionURL(c, instance_id, "reboot")
}

func getAIInstanceConsoleURL(c *gcorecloud.ServiceClient, instance_id string) string {
	return resourceAIInstanceActionURL(c, instance_id, "get_console")
}

func suspendAIURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "suspend")
}

func resumeAIURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resume")
}

func resizeAIURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resize")
}

func metadataURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metadata")
}

func metadataItemURL(c *gcorecloud.ServiceClient, id string, key string) string {
	return resourceActionURL(c, id, fmt.Sprintf("metadata_item?key=%s", key))
}
