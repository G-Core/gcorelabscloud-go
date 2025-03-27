package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
