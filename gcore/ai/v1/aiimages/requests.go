package aiimages

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)


type VisibilityType string

const (
	PRIVATE VisibilityType = "private"
    PUBLIC VisibilityType = "public"
    SHARED VisibilityType = "shared"
)

type ListOptsBuilder interface {
	ToAIImageListQuery() (string, error)
}

type AIImageListOpts struct {
	Visibility  string `q:"visibility" validate:"omitempty,enum"`
	Private	string `q:"private" validate:"omitempty"`
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts AIImageListOpts) ToAIImageListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}



// List retrieves list of flavors
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listAIImagesURL(c)
	if opts != nil {
		query, err := opts.ToAIImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AIImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll retrieves list of flavors
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]AIImage, error) {
	results, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractAIImages(results)
}