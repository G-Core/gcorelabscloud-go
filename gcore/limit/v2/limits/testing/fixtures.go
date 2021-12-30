package testing

import (
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v2/limits"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v2/types"
)

const CreateRequest = `
{
  "description": "test",	
  "requested_limits": {
	"global_limits": {},
	"regional_limits": [
 	  {
		"baremetal_basic_count_limit": 0,
		"baremetal_hf_count_limit": 0,
		"baremetal_infrastructure_count_limit": 0,
		"baremetal_network_count_limit": 0,
		"baremetal_storage_count_limit": 0,
		"cluster_count_limit": 0,
		"cpu_count_limit": 1,
		"external_ip_count_limit": 0,
		"firewall_count_limit": 13,
		"floating_count_limit": 0,
		"gpu_count_limit": 0,
		"image_count_limit": 0,
		"image_size_limit": 0,
		"loadbalancer_count_limit": 0,
		"network_count_limit": 0,
		"ram_limit": 0,
		"region_id": 1,
		"router_count_limit": 0,
		"secret_count_limit": 0,
		"servergroup_count_limit": 0,
		"shared_vm_count_limit": 0,
		"snapshot_schedule_count_limit": 0,
		"subnet_count_limit": 0,
		"vm_count_limit": 0,
		"volume_count_limit": 0,
		"volume_size_limit": 0,
		"volume_snapshots_count_limit": 0,
		"volume_snapshots_size_limit": 0
	  }
	]
  }
}
`

const ListResponse = `
{
  "count": 0,
  "results": [
    {
      "id": 1,
      "client_id": 1,
	  "requested_limits": {
		"regional_limits": [{"region_id": 1, "firewall_count_limit": 13, "cpu_count_limit": 1}]
	  },
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
  "requested_limits": {
    "regional_limits": [{"region_id": 1, "firewall_count_limit": 13, "cpu_count_limit": 1}]
  },
  "status": "in progress",
  "created_at": "2019-07-26T13:25:03"
}
`
const CreateResponse = GetResponse

var (
	createdTimeString    = "2019-07-26T13:25:03"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339NoZ, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339NoZ{Time: createdTimeParsed}
	LimitRequest1        = limits.LimitResponse{
		ID:       1,
		ClientID: 1,
		RequestedLimits: limits.Limit{RegionalLimits: []limits.RegionalLimits{
			{RegionID: 1, FirewallCountLimit: 13, CPUCountLimit: 1},
		}},
		Status:    types.LimitRequestInProgress,
		CreatedAt: createdTime,
	}

	ExpectedLimitRequestSlice = []limits.LimitResponse{LimitRequest1}
)
