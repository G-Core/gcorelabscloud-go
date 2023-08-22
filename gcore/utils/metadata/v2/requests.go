package metadataV2

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

// MetadataCreateOrUpdate creates or update a metadata for a resource.
func MetadataCreateOrUpdate(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r tasks.Result) {
	_, r.Err = client.Post(metadata.MetadataURL(client, id), opts, &r.Body, nil)
	return
}

// MetadataReplace replace a metadata for a resource.
func MetadataReplace(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r tasks.Result) {
	_, r.Err = client.Put(metadata.MetadataURL(client, id), opts, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for a resource.
func MetadataDelete(client *gcorecloud.ServiceClient, id string, key string) (r tasks.Result) {
	_, r.Err = client.DeleteWithResponse(metadata.MetadataItemURL(client, id, key), &r.Body, nil)
	return
}
