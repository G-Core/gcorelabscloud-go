package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/keystone/v1/keystones"

	keystonetypes "github.com/G-Core/gcorelabscloud-go/gcore/keystone/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/types"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "display_name": "ED-10",
      "created_on": "2020-04-10T11:37:58",
      "keystone_id": 1,
      "id": 1,
      "keystone": {
        "created_on": "2020-04-10T11:37:58",
        "admin_password": "******",
        "id": 1,
        "keystone_federated_domain_id": "5ac0a17e556d4a9c8f946928a7953990",
        "state": "NEW",
        "url": "https://ed-10.cloud.core.pw:5000/v3"
      },
      "state": "ACTIVE",
	  "task_id": null,
      "external_network_id": "0521f854-8e34-4e67-8827-2aeb27fb3872",
      "spice_proxy_url": "https://ed-10.cloud.core.pw:6062",
      "endpoint_type": "public",
      "keystone_name": "ED-10"
    }
  ]
}
`

const GetResponse = `
{
  "display_name": "ED-10",
  "created_on": "2020-04-10T11:37:58",
  "keystone_id": 1,
  "id": 1,
  "keystone": {
    "created_on": "2020-04-10T11:37:58",
    "admin_password": "******",
    "id": 1,
    "keystone_federated_domain_id": "5ac0a17e556d4a9c8f946928a7953990",
    "state": "NEW",
    "url": "https://ed-10.cloud.core.pw:5000/v3"
  },
  "state": "ACTIVE",
  "task_id": null,
  "external_network_id": "0521f854-8e34-4e67-8827-2aeb27fb3872",
  "spice_proxy_url": "https://ed-10.cloud.core.pw:6062",
  "endpoint_type": "public",
  "keystone_name": "ED-10"
}
`

const CreateRequest = `
{
	"display_name": "ED-10",
	"endpoint_type": "public",
	"external_network_id": "0521f854-8e34-4e67-8827-2aeb27fb3872",
	"keystone_id": 1,
	"keystone_name": "ED-10",
	"state": "ACTIVE"
}
`

const UpdateRequest = `
{
	"display_name": "ED-10",
	"state": "DELETED"
}	
`

const CreateResponse = GetResponse
const UpdateResponse = GetResponse

var (
	createdTimeString    = "2020-04-10T11:37:58"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339NoZ, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339NoZ{Time: createdTimeParsed}
	keystoneURL, _       = gcorecloud.ParseURL("https://ed-10.cloud.core.pw:5000/v3")
	spiceURL, _          = gcorecloud.ParseURL("https://ed-10.cloud.core.pw:6062")

	Region1 = regions.Region{
		ID:                1,
		DisplayName:       "ED-10",
		KeystoneName:      "ED-10",
		State:             types.RegionStateActive,
		TaskID:            nil,
		EndpointType:      types.EndpointTypePublic,
		ExternalNetworkID: "0521f854-8e34-4e67-8827-2aeb27fb3872",
		SpiceProxyURL:     *spiceURL,
		CreatedOn:         createdTime,
		KeystoneID:        1,
		Keystone: keystones.Keystone{
			ID:                        1,
			URL:                       *keystoneURL,
			State:                     keystonetypes.KeystoneStateNew,
			KeystoneFederatedDomainID: "5ac0a17e556d4a9c8f946928a7953990",
			CreatedOn:                 createdTime,
			AdminPassword:             "******",
		},
	}

	ExpectedRegionSlice = []regions.Region{Region1}
)
