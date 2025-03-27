package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// Delete removes a specific GPU cluster by its ID.
func Delete(client *gcorecloud.ServiceClient, clusterID string) (r tasks.Result) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}
