package securitygroups

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Limit      int               `q:"limit" validate:"omitempty,gt=0"`
	Offset     int               `q:"offset" validate:"omitempty,gte=0"`
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToSequirityGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSecurityGroupListQuery() (string, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToSecurityGroupListQuery() (string, error)
}

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)

	if opts != nil {
		query, err := opts.ToSecurityGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SecurityGroupPage{pagination.OffsetPageBase{PageResult: r}}
	})
}

// Get retrieves a specific security group based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the request.
type CreateOptsBuilder interface {
	ToSecurityGroupCreateMap() (map[string]interface{}, error)
}

// CreateRuleOptsBuilder allows extensions to add additional parameters to the request.
type CreateRuleOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
}

// CreateSecurityGroupRuleOpts represents options used to create a security group rule.
type CreateSecurityGroupRuleOpts struct {
	Direction       types.RuleDirection `json:"direction" required:"true"`
	EtherType       types.EtherType     `json:"ethertype,omitempty" required:"true"`
	Protocol        types.Protocol      `json:"protocol,omitempty" required:"true"`
	SecurityGroupID *string             `json:"security_group_id,omitempty"`
	RemoteGroupID   *string             `json:"remote_group_id,omitempty"`
	PortRangeMax    *int                `json:"port_range_max,omitempty"`
	PortRangeMin    *int                `json:"port_range_min,omitempty"`
	Description     *string             `json:"description,omitempty"`
	RemoteIPPrefix  *string             `json:"remote_ip_prefix,omitempty"`
}

// ToRuleCreateMap builds a request body from CreateSecurityGroupRuleOpts.
func (opts CreateSecurityGroupRuleOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

type CreateSecurityGroupOpts struct {
	Name               string                        `json:"name" required:"true"`
	Description        *string                       `json:"description,omitempty"`
	SecurityGroupRules []CreateSecurityGroupRuleOpts `json:"security_group_rules"`
	Metadata           map[string]interface{}        `json:"metadata,omitempty"`
}

// CreateOpts represents options used to create a security group.
type CreateOpts struct {
	SecurityGroup CreateSecurityGroupOpts `json:"security_group" required:"true"`
	Instances     []string                `json:"instances,omitempty"`
}

// ToSecurityGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSecurityGroupCreateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new security group using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecurityGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToSecurityGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a security group.
type UpdateOpts struct {
	Name         string                        `json:"name,omitempty"`
	ChangedRules []UpdateSecurityGroupRuleOpts `json:"changed_rules,omitempty"`
}

// UpdateSecurityGroupRuleOpts represents options used to change a security group rule.
type UpdateSecurityGroupRuleOpts struct {
	Action              types.Action        `json:"action" required:"true"`
	SecurityGroupRuleID string              `json:"security_group_rule_id,omitempty"`
	Direction           types.RuleDirection `json:"direction,omitempty"`
	EtherType           types.EtherType     `json:"ethertype,omitempty"`
	Protocol            types.Protocol      `json:"protocol,omitempty"`
	SecurityGroupID     *string             `json:"security_group_id,omitempty"`
	RemoteGroupID       *string             `json:"remote_group_id,omitempty"`
	PortRangeMax        *int                `json:"port_range_max,omitempty"`
	PortRangeMin        *int                `json:"port_range_min,omitempty"`
	Description         *string             `json:"description,omitempty"`
	RemoteIPPrefix      *string             `json:"remote_ip_prefix,omitempty"`
}

// ToSecurityGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSecurityGroupUpdateMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing security group using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, securityGroupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSecurityGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, securityGroupID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the security group associated with it.
func Delete(c *gcorecloud.ServiceClient, securityGroupID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, securityGroupID), nil)
	return
}

// ListAll returns all SGs
func ListAll(c *gcorecloud.ServiceClient, opts ListOptsBuilder) ([]SecurityGroup, error) {
	page, err := List(c, opts).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroups(page)
}

// AddRule accepts a CreateSecurityGroupRuleOpts struct and add rule to existed group.
func AddRule(c *gcorecloud.ServiceClient, securityGroupID string, opts CreateRuleOptsBuilder) (r CreateRuleResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(addRulesURL(c, securityGroupID), b, &r.Body, nil)
	return
}

// ListInstances returns page of instances for SG
func ListInstances(c *gcorecloud.ServiceClient, securityGroupID string) pagination.Pager {
	url := listInstancesURL(c, securityGroupID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SecurityGroupInstancesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllInstances returns all instances for SG
func ListAllInstances(c *gcorecloud.ServiceClient, securityGroupID string) ([]instances.Instance, error) {
	page, err := ListInstances(c, securityGroupID).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroupInstances(page)
}

// IDFromName is a convenience function that returns a security group's ID,
// given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	sgs, err := ListAll(client, nil)
	if err != nil {
		return "", err
	}

	for _, s := range sgs {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "security group"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "security group"}
	}
}

// DeepCopyOptsBuilder allows extensions to add additional parameters to the request.
type DeepCopyOptsBuilder interface {
	ToDeepCopyMap() (map[string]interface{}, error)
}

// DeepCopyOpts represents options used to deep copy a security group.
type DeepCopyOpts struct {
	Name string `json:"name" required:"true"`
}

// ToDeepCopyMap builds a request body from DeepCopyOpts.
func (opts DeepCopyOpts) ToDeepCopyMap() (map[string]interface{}, error) {
	if err := gcorecloud.ValidateStruct(opts); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

// DeepCopy accepts a DeepCopyOpts struct and create a deep copy of security group.
func DeepCopy(c *gcorecloud.ServiceClient, securityGroupID string, opts DeepCopyOptsBuilder) (r DeepCopyResult) {
	b, err := opts.ToDeepCopyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(deepCopyURL(c, securityGroupID), b, nil, nil)
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
