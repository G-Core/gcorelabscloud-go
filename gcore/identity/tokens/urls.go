package tokens

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func tokenURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("auth", "jwt", "login")
}
func refreshURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("auth", "jwt", "refresh")
}
func refreshGCloudURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("v1", "token", "refresh")
}
func selectAccountURL(c *gcorecloud.ServiceClient, clientID string) string {
	return c.ServiceURL("auth", "jwt", "clients", clientID, "login")
}
