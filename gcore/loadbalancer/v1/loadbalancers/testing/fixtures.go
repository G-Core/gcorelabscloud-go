package testing

import (
	"fmt"
	"net"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata/testing"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const VrrpIPsResponse = `
{
  "ip_address": "10.94.79.54",
  "subnet_id": "db5ebada-a86a-4702-8a19-00b23a1acb05"
},
{
  "ip_address": "aa:bb:cc:dd::1d4",
  "subnet_id": "abd99b68-e139-4715-b8c2-37ca324285b8"
},
{
  "ip_address": "10.94.76.179",
  "subnet_id": "db5ebada-a86a-4702-8a19-00b23a1acb05"
},
{
  "ip_address": "aa:bb:cc:dd::2b5",
  "subnet_id": "abd99b68-e139-4715-b8c2-37ca324285b8"
}
`

var ListResponse = fmt.Sprintf(`
{
  "count": 1,
  "results": [
    {
      "region": "RegionOne",
      "created_at": "2020-01-24T13:57:12+0000",
      "name": "lbname",
      "id": "e8ab1be4-1521-4266-be69-28dad4148a30",
      "provisioning_status": "ACTIVE",
      "updated_at": "2020-01-24T13:57:35+0000",
      "listeners": [
        {
          "id": "43658ea9-54bd-4807-90b1-925921c9a0d1"
        }
      ],
      "task_id": null,
      "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
      "vip_address": "5.5.5.5",
      "operating_status": "ONLINE",
      "project_id": 1,
      "region_id": 1,
 	    "metadata": [%s],
      "vrrp_ips": [%s],
      "vip_ip_family": "dual",
      "additional_vips": [
        {
          "ip_address": "aa:bb:cc:dd::29d",
          "subnet_id": "abd99b68-e139-4715-b8c2-37ca324285b8"
        }
      ]
    }
  ]
}
`, testing.MetadataResponse, VrrpIPsResponse)

const ListCustomSecurityGroupResponse = `
{
  "count": 1,
  "results": [
    {
      "id": "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
      "name": "Test"
    }
  ]
}
`

var GetResponse = fmt.Sprintf(`
{
  "region": "RegionOne",
  "created_at": "2020-01-24T13:57:12+0000",
  "name": "lbname",
  "id": "e8ab1be4-1521-4266-be69-28dad4148a30",
  "provisioning_status": "ACTIVE",
  "updated_at": "2020-01-24T13:57:35+0000",
  "listeners": [
    {
      "id": "43658ea9-54bd-4807-90b1-925921c9a0d1"
    }
  ],
  "task_id": null,
  "creator_task_id": "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22",
  "vip_address": "5.5.5.5",
  "operating_status": "ONLINE",
  "project_id": 1,
  "region_id": 1,
  "metadata": [%s],
  "vrrp_ips": [%s],
  "vip_ip_family": "dual",
  "additional_vips": [
    {
      "ip_address": "aa:bb:cc:dd::29d",
      "subnet_id": "abd99b68-e139-4715-b8c2-37ca324285b8"
    }
  ]
}
`, testing.MetadataResponse, VrrpIPsResponse)

const CreateRequest = `
{
  "name": "lbname",
  "vip_port_id": "169942e0-9b53-42df-95ef-1a8b6525c2bd",
  "vip_ip_family": "dual",
  "listeners": [
    {
      "name": "listener_name",
      "insert_x_forwarded": false,
      "protocol": "HTTP",
      "protocol_port": 80,
      "pools": [
        {
          "name": "pool_name",
          "protocol": "HTTP",
          "members": [
            {
              "instance_id": "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9",
              "address": "192.168.1.101",
              "weight": 2,
              "protocol_port": 8000
            },
            {
              "instance_id": "169942e0-9b53-42df-95ef-1a8b6525c2bd",
              "address": "192.168.1.102",
              "weight": 4,
              "protocol_port": 8000
            }
          ],
          "healthmonitor": {
            "type": "HTTP",
            "delay": 10,
            "max_retries": 3,
            "timeout": 5,
            "max_retries_down": 3,
            "http_method": "GET",
            "url_path": "/"
          },
          "lb_algorithm": "ROUND_ROBIN"
        }
      ]
    }
  ]
}
`

const UpdateRequest = `
{
	"name": "lbname"
}
`

const ResizeRequest = `
{
	"flavor": "lb1-4-8"
}
`

const CreateResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const DeleteResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ResizeResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

var (
	createdTimeString    = "2020-01-24T13:57:12+0000"
	updatedTimeString    = "2020-01-24T13:57:35+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}
	updatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, updatedTimeString)
	updatedTime          = gcorecloud.JSONRFC3339Z{Time: updatedTimeParsed}
	creatorTaskID        = "9f3ec11e-bcd4-4fe6-924a-a4439a56ad22"

	LoadBalancer1 = loadbalancers.LoadBalancer{
		Name:               "lbname",
		ID:                 "e8ab1be4-1521-4266-be69-28dad4148a30",
		ProvisioningStatus: types.ProvisioningStatusActive,
		OperationStatus:    types.OperatingStatusOnline,
		VipAddress:         net.ParseIP("5.5.5.5"),
		Listeners: []gcorecloud.ItemID{{
			ID: "43658ea9-54bd-4807-90b1-925921c9a0d1",
		}},
		CreatorTaskID: &creatorTaskID,
		TaskID:        nil,
		CreatedAt:     createdTime,
		UpdatedAt:     &updatedTime,
		ProjectID:     fake.ProjectID,
		RegionID:      fake.RegionID,
		Region:        "RegionOne",
		Metadata:      []metadata.Metadata{testing.ResourceMetadataReadOnly},
		VrrpIPs: []loadbalancers.NetworkPortFixedIP{
			{IpAddress: net.ParseIP("10.94.79.54"), SubnetID: "db5ebada-a86a-4702-8a19-00b23a1acb05"},
			{IpAddress: net.ParseIP("aa:bb:cc:dd::1d4"), SubnetID: "abd99b68-e139-4715-b8c2-37ca324285b8"},
			{IpAddress: net.ParseIP("10.94.76.179"), SubnetID: "db5ebada-a86a-4702-8a19-00b23a1acb05"},
			{IpAddress: net.ParseIP("aa:bb:cc:dd::2b5"), SubnetID: "abd99b68-e139-4715-b8c2-37ca324285b8"},
		},
		VipIPFamilyType: types.DualStackIPFamilyType,
		AdditionalVips: []loadbalancers.NetworkPortFixedIP{
			{IpAddress: net.ParseIP("aa:bb:cc:dd::29d"), SubnetID: "abd99b68-e139-4715-b8c2-37ca324285b8"},
		},
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedLoadBalancerSlice = []loadbalancers.LoadBalancer{LoadBalancer1}
	LbSecurityGroup1          = loadbalancers.CustomSecurityGroup{
		ID:   "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
		Name: "Test",
	}
	ExpectedLbSecurityGroupSlice = []loadbalancers.CustomSecurityGroup{LbSecurityGroup1}
	Flavor                       = "lb1-4-8"
)
