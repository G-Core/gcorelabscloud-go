package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// DeleteServerOptsBuilder allows extensions to add parameters to delete server options.
type DeleteServerOptsBuilder interface {
	ToServerDeleteQuery() (string, error)
}

// DeleteServerOpts specifies the parameters for the Delete method.
type DeleteServerOpts struct {
	AllFloatingIPs      bool     `q:"all_floating_ips" validate:"omitempty,allowed_without=FloatingIPIDs"`
	AllReservedFixedIPs bool     `q:"all_reserved_fixed_ips" validate:"omitempty,allowed_without=ReservedFixedIPIDs"`
	AllVolumes          bool     `q:"all_volumes" validate:"omitempty,allowed_without=VolumeIDs"`
	FloatingIPIDs       []string `q:"floating_ip_ids" validate:"omitempty,allowed_without=AllFloatingIPs,dive,uuid4" delimiter:"comma"`
	ReservedFixedIPIDs  []string `q:"reserved_fixed_ip_ids" validate:"omitempty,allowed_without=AllReservedFixedIPs,dive,uuid4" delimiter:"comma"`
	VolumeIDs           []string `q:"volume_ids" validate:"omitempty,allowed_without=AllVolumes,dive,uuid4" delimiter:"comma"`
}

// Validate checks if the provided options are valid.
func (opts DeleteServerOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToDeleteClusterActionMap builds a request body from DeleteInstanceOpts.
func (opts DeleteServerOpts) ToServerDeleteQuery() (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List retrieves servers of a specific GPU cluster.
func List(client *gcorecloud.ServiceClient, clusterID string) pagination.Pager {
	url := ClusterServersURL(client, clusterID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all servers of a specific GPU cluster.
func ListAll(client *gcorecloud.ServiceClient, clusterID string) ([]Server, error) {
	pages, err := List(client, clusterID).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractServers(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Delete removes a specific server from a GPU cluster by its ID.
func Delete(client *gcorecloud.ServiceClient, clusterID, serverID string, opts DeleteServerOptsBuilder) (r tasks.Result) {
	url := ClusterServerURL(client, clusterID, serverID)
	if opts != nil {
		query, err := opts.ToServerDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}
