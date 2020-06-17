package testing

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	uuid "github.com/satori/go.uuid"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
      "health_status": "HEALTHY",
      "name": "test",
      "create_timeout": 7200,
      "flavor_id": "g1-standard-1-2",
      "uuid": "5e09faed-e742-404f-8a75-0ea5eb3c435f",
      "docker_volume_size": 10,
      "labels": {
        "gcloud_project_id": "1",
        "gcloud_region_id": "2",
        "fixed_network_cidr": "192.168.0.0/18",
        "calico_ipv4pool": "192.168.64.0/18",
        "service_cluster_ip_range": "192.168.128.0/18"
      },
      "pools": [
        {
          "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
          "name": "test1",
          "max_node_count": 5,
          "min_node_count": 1,
          "is_default": true,
		  "docker_volume_size": 10,
		  "docker_volume_type": "standard",
          "flavor_id": "g1-standard-1-2",
          "uuid": "908338b2-9217-4673-af0e-f0093139fbac",
          "status": "CREATE_COMPLETE",
          "role": "worker",
          "image_id": "fedora-coreos",
          "node_count": 1
        }
      ],
      "master_flavor_id": "g1-standard-1-2",
      "status": "UPDATE_COMPLETE",
      "keypair": "keypair",
      "master_count": 1,
      "cluster_template_id": "2c884df0-a312-4950-a6ec-9405803affc9",
      "node_count": 1,
	  "version": "1.17"	
    }
  ]
}
`

const ClusterCsrRequest = `
{
  "csr": "string"
}
`

var ClusterSignResponse = fmt.Sprintf(`
{
  "cluster_uuid": "%s",
  "pem": "string",
  "csr": "string"
}
`, ClusterList1.UUID,
)

var ClusterCAResponse = fmt.Sprintf(`
{
  "cluster_uuid": "%s",
  "pem": "string"
}
`, ClusterList1.UUID,
)

const GetResponse = `
{
  "health_status": "HEALTHY",
  "fixed_subnet": "46beed39-38e6-4743-90b5-30fefd6173d2",
  "project_id": "46beed39-38e6-4743-90b5-30fefd6173d2",
  "pools": [
    {
      "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
      "name": "test1",
      "max_node_count": 5,
      "cluster_id": "5e09faed-e742-404f-8a75-0ea5eb3c435f",
      "role": "worker",
      "min_node_count": 1,
      "is_default": true,
	  "docker_volume_size": 10,
	  "docker_volume_type": "standard",
      "uuid": "908338b2-9217-4673-af0e-f0093139fbac",
      "labels": {
        "gcloud_project_id": "1",
        "gcloud_region_id": "2",
        "fixed_network_cidr": "192.168.0.0/18",
        "calico_ipv4pool": "192.168.64.0/18",
        "service_cluster_ip_range": "192.168.128.0/18"
      },
      "node_addresses": [
        "192.168.0.5"
      ],
      "project_id": "46beed3938e6474390b530fefd6173d2",
      "status": "CREATE_COMPLETE",
      "node_count": 1,
      "image_id": "fedora-coreos",
      "status_reason": "Stack CREATE completed successfully",
      "flavor_id": "g1-standard-1-2"
    }
  ],
  "updated_at": "2020-04-20T14:27:44+00:00",
  "created_at": "2020-04-20T08:32:33+00:00",
  "cluster_template_id": "2c884df0-a312-4950-a6ec-9405803affc9",
  "coe_version": "v1.17.4",
  "name": "test",
  "create_timeout": 7200,
  "uuid": "5e09faed-e742-404f-8a75-0ea5eb3c435f",
  "labels": {
	"gcloud_project_id": "1",
	"gcloud_region_id": "2",
	"fixed_network_cidr": "192.168.0.0/18",
	"calico_ipv4pool": "192.168.64.0/18",
	"service_cluster_ip_range": "192.168.128.0/18"
  },
  "discovery_url": "https://discovery.etcd.io/6fba601fce5c2b84eebd7c472ab36650",
  "health_status_reason": {
    "api": "ok"
  },
  "stack_id": "2f0d5d97-fb3c-4218-9201-34f804299510",
  "floating_ip_enabled": true,
  "docker_volume_size": 10,
  "master_flavor_id": "g1-standard-1-2",
  "node_addresses": [
    "192.168.0.5"
  ],
  "fixed_network": "46beed39-38e6-4743-90b5-30fefd6173d2",
  "container_version": "1.12.6",
  "api_address": "https://172.24.4.3:6443",
  "node_count": 1,
  "status": "UPDATE_COMPLETE",
  "user_id": "3ed399cd-ecd1-4403-b3e7-e0029e4f694f",
  "keypair": "keypair",
  "master_addresses": [
    "192.168.0.11"
  ],
  "master_count": 1,
  "status_reason": "Stack CREATE completed successfully",
  "flavor_id": "g1-standard-1-2",
  "version": "1.17"	
}
`

var CreateRequest = fmt.Sprintf(`
{
  "auto_healing_enabled": false,
  "fixed_network": "%s",
  "fixed_subnet": "%s",
  "keypair": "keypair",
  "master_count": 1,
  "name": "%s",
  "pools": [
	{
	  "docker_volume_size": 10,
	  "docker_volume_type": "standard",
	  "flavor_id": "g1-standard-1-2",
	  "max_node_count": 2,
	  "min_node_count": 1,
	  "name": "test1",
	  "node_count": 1
	}
  ],
  "version": "%s"
}
`, fixedNetwork, fixedSubnet, Cluster1.Name, version)

var ResizeRequest = fmt.Sprintf(`
{
    "node_count": 2,
    "pool": "%s"
}
`, pool)

var UpgradeRequest = fmt.Sprintf(`
{
    "version": "1.17",
    "pool": "%s"
}
`, pool)

var VersionResponse = `
{
  "count": 2,
  "results": [
    "1.14",
    "1.17"
  ]
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

const UpgradeResponse = `
{
  "tasks": [
    "50f53a35-42ed-40c4-82b2-5a37fb3e00bc"
  ]
}
`

const ConfigStringResponse = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ca
    server: 10.0.0.1
  name: name
contexts:
- context:
    cluster: name
    user: admin
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: cert
    client-key-data: key
`

var (
	c                 = clusters.Config{Config: ConfigStringResponse}
	ConfigResponse, _ = json.Marshal(c)
	createdTimeString = "2020-04-20T08:32:33+00:00"
	updatedTimeString = "2020-04-20T14:27:44+00:00"
	createdTime, _    = time.Parse(time.RFC3339, createdTimeString)
	updatedTime, _    = time.Parse(time.RFC3339, updatedTimeString)
	pool              = uuid.NewV4().String()
	fixedSubnet       = uuid.NewV4().String()
	fixedNetwork      = uuid.NewV4().String()
	version           = "1.17"

	nodeAddresses   = []net.IP{net.ParseIP("192.168.0.5")}
	masterAddresses = []net.IP{net.ParseIP("192.168.0.11")}
	Config1         = clusters.Config{Config: ConfigStringResponse}
	labels          = map[string]string{
		"gcloud_project_id":        "1",
		"gcloud_region_id":         "2",
		"fixed_network_cidr":       "192.168.0.0/18",
		"calico_ipv4pool":          "192.168.64.0/18",
		"service_cluster_ip_range": "192.168.128.0/18",
	}

	listPool = pools.ClusterListPool{
		StackID:          "2f0d5d97-fb3c-4218-9201-34f804299510",
		Name:             "test1",
		MaxNodeCount:     5,
		MinNodeCount:     1,
		IsDefault:        true,
		FlavorID:         "g1-standard-1-2",
		UUID:             "908338b2-9217-4673-af0e-f0093139fbac",
		Status:           "CREATE_COMPLETE",
		Role:             "worker",
		ImageID:          "fedora-coreos",
		DockerVolumeSize: 10,
		DockerVolumeType: volumes.Standard,
		NodeCount:        1,
	}
	clusterList = &clusters.ClusterList{
		UUID:              "5e09faed-e742-404f-8a75-0ea5eb3c435f",
		Name:              "test",
		ClusterTemplateID: "2c884df0-a312-4950-a6ec-9405803affc9",
		KeyPair:           "keypair",
		NodeCount:         1,
		MasterCount:       1,
		DockerVolumeSize:  10,
		Labels:            labels,
		MasterFlavorID:    "g1-standard-1-2",
		FlavorID:          "g1-standard-1-2",
		CreateTimeout:     7200,
		StackID:           "2f0d5d97-fb3c-4218-9201-34f804299510",
		Status:            "UPDATE_COMPLETE",
		HealthStatus:      "HEALTHY",
		Version:           version,
	}
	apiAddress, _   = gcorecloud.ParseURL("https://172.24.4.3:6443")
	discoveryURL, _ = gcorecloud.ParseURL("https://discovery.etcd.io/6fba601fce5c2b84eebd7c472ab36650")
	ClusterList1    = clusters.ClusterListWithPool{
		Pools:       []pools.ClusterListPool{listPool},
		ClusterList: clusterList,
	}
	Cluster1 = clusters.ClusterWithPool{
		Cluster: &clusters.Cluster{
			StatusReason:     "Stack CREATE completed successfully",
			APIAddress:       apiAddress,
			CoeVersion:       "v1.17.4",
			ContainerVersion: "1.12.6",
			DiscoveryURL:     discoveryURL,
			HealthStatusReason: map[string]string{
				"api": "ok",
			},
			ProjectID:         "46beed39-38e6-4743-90b5-30fefd6173d2",
			UserID:            "3ed399cd-ecd1-4403-b3e7-e0029e4f694f",
			NodeAddresses:     nodeAddresses,
			MasterAddresses:   masterAddresses,
			FixedNetwork:      "46beed39-38e6-4743-90b5-30fefd6173d2",
			FixedSubnet:       "46beed39-38e6-4743-90b5-30fefd6173d2",
			FloatingIPEnabled: true,
			CreatedAt:         createdTime,
			UpdatedAt:         &updatedTime,
			Faults:            nil,
			ClusterList:       clusterList,
		},
		Pools: []pools.ClusterPool{{
			ClusterID:       "5e09faed-e742-404f-8a75-0ea5eb3c435f",
			ProjectID:       "46beed3938e6474390b530fefd6173d2",
			Labels:          labels,
			NodeAddresses:   nodeAddresses,
			StatusReason:    "Stack CREATE completed successfully",
			ClusterListPool: &listPool,
		},
		},
	}
	ClusterCertificate = clusters.ClusterCACertificate{
		ClusterUUID: ClusterList1.UUID,
		PEM:         "string",
	}
	ClusterSignedCertificate = clusters.ClusterSignCertificate{
		ClusterUUID: ClusterList1.UUID,
		PEM:         "string",
		CSR:         "string",
	}
	Tasks1 = tasks.TaskResults{
		Tasks: []tasks.TaskID{"50f53a35-42ed-40c4-82b2-5a37fb3e00bc"},
	}

	ExpectedClusterSlice = []clusters.ClusterListWithPool{ClusterList1}
)
