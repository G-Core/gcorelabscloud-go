package inferences

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type InferenceDeployment struct {
	ProjectID       int               `json:"project_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Image           string            `json:"image"`
	ListeningPort   int               `json:"listening_port"`
	CreatedAt       *string           `json:"created_at"`
	AuthEnabled     bool              `json:"auth_enabled"`
	Address         string            `json:"address"`
	Timeout         int               `json:"timeout"`
	Envs            map[string]string `json:"envs"`
	FlavorName      string            `json:"flavor_name"`
	Command         *string           `json:"command"`
	CredentialsName string            `json:"credentials_name"`
	Logging         *Logging          `json:"logging"`
	Probes          *Probes           `json:"probes"`
	Containers      []Container       `json:"containers"`
	Status          string            `json:"status"`
}

type Probes struct {
	LivenessProbe  *ProbeConfiguration `json:"liveness_probe"`
	ReadinessProbe *ProbeConfiguration `json:"readiness_probe"`
	StartupProbe   *ProbeConfiguration `json:"startup_probe"`
}

type ProbeConfiguration struct {
	Enabled bool   `json:"enabled"`
	Probe   *Probe `json:"probe,omitempty"`
}

type Probe struct {
	FailureThreshold    int             `json:"failure_threshold"`
	InitialDelaySeconds int             `json:"initial_delay_seconds"`
	PeriodSeconds       int             `json:"period_seconds"`
	TimeoutSeconds      int             `json:"timeout_seconds"`
	SuccessThreshold    int             `json:"success_threshold"`
	Exec                *ExecProbe      `json:"exec"`
	TcpSocket           *TcpSocketProbe `json:"tcp_socket"`
	HttpGet             *HttpGetProbe   `json:"http_get"`
}

type ExecProbe struct {
	Command []string `json:"command"`
}

type TcpSocketProbe struct {
	Port int `json:"port"`
}

type HttpGetProbe struct {
	Headers map[string]string `json:"headers"`
	Host    *string           `json:"host"`
	Path    string            `json:"path"`
	Port    int               `json:"port"`
	Schema  string            `json:"schema"`
}

type Container struct {
	RegionID     int                   `json:"region_id"`
	Address      string                `json:"address"`
	Scale        ContainerScale        `json:"scale"`
	DeployStatus ContainerDeployStatus `json:"deploy_status"`
	ErrorMessage string                `json:"error_message"`
}

type ContainerDeployStatus struct {
	Total int `json:"total"`
	Ready int `json:"ready"`
}

type ContainerScale struct {
	Min             int                   `json:"min"`
	Max             int                   `json:"max"`
	CooldownPeriod  *int                  `json:"cooldown_period"`
	PollingInterval *int                  `json:"polling_interval"`
	Triggers        ContainerScaleTrigger `json:"triggers"`
}

type ContainerScaleTrigger struct {
	Cpu            *ScaleTriggerThreshold `json:"cpu,omitempty"`
	GpuMemory      *ScaleTriggerThreshold `json:"gpu_memory,omitempty"`
	GpuUtilization *ScaleTriggerThreshold `json:"gpu_utilization,omitempty"`
	Memory         *ScaleTriggerThreshold `json:"memory,omitempty"`
	Http           *ScaleTriggerHttp      `json:"http,omitempty"`
	Sqs            *ScaleTriggerSqs       `json:"sqs,omitempty"`
}

type ScaleTriggerSqs struct {
	QueueURL              string  `json:"queue_url"`
	QueueLength           int     `json:"queue_length"`
	ActivationQueueLength int     `json:"activation_queue_length"`
	ScaleOnFlight         bool    `json:"scale_on_flight"`
	ScaleOnDelayed        bool    `json:"scale_on_delayed"`
	AwsRegion             string  `json:"aws_region"`
	AwsEndpoint           *string `json:"aws_endpoint"`
	SecretName            string  `json:"secret_name"`
}

type ScaleTriggerThreshold struct {
	Threshold int `json:"threshold"`
}

type ScaleTriggerHttp struct {
	Rate   *int `json:"rate"`
	Window *int `json:"window"`
}

type Logging struct {
	Enabled                  bool    `json:"enabled"`
	DestinationRegionID      *int    `json:"destination_region_id"`
	TopicName                *string `json:"topic_name"`
	RetentionPolicy          *int    `json:"retention_policy"`
	OpensearchDashboardsLink string  `json:"opensearch_dashboards_link"`
}

type LoggingRetentionPolicy struct {
	Period *int `json:"period"`
}

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a inference resource.
func (r commonResult) Extract() (*InferenceDeployment, error) {
	var s InferenceDeployment
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type ListResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a inference deployment resource.
func (r ListResult) Extract() ([]InferenceDeployment, error) {
	var s []InferenceDeployment
	err := r.ExtractInto(&s)
	return s, err
}

func (r ListResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "results")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a InferenceDeployment.
type GetResult struct {
	commonResult
}
