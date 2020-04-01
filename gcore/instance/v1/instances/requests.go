package instances

import (
	"net/http"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	ExcludeSecGroup   *string `q:"exclude_secgroup"`
	AvailableFloating *string `q:"available_floating"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// SecurityGroupOptsBuilder allows extensions to add parameters to the security groups request.
type SecurityGroupOptsBuilder interface {
	ToSecurityGroupActionMap() (map[string]interface{}, error)
}

type SecurityGroupOpts struct {
	Name string `json:"name" required:"true"`
}

// ToSecurityGroupActionMap builds a request body from SecurityGroupOpts.
func (opts SecurityGroupOpts) ToSecurityGroupActionMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific instance based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	var resp *http.Response
	resp, r.Err = c.Get(url, &r.Body, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

// ListInterfaces retrieves network interfaces for instance
func ListInterfaces(c *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := interfacesListURL(c, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return InstanceInterfacePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all instances.
func ListAll(client *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]Instance, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstances(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListInterfacesAll is a convenience function that returns all instance interfaces.
func ListInterfacesAll(client *gcorecloud.ServiceClient, id string) ([]Interface, error) {
	pages, err := ListInterfaces(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceInterfaces(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// ListSecurityGroups retrieves security groups interfaces for instance
func ListSecurityGroups(c *gcorecloud.ServiceClient, id string) pagination.Pager {
	url := securityGroupsListURL(c, id)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return InstanceSecurityGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListSecurityGroupsAll is a convenience function that returns all instance security groups.
func ListSecurityGroupsAll(client *gcorecloud.ServiceClient, id string) ([]gcorecloud.ItemIDName, error) {
	pages, err := ListSecurityGroups(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceSecurityGroups(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// AssignSecurityGroup adds a security groups to the instance.
func AssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(addSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	defer func() {
		_ = resp.Body.Close()
	}()

	return
}

// UnAssignSecurityGroup removes a security groups from the instance.
func UnAssignSecurityGroup(client *gcorecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	var resp *http.Response
	resp, r.Err = client.Post(deleteSecurityGroupsURL(client, id), b, nil, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	defer func() {
		_ = resp.Body.Close()
	}()
	return
}
