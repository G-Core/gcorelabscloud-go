package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/extensions"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "name": "name",
      "alias": "alias",
      "description": "description",
      "links": [
        "http://test.com"
      ],
      "updated": "2006-01-02T15:04:05-0700"
    }
  ]
}
`

const GetResponse = `
{
  "name": "name",
  "alias": "alias",
  "description": "description",
  "links": [
	"http://test.com"
  ],
  "updated": "2006-01-02T15:04:05-0700"
}
`

var updatedTimeString = "2006-01-02T15:04:05-0700"
var updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
var updatedTime = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}

var (
	Extension1 = extensions.Extension{
		Name:        "name",
		Alias:       "alias",
		Links:       []string{"http://test.com"},
		Description: "description",
		Updated:     updatedTime,
	}

	ExpectedExtensionSlice = []extensions.Extension{Extension1}
)
