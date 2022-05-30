package laas

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func statusURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("status")
}

func usersURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("users")
}

func topicsURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("topics")
}

func deleteTopicURL(c *gcorecloud.ServiceClient, name string) string {
	return c.ServiceURL("topics", name)
}

func kafkaURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("laas", strconv.Itoa(c.RegionID), "kafka_hosts")
}

func openSearchURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("laas", strconv.Itoa(c.RegionID), "opensearch_hosts")
}
