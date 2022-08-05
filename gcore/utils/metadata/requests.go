package metadata

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
	"net/http"
)

func MetadataList(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := MetadataURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return MetadataPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func MetadataListAll(client *gcorecloud.ServiceClient, id string) ([]Metadata, error) {
	pages, err := MetadataList(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := ExtractMetadata(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// MetadataCreateOrUpdate creates or update a metadata for a resource.
func MetadataCreateOrUpdate(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r MetadataActionResult) {
	_, r.Err = client.Post(MetadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataReplace replace a metadata for a resource.
func MetadataReplace(client *gcorecloud.ServiceClient, id string, opts map[string]string) (r MetadataActionResult) {
	_, r.Err = client.Put(MetadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for a resource.
func MetadataDelete(client *gcorecloud.ServiceClient, id string, key string) (r MetadataActionResult) {
	_, r.Err = client.Delete(MetadataItemURL(client, id, key), &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataGet gets defined metadata key for a security group.
func MetadataGet(client *gcorecloud.ServiceClient, id string, key string) (r MetadataResult) {
	url := MetadataItemURL(client, id, key)

	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}
