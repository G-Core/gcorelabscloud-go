package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "description": "Default security group",
      "updated_at": "2019-07-26T13:25:03+0000",
      "name": "default",
      "security_group_rules": [
        {
          "description": null,
          "direction": "egress",
          "port_range_max": null,
          "updated_at": "2019-07-26T13:25:03+0000",
          "remote_group_id": null,
          "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
          "protocol": "51",
          "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
          "remote_ip_prefix": null,
          "port_range_min": null,
          "revision_number": 0,
          "created_at": "2019-07-26T13:25:03+0000",
          "ethertype": "IPv4"
        },
        {
          "description": null,
          "direction": "egress",
          "port_range_max": null,
          "updated_at": "2019-07-26T13:25:03+0000",
          "remote_group_id": null,
          "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
          "protocol": "50",
          "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
          "remote_ip_prefix": null,
          "port_range_min": null,
          "revision_number": 0,
          "created_at": "2019-07-26T13:25:03+0000",
          "ethertype": "IPv4"
        }
      ],
      "id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
      "revision_number": 4,
      "created_at": "2019-07-26T13:25:03+0000",
      "region": "Luxembourg 1",
      "project_id": 1,
      "region_id": 1
    }
  ]
}
`

const GetResponse = `
{
  "description": "Default security group",
  "updated_at": "2019-07-26T13:25:03+0000",
  "name": "default",
  "security_group_rules": [
	{
	  "description": null,
	  "direction": "egress",
	  "port_range_max": null,
	  "updated_at": "2019-07-26T13:25:03+0000",
	  "remote_group_id": null,
	  "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
	  "protocol": "51",
	  "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
	  "remote_ip_prefix": null,
	  "port_range_min": null,
	  "revision_number": 0,
	  "created_at": "2019-07-26T13:25:03+0000",
	  "ethertype": "IPv4"
	},
	{
	  "description": null,
	  "direction": "egress",
	  "port_range_max": null,
	  "updated_at": "2019-07-26T13:25:03+0000",
	  "remote_group_id": null,
	  "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
	  "protocol": "50",
	  "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
	  "remote_ip_prefix": null,
	  "port_range_min": null,
	  "revision_number": 0,
	  "created_at": "2019-07-26T13:25:03+0000",
	  "ethertype": "IPv4"
	}
  ],
  "id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
  "revision_number": 4,
  "created_at": "2019-07-26T13:25:03+0000",
  "region": "Luxembourg 1",
  "project_id": 1,
  "region_id": 1
}
`

const CreateRequest = `
{
  "instances": [
    "8dc30d49-bb34-4920-9bbd-03a2587ec0ad",
    "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
  ],
  "security_group": {
    "description": "Default security group",
    "name": "default",
    "security_group_rules": []
  }
}
`

const CreateRuleRequest = `
{
  "description": "Default security group",
  "direction": "egress",
  "ethertype": "IPv4",
  "protocol": "tcp"
}
`

const UpdateRequest = `
{
  "name": "default"
}
`

const CreateResponse = `
{
  "description": "Default security group",
  "updated_at": "2019-07-26T13:25:03+0000",
  "name": "default",
  "security_group_rules": [
	{
	  "description": null,
	  "direction": "egress",
	  "port_range_max": null,
	  "updated_at": "2019-07-26T13:25:03+0000",
	  "remote_group_id": null,
	  "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
	  "protocol": "51",
	  "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
	  "remote_ip_prefix": null,
	  "port_range_min": null,
	  "revision_number": 0,
	  "created_at": "2019-07-26T13:25:03+0000",
	  "ethertype": "IPv4"
	},
	{
	  "description": null,
	  "direction": "egress",
	  "port_range_max": null,
	  "updated_at": "2019-07-26T13:25:03+0000",
	  "remote_group_id": null,
	  "id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
	  "protocol": "50",
	  "security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
	  "remote_ip_prefix": null,
	  "port_range_min": null,
	  "revision_number": 0,
	  "created_at": "2019-07-26T13:25:03+0000",
	  "ethertype": "IPv4"
	}
  ],
  "id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
  "revision_number": 4,
  "created_at": "2019-07-26T13:25:03+0000",
  "region": "Luxembourg 1",
  "project_id": 1,
  "region_id": 1
}
`

const CreateRuleResponse = `
{
	"description": null,
	"direction": "egress",
	"port_range_max": null,
	"updated_at": "2019-07-26T13:25:03+0000",
	"remote_group_id": null,
	"id": "253c1ad7-8061-44b9-9f33-5616ad8ba5b6",
	"protocol": "tcp",
	"security_group_id": "3addc7a1-e926-46da-b5a2-eb4b2f935230",
	"remote_ip_prefix": null,
	"port_range_min": null,
	"revision_number": 0,
	"created_at": "2019-07-26T13:25:03+0000",
	"ethertype": "IPv4"
}
`

var (
	groupCreatedTimeString    = "2019-07-26T13:25:03+0000"
	groupUpdatedTimeString    = "2019-07-26T13:25:03+0000"
	groupCreatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, groupCreatedTimeString)
	groupCreatedTime          = gcorecloud.JSONRFC3339Z{Time: groupCreatedTimeParsed}
	groupUpdatedTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, groupUpdatedTimeString)
	groupUpdatedTime          = gcorecloud.JSONRFC3339Z{Time: groupUpdatedTimeParsed}
	groupName                 = "default"
	groupID                   = "3addc7a1-e926-46da-b5a2-eb4b2f935230"
	groupRuleID               = "253c1ad7-8061-44b9-9f33-5616ad8ba5b6"
	groupDescription          = "Default security group"
	eitherType                = types.EtherTypeIPv4
	direction                 = types.RuleDirectionEgress
	protocol                  = types.ProtocolTCP
	sgProto                   = types.Protocol51
	sgProto2                  = types.Protocol50

	securityGroupRule1 = securitygroups.SecurityGroupRule{
		ID:              groupRuleID,
		SecurityGroupID: groupID,
		RemoteGroupID:   nil,
		Direction:       direction,
		EtherType:       &eitherType,
		Protocol:        &protocol,
		PortRangeMax:    nil,
		PortRangeMin:    nil,
		Description:     nil,
		RemoteIPPrefix:  nil,
		CreatedAt:       groupCreatedTime,
		UpdatedAt:       &groupUpdatedTime,
		RevisionNumber:  0,
	}

	SecurityGroup1 = securitygroups.SecurityGroup{
		Name:           groupName,
		Description:    groupDescription,
		ID:             groupID,
		CreatedAt:      groupCreatedTime,
		UpdatedAt:      &groupUpdatedTime,
		RevisionNumber: 4,
		SecurityGroupRules: []securitygroups.SecurityGroupRule{
			{
				ID:              groupRuleID,
				SecurityGroupID: groupID,
				Direction:       direction,
				RemoteGroupID:   nil,
				EtherType:       &eitherType,
				Protocol:        &sgProto,
				PortRangeMax:    nil,
				PortRangeMin:    nil,
				Description:     nil,
				RemoteIPPrefix:  nil,
				CreatedAt:       groupCreatedTime,
				UpdatedAt:       &groupUpdatedTime,
				RevisionNumber:  0,
			},
			{
				ID:              groupRuleID,
				SecurityGroupID: groupID,
				Direction:       direction,
				RemoteGroupID:   nil,
				EtherType:       &eitherType,
				Protocol:        &sgProto2,
				PortRangeMax:    nil,
				PortRangeMin:    nil,
				Description:     nil,
				RemoteIPPrefix:  nil,
				CreatedAt:       groupCreatedTime,
				UpdatedAt:       &groupUpdatedTime,
				RevisionNumber:  0,
			},
		},
		ProjectID: fake.ProjectID,
		RegionID:  fake.RegionID,
		Region:    "Luxembourg 1",
	}

	ExpectedSecurityGroupSlice = []securitygroups.SecurityGroup{SecurityGroup1}
)
