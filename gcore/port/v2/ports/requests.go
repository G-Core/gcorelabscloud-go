package ports

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	ports1 "github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

// AllowAddressPairs assign allowed address pairs for instance port (v2 passes task to process for worker)
func AllowAddressPairs(c *gcorecloud.ServiceClient, portID string, opts ports1.AllowAddressPairsOptsBuilder) (r tasks.Result) {
	b, err := opts.ToAllowAddressPairsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(assignAllowedAddressPairsURL(c, portID), b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}
