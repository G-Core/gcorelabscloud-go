package testing

import (
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/limits"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/types"
)

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "id": 1,
      "client_id": 1,
      "limits": "{\"volume_snapshots_count_limit\": 12, \"router_count_limit\": 1}",
      "status": "in progress",
      "created_at": "2019-07-26T13:25:03"
    }
  ]
}
`

const GetResponse = `
{
  "id": 1,
  "client_id": 1,
  "limits": "{\"volume_snapshots_count_limit\": 12, \"router_count_limit\": 1}",
  "status": "in progress",
  "created_at": "2019-07-26T13:25:03"
}
`

const CreateRequest = `
{
  "description": "test",	
  "requested_quotas": {
    "external_ip_count_limit": 4
  }
}	
`

const UpdateRequest = `
{
  "external_ip_count_limit": 4
}	
`

const StatusRequest = `
{
  "status": "done"
}	
`

const CreateResponse = GetResponse
const UpdateResponse = GetResponse
const StatusResponse = GetResponse

var (
	createdTimeString    = "2019-07-26T13:25:03"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339NoZ, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339NoZ{Time: createdTimeParsed}
	LimitRequest1        = limits.LimitResponse{
		ID:        1,
		ClientID:  1,
		Limits:    "{\"volume_snapshots_count_limit\": 12, \"router_count_limit\": 1}",
		Status:    types.LimitRequestInProgress,
		CreatedAt: createdTime,
	}

	ExpectedLimitRequestSlice = []limits.LimitResponse{LimitRequest1}
)
