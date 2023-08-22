package metadataV1

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func MetadataList(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := metadata.MetadataURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return metadata.MetadataPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func MetadataListAll(client *gcorecloud.ServiceClient, id string) ([]metadata.Metadata, error) {
	pages, err := MetadataList(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := metadata.ExtractMetadata(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// MetadataCreateOrUpdate creates or update a metadata for a resource.
func MetadataCreateOrUpdate(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r metadata.MetadataActionResult) {
	_, r.Err = client.Post(metadata.MetadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataReplace replace a metadata for a resource.
func MetadataReplace(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r metadata.MetadataActionResult) {
	_, r.Err = client.Put(metadata.MetadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for a resource.
func MetadataDelete(client *gcorecloud.ServiceClient, id string, key string) (r metadata.MetadataActionResult) {
	_, r.Err = client.Delete(metadata.MetadataItemURL(client, id, key), &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataGet gets defined metadata key for a resource.
func MetadataGet(client *gcorecloud.ServiceClient, id string, key string) (r metadata.MetadataResult) {
	url := metadata.MetadataItemURL(client, id, key)

	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}
