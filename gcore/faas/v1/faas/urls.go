package faas

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("")
}

func namespaceListURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func namespaceCreateURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func namespaceURL(c *gcorecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName)
}

func functionListURL(c *gcorecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName, "functions")
}

func functionCreateURL(c *gcorecloud.ServiceClient, namespaceName string) string {
	return c.ServiceURL(namespaceName, "functions")
}

func functionURL(c *gcorecloud.ServiceClient, namespaceName, functionName string) string {
	return c.ServiceURL(namespaceName, "functions", functionName)
}
