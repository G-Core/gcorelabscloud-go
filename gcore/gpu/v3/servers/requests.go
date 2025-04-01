package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// List retrieves servers of a specific GPU cluster.
func List(client *gcorecloud.ServiceClient, clusterID string) (r ListResult) {
	url := ClusterServersURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
