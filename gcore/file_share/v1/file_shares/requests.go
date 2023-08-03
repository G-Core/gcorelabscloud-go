package file_shares

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type ListOptsBuilder interface {
	ToFileShareListQuery() (string, error)
}

// List returns a Pager which allows you to iterate over a collection of
// file shares. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToFileShareListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FileSharePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific file share based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToFileShareCreateMap() (map[string]interface{}, error)
}

type FileShareNetworkOpts struct {
	NetworkID string `json:"network_id,omitempty" validate:"uuid4"`
	SubnetID  string `json:"subnet_id,omitempty" validate:"uuid4"`
}

type CreateAccessRuleOpts struct {
	IPAddress  string `json:"ip_address,omitempty" validate:"required,ipv4|cidrv4,omitempty"`
	AccessMode string `json:"access_mode" validate:"required,oneof=ro rw"`
}

// CreateOpts represents options used to create a file share.
type CreateOpts struct {
	Name     string                 `json:"name" required:"true" validate:"required"`
	Protocol string                 `json:"protocol" required:"true" validate:"required,oneof=NFS"`
	Size     int                    `json:"size" required:"true" validate:"required,gt=1"`
	Network  FileShareNetworkOpts   `json:"network" required:"true" validate:"required,dive"`
	Access   []CreateAccessRuleOpts `json:"access,omitempty" validate:"dive"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// ToFileShareCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFileShareCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// ToFileShareListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFileShareListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Validate
func (opts CreateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Create accepts a CreateOpts struct and creates a new file share using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToFileShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToFileShareUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a file share.
type UpdateOpts struct {
	Name string `json:"name" required:"true" validate:"required"`
}

// ResizeOptsBuilder has parameters for resize request.
type ResizeOptsBuilder interface {
	ToFileShareResizeMap() (map[string]interface{}, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct{}

// ToFileShareUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToFileShareUpdateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts UpdateOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Update accepts a UpdateOpts struct and updates an existing file share using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, fileShareID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFileShareUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, fileShareID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the file share associated with it.
func Delete(c *gcorecloud.ServiceClient, fileShareID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, fileShareID), &r.Body, nil)
	return
}

// ExtendOpts represents options used to resize a file share.
type ExtendOpts struct {
	Size int `json:"size" required:"true" validate:"required,gt=1"`
}

// ToFileShareResizeMap builds a request body from ResizeOpts.
func (opts ExtendOpts) ToFileShareResizeMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Validate
func (opts ExtendOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// Extend accepts a ExtendOpts struct and resize an existing file share using the
// values provided.
func Extend(c *gcorecloud.ServiceClient, fileShareID string, opts ResizeOptsBuilder) (r tasks.Result) {
	b, err := opts.ToFileShareResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(extendResourceUrl(c, fileShareID), b, &r.Body, nil)
	return
}

// ListAll is a convenience function that returns all file shares.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]FileShare, error) {
	pages, err := List(client, opts).AllPages()

	if err != nil {
		return nil, err
	}

	all, err := ExtractFileShares(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

type AccessRuleListOptsBuilder interface {
	ToAccessRuleListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type AccessRuleListOpts struct{}

// ToAccessRuleListQuery formats a ListOpts into a query string.
func (opts AccessRuleListOpts) ToAccessRuleListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// file shares. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListAccessRules(c *gcorecloud.ServiceClient, fileShareID string, opts AccessRuleListOptsBuilder) pagination.Pager {
	url := accessRuleURL(c, fileShareID)
	if opts != nil {
		query, err := opts.ToAccessRuleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AccessRulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateAccessRuleOptsBuilder allows extensions to add additional parameters to the
// AccessRuleCreate request.
type CreateAccessRuleOptsBuilder interface {
	ToAccessRuleCreateMap() (map[string]interface{}, error)
}

// Validate
func (opts CreateAccessRuleOpts) Validate() error {
	return gcorecloud.TranslateValidationError(gcorecloud.Validate.Struct(opts))
}

// ToAccessRuleCreateMap builds a request body from CreateAccessRuleOpts.
func (opts CreateAccessRuleOpts) ToAccessRuleCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateAccessRule accepts a CreateAccessRuleOpts struct and creates a new file share access rule using the values provided.
func CreateAccessRule(c *gcorecloud.ServiceClient, fileShareID string, opts CreateAccessRuleOptsBuilder) (r CreateAccessRuleResult) {
	b, err := opts.ToAccessRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(accessRuleURL(c, fileShareID), b, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}

// Delete accepts a unique ID and deletes the file share access rule associated with it.
func DeleteAccessRule(c *gcorecloud.ServiceClient, fileShareID string, ruleID string) (r DeleteResult) {
	_, r.Err = c.Delete(accessRuleItemURL(c, fileShareID, ruleID), &gcorecloud.RequestOpts{OkCodes: []int{http.StatusNoContent}})
	return
}

func MetadataList(client *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := metadataURL(client, id)
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

// MetadataCreateOrUpdate creates or update a metadata for an security group.
func MetadataCreateOrUpdate(client *gcorecloud.ServiceClient, id string, opts map[string]interface{}) (r MetadataActionResult) {
	_, r.Err = client.Post(metadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataReplace replace a metadata for an security group.
func MetadataReplace(client *gcorecloud.ServiceClient, id string, opts map[string]interface{}) (r MetadataActionResult) {
	_, r.Err = client.Put(metadataURL(client, id), opts, nil, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for a security group.
func MetadataDelete(client *gcorecloud.ServiceClient, id string, key string) (r MetadataActionResult) {
	_, r.Err = client.Delete(metadataItemURL(client, id, key), &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataGet gets defined metadata key for a security group.
func MetadataGet(client *gcorecloud.ServiceClient, id string, key string) (r MetadataResult) {
	url := metadataItemURL(client, id, key)

	_, r.Err = client.Get(url, &r.Body, nil) // nolint
	return
}
