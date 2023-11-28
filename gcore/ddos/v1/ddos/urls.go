package ddos

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("")
}

func getAccessStatusURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("ddos", "accessibility", strconv.Itoa(c.RegionID))
}

func checkRegionCoverageURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("ddos", "region_coverage", strconv.Itoa(c.RegionID))
}

func getProfileTemplatesURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("ddos", "profile-templates", strconv.Itoa(c.RegionID))
}

func listProfilesURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createProfileURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateProfileURL(c *gcorecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id))
}

func deleteProfileURL(c *gcorecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id))
}

func activateProfileURL(c *gcorecloud.ServiceClient, id int) string {
	return c.ServiceURL(strconv.Itoa(id), "action")
}
