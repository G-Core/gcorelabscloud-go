package ddos

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// CreateProfileOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateProfileOptsBuilder interface {
	ToProfileCreateMap() (map[string]interface{}, error)
}

type CreateProfileOpts struct {
	ProfileTemplate     int            `json:"profile_template" required:"true" validate:"required"`
	BaremetalInstanceID string         `json:"bm_instance_id" required:"true" validate:"required"`
	IPAddress           string         `json:"ip_address" required:"true" validate:"required,ip4_addr"`
	Fields              []ProfileField `json:"fields"`
}

// ToProfileCreateMap builds a request body from CreateProfileOpts.
func (opts CreateProfileOpts) ToProfileCreateMap() (map[string]interface{}, error) {
	result, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateProfileOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateProfileOptsBuilder interface {
	ToProfileUpdateMap() (map[string]interface{}, error)
}

type UpdateProfileOpts struct {
	ProfileTemplate     int            `json:"profile_template" required:"true" validate:"required"`
	BaremetalInstanceID string         `json:"bm_instance_id" required:"true" validate:"required"`
	IPAddress           string         `json:"ip_address" required:"true" validate:"required,ip4_addr"`
	Fields              []ProfileField `json:"fields"`
}

// ToProfileUpdateMap builds a request body from UpdateProfileOpts.
func (opts UpdateProfileOpts) ToProfileUpdateMap() (map[string]interface{}, error) {
	result, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return result, nil
}

type ActivateProfileOptsBuilder interface {
	ToActivateProfileMap() (map[string]interface{}, error)
}

type ActivateProfileOpts struct {
	BGP    bool `json:"bgp"`
	Active bool `json:"active"`
}

// ToActivateProfileMap builds a request bode from ActivateProfileOptsBuilder.
func (opts ActivateProfileOpts) ToActivateProfileMap() (map[string]interface{}, error) {
	result, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccessibility retrieves DDoS protection service status
func GetAccessibility(c *gcorecloud.ServiceClient) (r GetAccessStatusResult) {
	url := getAccessStatusURL(c)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CheckRegionCoverage retrieves region coverage by the DDoS protection features
func CheckRegionCoverage(c *gcorecloud.ServiceClient) (r CheckRegionCoverageResult) {
	url := checkRegionCoverageURL(c)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

func ListProfileTemplates(c *gcorecloud.ServiceClient) pagination.Pager {
	url := getProfileTemplatesURL(c)

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ProfileTemplatesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListProfiles(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listProfilesURL(c)

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ProfilesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAllProfileTemplates returns all DDoS protection profile templates
func ListAllProfileTemplates(c *gcorecloud.ServiceClient) ([]ProfileTemplate, error) {
	page, err := ListProfileTemplates(c).AllPages()
	if err != nil {
		return nil, err
	}

	return ExtractProfileTemplates(page)
}

// ListAllProfiles returns active clients DDoS protection profiles
func ListAllProfiles(c *gcorecloud.ServiceClient) ([]Profile, error) {
	page, err := ListProfiles(c).AllPages()
	if err != nil {
		return nil, err
	}

	return ExtractProfiles(page)
}

// CreateProfile accepts a CreateProfileOpts struct and creates a new profile using the values provided.
func CreateProfile(c *gcorecloud.ServiceClient, opts CreateProfileOptsBuilder) (r tasks.Result) {
	b, err := opts.ToProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createProfileURL(c), b, &r.Body, nil)
	return
}

// UpdateProfile accepts an UpdateProfileOpts struct and updates a profile with given ID using the values provided.
func UpdateProfile(c *gcorecloud.ServiceClient, id int, opts UpdateProfileOptsBuilder) (r tasks.Result) {
	b, err := opts.ToProfileUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Patch(updateProfileURL(c, id), b, &r.Body, nil)
	return
}

// DeleteProfile accepts a unique ID and deletes the DDoS protection profile associated with it.
func DeleteProfile(c *gcorecloud.ServiceClient, profileID int) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteProfileURL(c, profileID), &r.Body, nil)
	return
}

func ActivateProfile(c *gcorecloud.ServiceClient, id int, opts ActivateProfileOptsBuilder) (r tasks.Result) {
	b, err := opts.ToActivateProfileMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(activateProfileURL(c, id), b, &r.Body, nil)
	return
}
