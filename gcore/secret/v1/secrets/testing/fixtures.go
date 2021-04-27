package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var ListResponse = `
{
  "count": 1,
  "results": [
    {
      "status": "ACTIVE",
	  "created": "2025-12-28T19:14:44+00:00",
	  "updated": "2025-12-28T19:14:44+00:00",
	  "expiration": "2025-12-28T19:14:44+00:00",
      "algorithm": "aes",
      "bit_length": 256,
      "mode": "cbc",
      "name": "AES key",
      "id": "bfc7824b-31b6-4a28-a0c4-7df137139215",
      "secret_type": "opaque",
      "content_types": {
        "default": "application/octet-stream"
      }
    }
  ]
}
`

var GetResponse = `
{
  "status": "ACTIVE",
  "created": "2025-12-28T19:14:44+00:00",
  "updated": "2025-12-28T19:14:44+00:00",
  "expiration": "2025-12-28T19:14:44+00:00",
  "algorithm": "aes",
  "bit_length": 256,
  "mode": "cbc",
  "name": "AES key",
  "id": "bfc7824b-31b6-4a28-a0c4-7df137139215",
  "secret_type": "opaque",
  "content_types": {
    "default": "application/octet-stream"
  }
}
`

var CreateRequest = `
{
  "expiration": "2025-12-28T19:14:44",
  "name": "AES key",
  "payload": "aGVsbG8sIHRlc3Qgc3RyaW5nCg==",
  "payload_content_encoding": "base64",
  "payload_content_type": "application/octet-stream",
  "secret_type": "certificate"
}
`

var CreateResponse = `
{
  "tasks": [
    "d478ae29-dedc-4869-82f0-96104425f565"
  ]
}
`

var DeleteResponse = CreateResponse

var createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339ZColon, "2025-12-28T19:14:44+00:00")
var createdTime = gcorecloud.JSONRFC3339ZColon{Time: createdTimeParsed}

var (
	S1 = secrets.Secret{
		ID:           "bfc7824b-31b6-4a28-a0c4-7df137139215",
		Name:         "AES key",
		Status:       "ACTIVE",
		Algorithm:    "aes",
		BitLength:    256,
		ContentTypes: map[string]string{"default": "application/octet-stream"},
		Mode:         "cbc",
		Type:         secrets.OpaqueSecretType,
		CreatedAt:    createdTime,
		UpdatedAt:    createdTime,
		Expiration:   createdTime,
	}
	ExpectedSecretsSlice = []secrets.Secret{S1}
	Tasks1               = tasks.TaskResults{
		Tasks: []tasks.TaskID{"d478ae29-dedc-4869-82f0-96104425f565"},
	}
)
