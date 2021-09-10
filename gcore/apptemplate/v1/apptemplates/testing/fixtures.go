package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/apptemplate/v1/apptemplates"

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "app_config": [
        {
          "default": "Server name",
          "description": "tf2_hostname.description",
          "label": "tf2_hostname.label",
          "name": "tf2_hostname",
          "required": true,
          "type": "string"
        }
      ],
      "category": "gaming",
      "description": "tf2server.description",
      "developer": "Team Fortress 2 team",
      "display_name": "Team Fortress 2 server",
      "id": "tf2server",
      "image_name": "ubuntu-20.04-x64",
      "min_disk": 30,
      "min_ram": 1024,
      "min_vcpus": null,
      "os_name": "Ubuntu 20.04",
      "os_version": "20.04",
      "region_id": null,
      "short_description": "tf2server.short_description",
      "usage": "tf2server.usage",
      "version": "latest",
      "website": "https://www.teamfortress.com/"
    }
  ]
}
`

const GetResponse = `
{
  "app_config": [
    {
      "default": "Server name",
      "description": "tf2_hostname.description",
      "label": "tf2_hostname.label",
      "name": "tf2_hostname",
      "required": true,
      "type": "string"
    }
  ],
  "category": "gaming",
  "description": "tf2server.description",
  "developer": "Team Fortress 2 team",
  "display_name": "Team Fortress 2 server",
  "id": "tf2server",
  "image_name": "ubuntu-20.04-x64",
  "min_disk": 30,
  "min_ram": 1024,
  "min_vcpus": null,
  "os_name": "Ubuntu 20.04",
  "os_version": "20.04",
  "region_id": null,
  "short_description": "tf2server.short_description",
  "usage": "tf2server.usage",
  "version": "latest",
  "website": "https://www.teamfortress.com/"
}
`

var (
	AppTemplate1 = apptemplates.AppTemplate{
		ID:               "tf2server",
		OsName:           "Ubuntu 20.04",
		Developer:        "Team Fortress 2 team",
		OsVersion:        "20.04",
		Category:         "gaming",
		Website:          "https://www.teamfortress.com/",
		DisplayName:      "Team Fortress 2 server",
		ImageName:        "ubuntu-20.04-x64",
		Usage:            "tf2server.usage",
		Description:      "tf2server.description",
		ShortDescription: "tf2server.short_description",
		MinRam:           1024,
		AppConfig: []map[string]interface{}{
			{
				"default":     "Server name",
				"description": "tf2_hostname.description",
				"label":       "tf2_hostname.label",
				"name":        "tf2_hostname",
				"required":    true,
				"type":        "string",
			},
		},
		Version: "latest",
		MinDisk: 30,
	}
	ExpectedAppTemplateSlice = []apptemplates.AppTemplate{AppTemplate1}
)
