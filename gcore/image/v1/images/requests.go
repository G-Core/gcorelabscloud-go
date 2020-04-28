package images

import (
	"net/http"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images/types"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Private    bool             `q:"private"`
	Visibility types.Visibility `q:"visibility"`
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToImageCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a image.
type CreateOpts struct {
	URL        string             `json:"url" required:"true" validate:"required,url"`
	Name       string             `json:"name" required:"true" validate:"required"`
	CowFormat  bool               `json:"cow_format,omitempty"`
	Properties *map[string]string `json:"properties,omitempty"`
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.Validate.Struct(opts)
}

// ToImageCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

func List(client *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific image based on its unique ID.
func Get(client *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	var resp *http.Response
	resp, r.Err = client.Get(url, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Image, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractImages(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// Create an image.
func Create(client *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(createURL(client), b, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

// Delete an image.
func Delete(client *gcorecloud.ServiceClient, imageID string) (r tasks.Result) {
	url := deleteURL(client, imageID)
	var resp *http.Response
	resp, r.Err = client.DeleteWithResponse(url, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}
