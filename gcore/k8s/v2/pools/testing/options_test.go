package testing

import (
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"

	"github.com/stretchr/testify/require"
)

func TestCreateOpts(t *testing.T) {
	options := pools.CreateOpts{
		Name:         "",
		FlavorID:     "",
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err := options.ToClusterPoolCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "FlavorID")
	require.Contains(t, err.Error(), "NodeCount")

	options = pools.CreateOpts{
		Name:         "name",
		FlavorID:     "flavor",
		MinNodeCount: 4,
		MaxNodeCount: 3,
	}

	_, err = options.ToClusterPoolCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")

	options = pools.CreateOpts{
		Name:         "name",
		FlavorID:     "flavor",
		MinNodeCount: 6,
		MaxNodeCount: 8,
	}

	_, err = options.ToClusterPoolCreateMap()
	require.NoError(t, err)
}

func TestUpdateOpts(t *testing.T) {
	options := pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err := options.ToClusterPoolUpdateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "MaxNodeCount")

	options = pools.UpdateOpts{}

	_, err = options.ToClusterPoolUpdateMap()
	require.NoError(t, err)

	options = pools.UpdateOpts{
		MinNodeCount: 2,
		MaxNodeCount: 2,
	}

	_, err = options.ToClusterPoolUpdateMap()
	require.NoError(t, err)
}

func TestResizeOpts(t *testing.T) {
	options := pools.ResizeOpts{
		NodeCount: 0,
	}

	_, err := options.ToClusterPoolResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodeCount")

	options = pools.ResizeOpts{
		NodeCount: 1,
	}

	_, err = options.ToClusterPoolResizeMap()
	require.NoError(t, err)
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
