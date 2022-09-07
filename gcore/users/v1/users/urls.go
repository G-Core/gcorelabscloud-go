package users

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func createUserURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("users")
}

func createApiTokenURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("permanent_api_token")
}
