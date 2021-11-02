package apitokens

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, clientID, tokenID int) string {
	return c.ServiceURL("clients", strconv.Itoa(clientID), "tokens", strconv.Itoa(tokenID))
}

func rootURL(c *gcorecloud.ServiceClient, clientID int) string {
	return c.ServiceURL("clients", strconv.Itoa(clientID), "tokens")
}

func getURL(c *gcorecloud.ServiceClient, clientID, tokenID int) string {
	return resourceURL(c, clientID, tokenID)
}

func listURL(c *gcorecloud.ServiceClient, clientID int) string {
	return rootURL(c, clientID)
}

func createURL(c *gcorecloud.ServiceClient, clientID int) string {
	return rootURL(c, clientID)
}

func deleteURL(c *gcorecloud.ServiceClient, clientID, tokenID int) string {
	return resourceURL(c, clientID, tokenID)
}
