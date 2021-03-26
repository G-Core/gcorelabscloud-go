package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListPolicyResponse = `
{
  "count": 1,
  "results": [
    {
      "action": "REDIRECT_TO_URL",
      "created_at": "2020-09-14T14:45:30+0000",
      "id": "9b4b9a23-ccac-4945-bcdd-b0e793c12cd9",
      "listener_id": "0388b5e5-3393-4aa8-a88a-dbcdcedf9970",
      "name": "redirect-example.com",
      "operating_status": "ONLINE",
      "position": 1,
      "project_id": 1,
      "provisioning_status": "ACTIVE",
      "redirect_http_code": 301,
      "redirect_pool_id": null,
      "redirect_prefix": null,
      "redirect_url": "http://www.example.com",
      "region": "Luxembourg",
      "region_id": 1,
      "rules": [
        {
          "compare_type": "STARTS_WITH",
          "created_at": "2020-09-14T14:45:30+0000",
          "id": "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7",
          "invert": false,
          "key": null,
          "operating_status": "ONLINE",
          "project_id": 1,
          "provisioning_status": "ACTIVE",
          "region": "Luxembourg",
          "region_id": 1,
          "tags": [
            "test_tag"
          ],
          "type": "PATH",
          "updated_at": "2020-09-14T14:45:31+0000",
          "value": "/images*"
        }
      ],
      "tags": [
        "test_tag"
      ],
      "updated_at": "2020-09-14T14:45:31+0000"
    }
  ]
}
`

const GetPolicyResponse = `
    {
      "action": "REDIRECT_TO_URL",
      "created_at": "2020-09-14T14:45:30+0000",
      "id": "9b4b9a23-ccac-4945-bcdd-b0e793c12cd9",
      "listener_id": "0388b5e5-3393-4aa8-a88a-dbcdcedf9970",
      "name": "redirect-example.com",
      "operating_status": "ONLINE",
      "position": 1,
      "project_id": 1,
      "provisioning_status": "ACTIVE",
      "redirect_http_code": 301,
      "redirect_pool_id": null,
      "redirect_prefix": null,
      "redirect_url": "http://www.example.com",
      "region": "Luxembourg",
      "region_id": 1,
      "rules": [
        {
          "compare_type": "STARTS_WITH",
          "created_at": "2020-09-14T14:45:30+0000",
          "id": "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7",
          "invert": false,
          "key": null,
          "operating_status": "ONLINE",
          "project_id": 1,
          "provisioning_status": "ACTIVE",
          "region": "Luxembourg",
          "region_id": 1,
          "tags": [
            "test_tag"
          ],
          "type": "PATH",
          "updated_at": "2020-09-14T14:45:31+0000",
          "value": "/images*"
        }
      ],
      "tags": [
        "test_tag"
      ],
      "updated_at": "2020-09-14T14:45:31+0000"
    }
`

const CreatePolicyRequest = `
{
  "action": "REDIRECT_TO_URL",
  "listener_id": "023f2e34-7806-443b-bfae-16c324569a3d",
  "name": "redirect-example.com",
  "position": 1,
  "redirect_http_code": 301,
  "redirect_url": "http://www.example.com",
  "tags": [
    "test_tag"
  ]
}
`

const ReplacePolicyRequest = `
{
  "action": "REDIRECT_TO_URL",
  "name": "redirect-example.com",
  "position": 1,
  "redirect_http_code": 301,
  "redirect_url": "http://www.example.com",
  "tags": [
    "test_tag"
  ]
}
`

const TaskResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ListRuleResponse = `
{
  "count": 1,
  "results": [
	{
	  "compare_type": "STARTS_WITH",
	  "created_at": "2020-09-14T14:45:30+0000",
	  "id": "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7",
	  "invert": false,
	  "key": null,
	  "operating_status": "ONLINE",
	  "project_id": 1,
	  "provisioning_status": "ACTIVE",
	  "region": "Luxembourg",
	  "region_id": 1,
	  "tags": [
		"test_tag"
	  ],
	  "type": "PATH",
	  "updated_at": "2020-09-14T14:45:31+0000",
	  "value": "/images*"
	}
  ]
}
`

const GetRuleResponse = `
	{
	  "compare_type": "STARTS_WITH",
	  "created_at": "2020-09-14T14:45:30+0000",
	  "id": "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7",
	  "invert": false,
	  "key": null,
	  "operating_status": "ONLINE",
	  "project_id": 1,
	  "provisioning_status": "ACTIVE",
	  "region": "Luxembourg",
	  "region_id": 1,
	  "tags": [
		"test_tag"
	  ],
	  "type": "PATH",
	  "updated_at": "2020-09-14T14:45:31+0000",
	  "value": "/images*"
	}
`

const CreateRuleRequest = `
{
  "compare_type": "REGEX",
  "invert": false,
  "tags": [
    "test_tag"
  ],
  "type": "PATH",
  "value": "/images*"
}
`

var (
	createdTimeString    = "2020-09-14T14:45:30+0000"
	updatedTimeString    = "2020-09-14T14:45:31+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
	updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
	updatedTime          = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
	redirectHttpCode     = 301
	redirectURL          = "http://www.example.com"
	pid                  = "9b4b9a23-ccac-4945-bcdd-b0e793c12cd9"
	rid                  = "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7"

	Rule = l7policies.L7Rule{
		ID:                 "0ca7bebd-7a54-4977-bca7-e4ac1e612ec7",
		CompareType:        l7policies.CompareTypeStartWith,
		Invert:             false,
		OperatingStatus:    "ONLINE",
		ProvisioningStatus: "ACTIVE",
		CreatedAt:          createdTime,
		UpdatedAt:          &updatedTime,
		Type:               l7policies.TypePath,
		Value:              "/images*",
		Tags:               []string{"test_tag"},
		ProjectID:          1,
		RegionID:           1,
		Region:             "Luxembourg",
	}
	Policy = l7policies.L7Policy{
		ID:                 "9b4b9a23-ccac-4945-bcdd-b0e793c12cd9",
		Name:               "redirect-example.com",
		Action:             l7policies.ActionRedirectToURL,
		ListenerID:         "0388b5e5-3393-4aa8-a88a-dbcdcedf9970",
		Position:           1,
		ProjectID:          1,
		RegionID:           1,
		Region:             "Luxembourg",
		OperatingStatus:    "ONLINE",
		ProvisioningStatus: "ACTIVE",
		RedirectHttpCode:   &redirectHttpCode,
		RedirectURL:        &redirectURL,
		Tags:               []string{"test_tag"},
		CreatedAt:          createdTime,
		UpdatedAt:          &updatedTime,
		Rules:              []l7policies.L7Rule{Rule},
	}
	ExpectedPolicySlice = []l7policies.L7Policy{Policy}
	ExpectedRuleSlice   = []l7policies.L7Rule{Rule}
	Tasks1              = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}
)
