package testing

import (
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"

	"github.com/stretchr/testify/require"
)

func TestUpdateOpts(t *testing.T) {

	options := pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err := options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")

	options = pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)

	options = pools.UpdateOpts{}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "Name")

}

func TestCreateOpts(t *testing.T) {

	options := pools.CreateOpts{
		Name:             "",
		FlavorID:         "",
		NodeCount:        0,
		DockerVolumeSize: 0,
		MinNodeCount:     5,
		MaxNodeCount:     3,
	}

	_, err := options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "FlavorID")
	require.Contains(t, err.Error(), "NodeCount")

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        5,
		DockerVolumeSize: 10,
		MinNodeCount:     4,
		MaxNodeCount:     3,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        5,
		DockerVolumeSize: 10,
		MinNodeCount:     6,
		MaxNodeCount:     8,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "NodeCount")

}

func TestDecodePoolTask(t *testing.T) {
	taskID := "732851e1-f792-4194-b966-4cbfa5f30093"
	rs := map[string]interface{}{"k8s_pools": []string{taskID}}
	taskInfo := tasks.Task{
		CreatedResources: &rs,
	}
	var result pools.PoolTaskResult
	err := gcorecloud.NativeMapToStruct(taskInfo.CreatedResources, &result)
	require.NoError(t, err)
	require.Equal(t, taskID, result.K8sPools[0])
}
