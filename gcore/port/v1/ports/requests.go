package ports

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// EnablePortSecurity
func EnablePortSecurity(c *gcorecloud.ServiceClient, portID string) (r UpdateResult) {
	_, r.Err = c.Post(enablePortSecurityURL(c, portID), nil, &r.Body, nil)
	return
}

// DisablePortSecurity
func DisablePortSecurity(c *gcorecloud.ServiceClient, portID string) (r UpdateResult) {
	_, r.Err = c.Post(disablePortSecurityURL(c, portID), nil, &r.Body, nil)
	return
}
