package testing

import (
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
)

const ReplaceRuleRequest = `
{
  "description": "Default security group",
  "direction": "egress",
  "ethertype": "IPv4",
  "protocol": "tcp"
}
`
const ReplaceRuleResponse = `
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
	groupID                   = "3addc7a1-e926-46da-b5a2-eb4b2f935230"
	groupRuleID               = "253c1ad7-8061-44b9-9f33-5616ad8ba5b6"
	groupDescription          = "Default security group"
	eitherType                = types.EtherTypeIPv4
	direction                 = types.RuleDirectionEgress
	protocol                  = types.ProtocolTCP

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
)
