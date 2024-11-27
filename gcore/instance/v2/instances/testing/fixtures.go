package testing

import (
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
)

const MetadataResponse = `
{
  "key": "db.name",
  "value": "pg",
  "read_only": false
}
`

var (
	instanceID = "ad1bb86e-2f83-4e0f-87c0-e1fd777d6352"
	Metadata   = metadata.Metadata{
		Key:      "db.name",
		Value:    "pg",
		ReadOnly: false,
	}
)
