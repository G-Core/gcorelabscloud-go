package clusters

import "slices"

type HighAvailabilityReplicationMode string

type PoolerMode string

type ClusterStatus string

type RoleAttribute string

const (
	HighAvailabilityReplicationModeAsync HighAvailabilityReplicationMode = "async"
	HighAvailabilityReplicationModeSync  HighAvailabilityReplicationMode = "sync"

	PoolerModeSession     PoolerMode = "session"
	PoolerModeStatement   PoolerMode = "statement"
	PoolerModeTransaction PoolerMode = "transaction"

	ClusterStatusDeleting  ClusterStatus = "DELETING"
	ClusterStatusFailed    ClusterStatus = "FAILED"
	ClusterStatusPreparing ClusterStatus = "PREPARING"
	ClusterStatusReady     ClusterStatus = "READY"
	ClusterStatusUnhealthy ClusterStatus = "UNHEALTHY"
	ClusterStatusUpdating  ClusterStatus = "UPDATING"
	ClusterStatusUnknown   ClusterStatus = "UNKNOWN"

	RoleAttributeBypassRLS  RoleAttribute = "BYPASSRLS"
	RoleAttributeCreateRole RoleAttribute = "CREATEROLE"
	RoleAttributeCreateDB   RoleAttribute = "CREATEDB"
	RoleAttributeInherit    RoleAttribute = "INHERIT"
	RoleAttributeLogin      RoleAttribute = "LOGIN"
	RoleAttributeNoLogin    RoleAttribute = "NOLOGIN"
)

func RoleAttributeList() []RoleAttribute {
	return []RoleAttribute{
		RoleAttributeBypassRLS,
		RoleAttributeCreateRole,
		RoleAttributeCreateDB,
		RoleAttributeInherit,
		RoleAttributeLogin,
		RoleAttributeNoLogin,
	}
}

func RoleAttributeStringList() []string {
	var s []string
	for _, v := range RoleAttributeList() {
		s = append(s, v.String())
	}
	return s
}

func RoleAttributeSliceToStrings(roleAttributes []RoleAttribute) []string {
	var s []string
	for _, v := range roleAttributes {
		s = append(s, v.String())
	}
	return s
}

func IsValidRoleAttribute(s string) bool {
	return slices.Contains(RoleAttributeStringList(), s)
}

func (it *RoleAttribute) String() string {
	return string(*it)
}
