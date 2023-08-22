package testing

import "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"

const MetadataResponse = `
{
  "key": "some_key",
  "value": "some_val",
  "read_only": false
}
`
const MetadataCreateRequest = `
{
"test1": "test1", 
"test2": "test2"
}
`
const MetadataListResponse = `
{
  "count": 2,
  "results": [
    {
      "key": "cost-center",
      "value": "Atlanta",
      "read_only": false
    },
    {
      "key": "data-center",
      "value": "A",
      "read_only": false
    }
  ]
}
`

var (
	ResourceMetadata = map[string]interface{}{
		"some_key": "some_val",
	}

	ResourceMetadataReadOnly = metadata.Metadata{
		Key:      "some_key",
		Value:    "some_val",
		ReadOnly: false,
	}

	Metadata1 = metadata.Metadata{
		Key:      "cost-center",
		Value:    "Atlanta",
		ReadOnly: false,
	}
	Metadata2 = metadata.Metadata{
		Key:      "data-center",
		Value:    "A",
		ReadOnly: false,
	}
	ExpectedMetadataList = []metadata.Metadata{Metadata1, Metadata2}
)
