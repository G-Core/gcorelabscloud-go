package inferences

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

// CreateInferenceDeploymentOptsBuilder allows extensions to add additional parameters to the request.
type CreateInferenceDeploymentOptsBuilder interface {
	ToRegistryCredentialCreateMap() (map[string]interface{}, error)
}

// CreateInferenceDeploymentOpts represents options used to create a function.
type CreateInferenceDeploymentOpts struct {
	Name            string                `json:"name"`
	Image           string                `json:"image"`
	ListeningPort   int                   `json:"listening_port"`
	Description     string                `json:"description"`
	AuthEnabled     bool                  `json:"auth_enabled"`
	Containers      []CreateContainerOpts `json:"containers"`
	Timeout         *int                  `json:"timeout,omitempty"`
	Envs            map[string]string     `json:"envs,omitempty"`
	Command         []string              `json:"command,omitempty"`
	CredentialsName *string               `json:"credentials_name,omitempty"`
	Logging         *CreateLoggingOpts    `json:"logging,omitempty"`
	Probes          *Probes               `json:"probes,omitempty"`
	FlavorName      string                `json:"flavor_name"`
}

// CreateContainerOpts represents options used to create a container.
type CreateContainerOpts struct {
	RegionID int            `json:"region_id"`
	Scale    ContainerScale `json:"scale"`
}

// CreateLoggingOpts represents options used to create a logging.
type CreateLoggingOpts struct {
	Enabled             bool                   `json:"enabled"`
	DestinationRegionID int                    `json:"destination_region_id"`
	TopicName           string                 `json:"topic_name"`
	RetentionPolicy     LoggingRetentionPolicy `json:"retention_policy"`
}

// ToInferenceInstanceCreateMap builds a request body from CreateInferenceDeploymentOpts.
func (opts CreateInferenceDeploymentOpts) ToRegistryCredentialCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// CreateInferenceDeployment create FaaS function.
func CreateInferenceDeployment(c *gcorecloud.ServiceClient, opts CreateInferenceDeploymentOptsBuilder) (r tasks.Result) {
	url := createURL(c)
	b, err := opts.ToRegistryCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(url, b, &r.Body, nil)
	return
}

// GetInferenceDeployment get inference deployment instance.
func GetInferenceDeployment(c *gcorecloud.ServiceClient, name string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, name), &r.Body, nil)
	return
}

// ListAllInferenceDeployments lists all inference deployments.
func ListAllInferenceDeployments(c *gcorecloud.ServiceClient) ([]InferenceDeployment, error) {
	var r ListResult
	_, r.Err = c.Get(listURL(c), &r.Body, nil)
	return r.Extract()
}

// DeleteInferenceDeployment accepts a unique ID and deletes the inference deployment associated with it.
func DeleteInferenceDeployment(c *gcorecloud.ServiceClient, name string) (r tasks.Result) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, name), &r.Body, nil)
	return
}

// UpdateInferenceDeploymentOptsBuilder allows extensions to add additional parameters to the request.
type UpdateInferenceDeploymentOptsBuilder interface {
	ToRegistryCredentialUpdateMap() (map[string]interface{}, error)
}

// UpdateInferenceDeploymentOpts represents options used to update a function.
type UpdateInferenceDeploymentOpts struct {
	Description     *string               `json:"description,omitempty"`
	Image           *string               `json:"image,omitempty"`
	ListeningPort   *int                  `json:"listening_port,omitempty"`
	AuthEnabled     *bool                 `json:"auth_enabled,omitempty"`
	Containers      []CreateContainerOpts `json:"containers,omitempty"`
	Timeout         *int                  `json:"timeout"`
	Envs            map[string]string     `json:"envs,omitempty"`
	Command         []string              `json:"command"`
	Logging         *CreateLoggingOpts    `json:"logging,omitempty"`
	Probes          *Probes               `json:"probes,omitempty"`
	FlavorName      *string               `json:"flavor_name"`
	CredentialsName *string               `json:"credentials_name"`
}

// ToInferenceDeploymentUpdateMap builds a request body from UpdateInferenceDeploymentOpts.
func (opts UpdateInferenceDeploymentOpts) ToRegistryCredentialUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// UpdateInferenceDeployment update existing inference deployment.
func UpdateInferenceDeployment(c *gcorecloud.ServiceClient, name string, opts UpdateInferenceDeploymentOptsBuilder) (r tasks.Result) {
	url := updateURL(c, name)
	b, err := opts.ToRegistryCredentialUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(url, b, &r.Body, nil)
	return
}
