package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// GetFlavor get inference flavor instance.
func GetFlavor(c *gcorecloud.ServiceClient, name string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, name), &r.Body, nil)
	return
}

// ListAllFlavor lists all inference flavors.
func ListAllFlavor(c *gcorecloud.ServiceClient) ([]Flavor, error) {
	var r ListResult
	_, r.Err = c.Get(listURL(c), &r.Body, nil)
	return r.Extract()
}
