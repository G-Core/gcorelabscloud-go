package testing

import (
	"fmt"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/dbaas/postgres/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

// prepareTestURL constructs the base URL for cluster operations
func prepareTestURL() string {
	return fmt.Sprintf("/v1/dbaas/postgres/clusters/%d/%d", fake.ProjectID, fake.RegionID)
}

// prepareClusterTestURL constructs the URL for specific cluster operations
func prepareClusterTestURL(clusterName string) string {
	return fmt.Sprintf("/v1/dbaas/postgres/clusters/%d/%d/%s", fake.ProjectID, fake.RegionID, clusterName)
}

var (
	createdTimeString    = "2024-01-15T10:00:00+0000"
	createdTimeParsed, _ = time.Parse(gcorecloud.RFC3339Z, createdTimeString)
	createdTime          = gcorecloud.JSONRFC3339Z{Time: createdTimeParsed}

	FirstClusterShort = clusters.PostgresSQLClusterShort{
		ClusterName: "test-cluster-1",
		CreatedAt:   createdTime,
		Status:      clusters.ClusterStatusReady,
		Version:     "15.0",
	}

	SecondClusterShort = clusters.PostgresSQLClusterShort{
		ClusterName: "test-cluster-2",
		CreatedAt:   createdTime,
		Status:      clusters.ClusterStatusPreparing,
		Version:     "14.0",
	}

	FirstClusterDetail = clusters.PostgresSQLCluster{
		ClusterName: "test-cluster-1",
		CreatedAt:   createdTime,
		Status:      clusters.ClusterStatusReady,
		Flavor: clusters.Flavor{
			CPU:       2,
			MemoryGiB: 4,
		},
		Storage: clusters.PGStorageConfiguration{
			SizeGiB: 50,
			Type:    "standard",
		},
		Network: clusters.Network{
			ACL:              []string{"0.0.0.0/0"},
			ConnectionString: "postgres://test-cluster-1.example.com:5432",
			Host:             "test-cluster-1.example.com",
			NetworkType:      "public",
		},
		PGServerConfiguration: clusters.PGServerConfiguration{
			PGConf:  "standard",
			Version: "15.0",
			Pooler: &clusters.Pooler{
				Mode: clusters.PoolerModeSession,
				Type: "pgbouncer",
			},
		},
		HighAvailability: &clusters.HighAvailability{
			ReplicationMode: clusters.HighAvailabilityReplicationModeAsync,
		},
		Databases: []clusters.DatabaseOverview{
			{
				Name:  "testdb",
				Owner: "testuser",
				Size:  1024,
			},
		},
		Users: []clusters.PgUserOverview{
			{
				Name:             "testuser",
				RoleAttributes:   []clusters.RoleAttribute{clusters.RoleAttributeLogin, clusters.RoleAttributeCreateDB},
				IsSecretRevealed: false,
			},
		},
	}

	ExpectedTaskResults = tasks.TaskResults{
		Tasks: []tasks.TaskID{"79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"},
	}
)

const ListResponse = `
{
  "count": 2,
  "results": [
    {
      "cluster_name": "test-cluster-1",
      "created_at": "2024-01-15T10:00:00+0000",
      "status": "READY",
      "version": "15.0"
    },
    {
      "cluster_name": "test-cluster-2",
      "created_at": "2024-01-15T10:00:00+0000",
      "status": "PREPARING",
      "version": "14.0"
    }
  ]
}
`

const ListResponsePagination = `
{
  "count": 2,
  "results": []
}
`

const GetResponse = `
{
  "cluster_name": "test-cluster-1",
  "created_at": "2024-01-15T10:00:00+0000",
  "status": "READY",
  "flavor": {
    "cpu": 2,
    "memory_gib": 4
  },
  "storage": {
    "size_gib": 50,
    "type": "standard"
  },
  "network": {
    "acl": ["0.0.0.0/0"],
    "connection_string": "postgres://test-cluster-1.example.com:5432",
    "host": "test-cluster-1.example.com",
    "network_type": "public"
  },
  "pg_server_configuration": {
    "pg_conf": "standard",
    "version": "15.0",
    "pooler": {
      "mode": "session",
      "type": "pgbouncer"
    }
  },
  "high_availability": {
    "replication_mode": "async"
  },
  "databases": [
    {
      "name": "testdb",
      "owner": "testuser",
      "size": 1024
    }
  ],
  "users": [
    {
      "name": "testuser",
      "role_attributes": ["LOGIN", "CREATEDB"],
      "is_secret_revealed": false
    }
  ]
}
`

const CreateRequest = `
{
  "cluster_name": "test-cluster-1",
  "databases": [
    {
      "name": "testdb",
      "owner": "testuser"
    }
  ],
  "flavor": {
    "cpu": 2,
    "memory_gib": 4
  },
  "high_availability": {
    "replication_mode": "async"
  },
  "network": {
    "acl": ["0.0.0.0/0"],
    "network_type": "public"
  },
  "pg_server_configuration": {
    "pg_conf": "standard",
    "version": "15.0",
    "pooler": {
      "mode": "session",
      "type": "pgbouncer"
    }
  },
  "storage": {
    "size_gib": 50,
    "type": "standard"
  },
  "users": [
    {
      "name": "testuser",
      "role_attributes": ["LOGIN", "CREATEDB"]
    }
  ]
}
`

const CreateResponse = `
{
  "tasks": ["79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"]
}
`

const UpdateRequest = `
{
  "flavor": {
    "cpu": 4,
    "memory_gib": 8
  },
  "storage": {
    "size_gib": 100
  }
}
`

const UpdateResponse = `
{
  "tasks": ["79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"]
}
`

const DeleteResponse = `
{
  "tasks": ["79dc7c30-44d2-4c5c-b5c1-e9a46f6bbf54"]
}
`
