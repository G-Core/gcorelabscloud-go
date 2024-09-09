package ports

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
)

// AllowAddressPairsOptsBuilder allows extensions to add additional parameters to the AllowAddressPairs request.
type AllowAddressPairsOptsBuilder interface {
	ToAllowAddressPairsMap() (map[string]interface{}, error)
}

// AllowAddressPairsOpts represents options used to allow address pairs.
type AllowAddressPairsOpts struct {
	AllowedAddressPairs []reservedfixedips.AllowedAddressPairs `json:"allowed_address_pairs,omitempty"`
}

// ToAllowAddressPairsMap builds a request body from AllowAddressPairsOpts.
func (opts AllowAddressPairsOpts) ToAllowAddressPairsMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts AllowAddressPairsOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

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

// AllowAddressPairs assign allowed address pairs for instance port
func AllowAddressPairs(c *gcorecloud.ServiceClient, portID string, opts AllowAddressPairsOptsBuilder) (r AssignResult) {
	b, err := opts.ToAllowAddressPairsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(assignAllowedAddressPairsURL(c, portID), b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}
