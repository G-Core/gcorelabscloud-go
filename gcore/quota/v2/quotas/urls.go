package quotas

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func getCombinedURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("client_quotas")
}

func getGlobalURL(c *gcorecloud.ServiceClient, clientID int) string {
	return c.BaseServiceURL("global_quotas", strconv.Itoa(clientID))
}

func getRegionURL(c *gcorecloud.ServiceClient, clientID, regionID int) string {
	return c.BaseServiceURL("regional_quotas", strconv.Itoa(clientID), strconv.Itoa(regionID))
}
