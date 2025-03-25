package flavors

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	flavorsPath = "flavors"
)

// FlavorsURL returns URL for GPU flavors operations
func FlavorsURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(flavorsPath)
}
