package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const (
	accessibilityResponse = `
{
    "is_accessible": false,
    "message": "Forbidden",
    "http_code": 403
}`

	profileTemplatesResponse = `
{
  "count": 1,
  "results": [
    {
      "id": 1,
      "name": "test_client_profile_template",
      "description": "test client profile template",
      "fields": [
        {
          "id": 0,
          "name": "test_field",
          "description": "test field",
          "field_type": "int",
          "required": true,
          "default": "string"
        }
      ]
    }
  ]
}`

	profilesResponse = `
{
  "count": 1,
  "results": [
    {
      "id": 1,
      "profile_template": {
        "id": 1,
        "name": "test_client_profile_template",
        "description": "test client profile template",
        "fields": [
          {
            "id": 0,
            "name": "test_field",
            "description": "test field",
            "field_type": "int",
            "required": true,
            "default": "string"
          }
        ]
      },
      "ip_address": "123.123.123.1",
      "site": "example.com",
      "options": {
        "bgp": true,
        "active": true
      },
      "fields": [
        {
          "id": 0,
          "name": "test_field",
          "description": "test field",
          "field_type": "int",
          "required": true,
          "default": "string",
          "value": "string",
          "base_field": 1
        }
      ],
      "protocols": [
        {
          "port": "12000-25000",
          "protocols": ["UDP", "TCP"]
        }
      ]
    }
  ]
}
`

	createProfileRequest = `
{
  "bm_instance_id": "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
  "profile_template": 1,
  "ip_address": "123.123.123.1",
  "fields": [
    {
      "value": "string",
      "base_field": 1
    }
  ]
}
`

	updateProfileRequest = `
{
  "bm_instance_id": "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
  "profile_template": 1,
  "ip_address": "123.123.123.1",
  "fields": [
    {
      "value": "string",
      "base_field": 1
    }
  ]
}
`

	activateProfileRequest = `
{
    "active": true,
    "bgp": true
}
`

	tasksList = `
{
  "tasks": [
    "d478ae29-dedc-4869-82f0-96104425f565"
  ]
}
`

	regionCoverageResponse = `
{
    "is_covered": true
}
`
)

var (
	accessStatus = ddos.AccessStatus{
		IsAccessible: false,
		Message:      "Forbidden",
		HTTPCode:     403,
	}

	regionCoverage = ddos.RegionCoverage{
		IsCovered: true,
	}

	defaultString    = "string"
	profileTemplates = ddos.ProfileTemplate{
		ID:          1,
		Name:        "test_client_profile_template",
		Description: "test client profile template",
		Fields: []ddos.TemplateField{
			{
				ID:          0,
				Name:        "test_field",
				Description: "test field",
				FieldType:   "int",
				Required:    true,
				Default:     defaultString,
			},
		},
	}

	profile = ddos.Profile{
		ID: 1,
		ProfileTemplate: ddos.ProfileTemplate{
			ID:          1,
			Name:        "test_client_profile_template",
			Description: "test client profile template",
			Fields: []ddos.TemplateField{
				{
					ID:          0,
					Name:        "test_field",
					Description: "test field",
					FieldType:   "int",
					Required:    true,
					Default:     defaultString,
				},
			},
		},
		IPAddress: "123.123.123.1",
		Options: ddos.Options{
			Active: true,
			BGP:    true,
		},
		Site: "example.com",
		Fields: []ddos.ProfileField{
			{
				ID:          0,
				Name:        "test_field",
				Description: "test field",
				FieldType:   "int",
				Required:    true,
				Default:     defaultString,
				Value:       "string",
				BaseField:   1,
			},
		},
		Protocols: []ddos.Protocol{
			{
				Port:      "12000-25000",
				Protocols: []string{"UDP", "TCP"},
			},
		},
	}

	task = tasks.TaskResults{
		Tasks: []tasks.TaskID{
			"d478ae29-dedc-4869-82f0-96104425f565",
		},
	}
)
