package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/inferences"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const TasksResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "project_id": 1,
      "name": "test-inf",
      "image": "nginx:latest",
      "listening_port": 8080,
      "created_at": "2023-08-22T11:21:00Z",
      "status": "ACTIVE",
      "auth_enabled": false,
      "address": "https://example.com",
      "containers": [
        {
          "deploy_status": {
            "ready": 1,
            "total": 3
          },
          "region_id": 1,
          "scale": {
            "cooldown_period": 60,
            "max": 3,
            "min": 1,
            "triggers": {
              "cpu": {
                "threshold": 80
              },
              "memory": {
                "threshold": 70
              }
            }
          }
        }
      ],
      "timeout": 120,
      "envs": {
        "DEBUG_MODE": "False",
        "KEY": "12345"
      },
      "flavor_name": "inference-16vcpu-232gib-1xh100-80gb",
      "command": [
        "nginx",
        "-g",
        "daemon off;"
      ],
      "credentials_name": "dockerhub",
      "logging": {
        "destination_region_id": 1,
        "enabled": true,
        "opensearch_dashboards_link": "https://opensearch.gcore.com",
        "retention_policy": 30,
        "topic_name": "mynamespace.topic"
      }
    }
  ]
}
`

const CreateRequest = `
{
  "name": "test-inf",
  "image": "nginx:latest",
  "listening_port": 8080,
  "description": "",
  "auth_enabled": false,
  "containers": [
    {
      "region_id": 1,
      "scale": {
        "cooldown_period": 60,
        "max": 3,
        "min": 1,
        "triggers": {
          "cpu": {
            "threshold": 80
          },
          "memory": {
            "threshold": 70
          }
        }
      }
    }
  ],
  "timeout": 120,
  "envs": {
    "DEBUG_MODE": "False",
    "KEY": "12345"
  },
  "flavor_name": "inference-16vcpu-232gib-1xh100-80gb",
  "command": [
    "nginx",
    "-g",
    "daemon off;"
  ],
  "credentials_name": "dockerhub",
  "logging": {
    "destination_region_id": 1,
    "enabled": true,
    "retention_policy": {
      "period": 30
    },
    "topic_name": "mynamespace.topic"
  }
}
`

const UpdateRequest = `
{
  "image": "nginx:latest",
  "listening_port": 8080,
  "auth_enabled": false,
  "containers": [
    {
      "region_id": 1,
      "scale": {
        "cooldown_period": 60,
        "max": 3,
        "min": 1,
        "triggers": {
          "cpu": {
            "threshold": 80
          },
          "memory": {
            "threshold": 70
          }
        }
      }
    }
  ],
  "timeout": 120,
  "envs": {
    "DEBUG_MODE": "False",
    "KEY": "12345"
  },
  "flavor_name": "inference-16vcpu-232gib-1xh100-80gb",
  "command": [
    "nginx",
    "-g",
    "daemon off;"
  ],
  "credentials_name": "dockerhub",
  "logging": {
    "destination_region_id": 1,
    "enabled": true,
    "retention_policy": {
      "period": 30
    },
    "topic_name": "mynamespace.topic"
  }
}
`

const GetResponse = `
{
  "project_id": 1,
  "name": "test-inf",
  "image": "nginx:latest",
  "listening_port": 8080,
  "created_at": "2023-08-22T11:21:00Z",
  "status": "ACTIVE",
  "auth_enabled": false,
  "address": "https://example.com",
  "containers": [
	{
	  "deploy_status": {
		"ready": 1,
		"total": 3
	  },
	  "region_id": 1,
	  "scale": {
		"cooldown_period": 60,
		"max": 3,
		"min": 1,
		"triggers": {
		  "cpu": {
			"threshold": 80
		  },
		  "memory": {
			"threshold": 70
		  }
		}
	  }
	}
  ],
  "timeout": 120,
  "envs": {
	"DEBUG_MODE": "False",
	"KEY": "12345"
  },
  "flavor_name": "inference-16vcpu-232gib-1xh100-80gb",
  "command": [
	"nginx",
	"-g",
	"daemon off;"
  ],
  "credentials_name": "dockerhub",
  "logging": {
	"destination_region_id": 1,
	"enabled": true,
	"opensearch_dashboards_link": "https://opensearch.gcore.com",
	"retention_policy": 30,
	"topic_name": "mynamespace.topic"
  }
}
`

var (
	image           = "nginx:latest"
	listeningPort   = 8080
	enableAuth      = false
	flavorName      = "inference-16vcpu-232gib-1xh100-80gb"
	credentialsName = "dockerhub"
	createdAt       = "2023-08-22T11:21:00Z"
	cooldownPeriod  = 60
	regionID        = fake.RegionID
	retentionPolicy = 30
	topicName       = "mynamespace.topic"
	timeout         = 120
	Inference1      = inferences.InferenceDeployment{
		ProjectID:     fake.ProjectID,
		Name:          "test-inf",
		Image:         "nginx:latest",
		ListeningPort: listeningPort,
		CreatedAt:     &createdAt,
		Status:        "ACTIVE",
		AuthEnabled:   false,
		Address:       "https://example.com",
		Containers: []inferences.Container{
			{
				DeployStatus: inferences.ContainerDeployStatus{
					Ready: 1,
					Total: 3,
				},
				RegionID: fake.RegionID,
				Scale: inferences.ContainerScale{
					CooldownPeriod: &cooldownPeriod,
					Max:            3,
					Min:            1,
					Triggers: inferences.ContainerScaleTrigger{
						Cpu: &inferences.ScaleTriggerThreshold{
							Threshold: 80,
						},
						Memory: &inferences.ScaleTriggerThreshold{
							Threshold: 70,
						},
					},
				},
			},
		},
		Timeout: timeout,
		Envs: map[string]string{
			"DEBUG_MODE": "False",
			"KEY":        "12345",
		},
		FlavorName: "inference-16vcpu-232gib-1xh100-80gb",
		Command: []string{
			"nginx",
			"-g",
			"daemon off;",
		},
		CredentialsName: "dockerhub",
		Logging: &inferences.Logging{
			DestinationRegionID:      &regionID,
			Enabled:                  true,
			OpensearchDashboardsLink: "https://opensearch.gcore.com",
			RetentionPolicy:          &retentionPolicy,
			TopicName:                &topicName,
		},
	}
	InferencesSlice = []inferences.InferenceDeployment{Inference1}
	Tasks1          = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)
