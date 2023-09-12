package aiflavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func listAIFlavorsURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}