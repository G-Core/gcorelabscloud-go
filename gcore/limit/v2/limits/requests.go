package limits

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"net/http"
)

// Delete deleted limit request
func Delete(c *gcorecloud.ServiceClient, id int) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, id), &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
