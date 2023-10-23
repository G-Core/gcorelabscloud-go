package faas

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateNamespaceOptsBuilder allows extensions to add additional parameters to the request.
type CreateNamespaceOptsBuilder interface {
	ToNamespaceCreateMap() (map[string]interface{}, error)
}

// CreateNamespaceOpts represents options used to create a namespace.
type CreateNamespaceOpts struct {
	Name        string            `json:"name" required:"true" validate:"required"`
	Description string            `json:"description"`
	Envs        map[string]string `json:"envs"`
}

// ToNamespaceCreateMap builds a request body from CreateNamespaceOpts.
func (opts CreateNamespaceOpts) ToNamespaceCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateNamespace create FaaS namespace.
func CreateNamespace(c *gcorecloud.ServiceClient, opts CreateNamespaceOptsBuilder) (r tasks.Result) {
	url := namespaceCreateURL(c)
	b, err := opts.ToNamespaceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// DeleteNamespace delete FaaS namespace.
func DeleteNamespace(c *gcorecloud.ServiceClient, name string) (r tasks.Result) {
	url := namespaceURL(c, name)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// Versions request.
type ListOptsBuilder interface {
	ToFaaSListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Offset  int    `q:"offset"`
	Search  string `q:"search"`
	OrderBy string `q:"order_by"`
}

// ToFaaSListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFaaSListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListNamespace returns a Pager which allows you to iterate over a collection of
// namespaces. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListNamespace(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := namespaceListURL(c)
	if opts != nil {
		query, err := opts.ToFaaSListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NamespacePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListNamespaceALL returns all namespaces.
func ListNamespaceALL(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Namespace, error) {
	page, err := ListNamespace(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractNamespaces(page)
}

// GetNamespace retrieves a specific namespace based on its name.
func GetNamespace(c *gcorecloud.ServiceClient, name string) (r NamespaceResult) {
	url := namespaceURL(c, name)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// UpdateNamespaceOptsBuilder allows extensions to add additional parameters to the request.
type UpdateNamespaceOptsBuilder interface {
	ToNamespaceUpdateMap() (map[string]interface{}, error)
}

// UpdateNamespaceOpts represents options used to update a namespace.
type UpdateNamespaceOpts struct {
	Description string            `json:"description,omitempty"`
	Envs        map[string]string `json:"envs,omitempty"`
}

// ToNamespaceUpdateMap builds a request body from UpdateNamespaceOpts.
func (opts UpdateNamespaceOpts) ToNamespaceUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateNamespace update FaaS namespace.
func UpdateNamespace(c *gcorecloud.ServiceClient, name string, opts UpdateNamespaceOptsBuilder) (r tasks.Result) {
	url := namespaceURL(c, name)
	b, err := opts.ToNamespaceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

// ListFunctions returns a Pager which allows you to iterate over a collection of
// functions. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListFunctions(c *gcorecloud.ServiceClient, nsName string, opts ListOptsBuilder) pagination.Pager {
	url := functionListURL(c, nsName)
	if opts != nil {
		query, err := opts.ToFaaSListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FunctionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListFunctionsALL returns all functions.
func ListFunctionsALL(c *gcorecloud.ServiceClient, nsName string, opts ListOptsBuilder) ([]Function, error) {
	page, err := ListFunctions(c, nsName, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractFunctions(page)
}

// CreateFunctionOptsBuilder allows extensions to add additional parameters to the request.
type CreateFunctionOptsBuilder interface {
	ToFunctionCreateMap() (map[string]interface{}, error)
}

// CreateFunctionOpts represents options used to create a function.
type CreateFunctionOpts struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Envs         map[string]string   `json:"envs"`
	Runtime      string              `json:"runtime"`
	Timeout      int                 `json:"timeout"`
	Flavor       string              `json:"flavor"`
	Autoscaling  FunctionAutoscaling `json:"autoscaling"`
	CodeText     string              `json:"code_text"`
	EnableApiKey *bool               `json:"enable_api_key,omitempty"`
	Keys         []string            `json:"keys,omitempty"`
	Disabled     *bool               `json:"disabled,omitempty"`
	MainMethod   string              `json:"main_method"`
	Dependencies string              `json:"dependencies,omitempty"`
}

// ToFunctionCreateMap builds a request body from CreateFunctionOpts.
func (opts CreateFunctionOpts) ToFunctionCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateFunction create FaaS function.
func CreateFunction(c *gcorecloud.ServiceClient, nsName string, opts CreateFunctionOptsBuilder) (r tasks.Result) {
	url := functionCreateURL(c, nsName)
	b, err := opts.ToFunctionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// DeleteFunction delete FaaS function.
func DeleteFunction(c *gcorecloud.ServiceClient, nsName, fName string) (r tasks.Result) {
	url := functionURL(c, nsName, fName)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, &gcorecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}

// GetFunction get FaaS function.
func GetFunction(c *gcorecloud.ServiceClient, nsName, fName string) (r FunctionResult) {
	url := functionURL(c, nsName, fName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// UpdateFunctionOptsBuilder allows extensions to add additional parameters to the request.
type UpdateFunctionOptsBuilder interface {
	ToFunctionUpdateMap() (map[string]interface{}, error)
}

// UpdateFunctionOpts represents options used to Update a function.
type UpdateFunctionOpts struct {
	Description  string               `json:"description,omitempty"`
	Envs         map[string]string    `json:"envs,omitempty"`
	Timeout      int                  `json:"timeout,omitempty"`
	Flavor       string               `json:"flavor,omitempty"`
	Autoscaling  *FunctionAutoscaling `json:"autoscaling,omitempty"`
	CodeText     string               `json:"code_text,omitempty"`
	EnableApiKey *bool                `json:"enable_api_key,omitempty"`
	Keys         *[]string            `json:"keys,omitempty"`
	Disabled     *bool                `json:"disabled,omitempty"`
	Dependencies string               `json:"dependencies,omitempty"`
	MainMethod   string               `json:"main_method,omitempty"`
}

// ToFunctionUpdateMap builds a request body from UpdateFunctionOpts.
func (opts UpdateFunctionOpts) ToFunctionUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateFunction update FaaS function.
func UpdateFunction(c *gcorecloud.ServiceClient, nsName, fName string, opts UpdateFunctionOptsBuilder) (r tasks.Result) {
	url := functionURL(c, nsName, fName)
	b, err := opts.ToFunctionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(url, b, &r.Body, nil)
	return
}

// ListKeys returns a Pager which allows you to iterate over a collection of
// keys. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListKeys(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := keysListURL(c)
	if opts != nil {
		query, err := opts.ToFaaSListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return KeyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListKeysAll returns all keys.
func ListKeysAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Key, error) {
	page, err := ListKeys(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractKeys(page)
}

// CreateKeyOptsBuilder allows extensions to add additional parameters to the request.
type CreateKeyOptsBuilder interface {
	ToKeyCreateMap() (map[string]any, error)
}

// CreateKeyOpts represents options used to create a key.
type CreateKeyOpts struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Expire      *gcorecloud.JSONRFC3339ZZ `json:"expire,omitempty"`
	Functions   []KeysFunction            `json:"functions,omitempty"`
}

func (opts CreateKeyOpts) ToKeyCreateMap() (map[string]any, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateKey create FaaS key.
func CreateKey(c *gcorecloud.ServiceClient, opts CreateKeyOptsBuilder) (r Key, err error) {
	url := keysCreateURL(c)
	b, err := opts.ToKeyCreateMap()
	if err != nil {
		return Key{}, err
	}

	_, err = c.Post(url, b, &r, nil)

	return
}

// DeleteKey delete FaaS key.
func DeleteKey(c *gcorecloud.ServiceClient, kName string) error {
	url := keyURL(c, kName)
	_, err := c.Delete(url, nil)

	return err
}

// GetKey get FaaS key.
func GetKey(c *gcorecloud.ServiceClient, kName string) (r KeyResult) {
	url := keyURL(c, kName)
	_, r.Err = c.Get(url, &r.Body, nil)

	return
}

// UpdateKeyOptsBuilder allows extensions to add additional parameters to the request.
type UpdateKeyOptsBuilder interface {
	ToKeyUpdateMap() (map[string]any, error)
}

// UpdateKeyOpts represents options used to Update a key.
type UpdateKeyOpts struct {
	Description string                    `json:"description,omitempty"`
	Expire      *gcorecloud.JSONRFC3339ZZ `json:"expire,omitempty"`
	Functions   []KeysFunction            `json:"functions,omitempty"`
}

// ToKeyUpdateMap builds a request body from UpdateKeyOpts.
func (opts UpdateKeyOpts) ToKeyUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateKey update FaaS key.
func UpdateKey(c *gcorecloud.ServiceClient, kName string, opts UpdateKeyOptsBuilder) (k Key, err error) {
	url := keyURL(c, kName)
	b, err := opts.ToKeyUpdateMap()
	if err != nil {
		return Key{}, err
	}

	_, err = c.Patch(url, b, &k, nil)

	return
}
