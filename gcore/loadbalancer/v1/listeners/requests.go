package listeners

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToListenerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ListenerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific listener based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToListenerCreateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToListenerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	LoadBalancerID *string `q:"loadbalancer_id"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToListenerListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOpts represents options used to create a listener pool.
type CreateOpts struct {
	Name                 string             `json:"name" required:"true" validate:"required,name"`
	Protocol             types.ProtocolType `json:"protocol" required:"true"`
	ProtocolPort         int                `json:"protocol_port" required:"true"`
	LoadBalancerID       string             `json:"loadbalancer_id" required:"true"`
	InsertXForwarded     bool               `json:"insert_x_forwarded"`
	SecretID             string             `json:"secret_id,omitempty"`
	SNISecretID          []string           `json:"sni_secret_id,omitempty"`
	AllowedCIDRS         []string           `json:"allowed_cidrs,omitempty" validate:"omitempty,dive,cidr"`
	TimeoutClientData    *int               `json:"timeout_client_data,omitempty"`
	TimeoutMemberData    *int               `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect *int               `json:"timeout_member_connect,omitempty"`
	ConnectionLimit		 *int				`json:"connection_limit,omitempty"`
}

// ToListenerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToListenerCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new listener using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToListenerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToListenerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a listener.
type UpdateOpts struct {
	Name                 string   `json:"name,omitempty"`
	SecretID             string   `json:"secret_id,omitempty"`
	SNISecretID          []string `json:"sni_secret_id,omitempty"`
	AllowedCIDRS         []string `json:"allowed_cidrs,omitempty" validate:"omitempty,dive,cidr"`
	TimeoutClientData    *int     `json:"timeout_client_data,omitempty"`
	TimeoutMemberData    *int     `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect *int     `json:"timeout_member_connect,omitempty"`
	ConnectionLimit		 *int	  `json:"connection_limit,omitempty"`
}

// ToListenerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToListenerUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing listener using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, listenerID string, opts UpdateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToListenerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, listenerID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the listener associated with it.
func Delete(c *gcorecloud.ServiceClient, listenerID string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, listenerID), &r.Body, nil)
	return
}

// ListAll returns all LBs
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Listener, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractListeners(page)
}
