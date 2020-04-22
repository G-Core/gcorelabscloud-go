package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/clusters"

	"github.com/stretchr/testify/require"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/clustertemplates"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/types"
)

func TestUpdateOpts(t *testing.T) {
	options := clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "/labels",
		Value: map[string]string{"one": "two"},
		Op:    types.ClusterUpdateOperationReplace,
	}}

	_, err := options.ToClusterTemplateUpdateMap()
	require.NoError(t, err)

	options = clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "/labels",
		Value: map[string]string{"one": "two"},
		Op:    "x",
	}}

	_, err = options.ToClusterTemplateUpdateMap()
	require.Error(t, err)

	options = clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "/labels",
		Value: nil,
		Op:    types.ClusterUpdateOperationAdd,
	}}

	_, err = options.ToClusterTemplateUpdateMap()
	require.Error(t, err)

	options = clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "/labels",
		Value: "",
		Op:    types.ClusterUpdateOperationAdd,
	}}

	_, err = options.ToClusterTemplateUpdateMap()
	require.Error(t, err)

	options = clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "labels",
		Value: "x",
		Op:    types.ClusterUpdateOperationAdd,
	}}

	_, err = options.ToClusterTemplateUpdateMap()
	require.Error(t, err)

	options = clustertemplates.UpdateOpts{clustertemplates.UpdateOptsElem{
		Path:  "/labels",
		Value: "x",
		Op:    types.ClusterUpdateOperationAdd,
	}}

	_, err = options.ToClusterTemplateUpdateMap()
	require.NoError(t, err)

}

func TestResizeOpts(t *testing.T) {
	options := clusters.ResizeOpts{
		NodeCount:     0,
		NodesToRemove: nil,
		NodeGroup:     "",
	}

	_, err := options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodeCount")

	options = clusters.ResizeOpts{
		NodeCount:     1,
		NodesToRemove: []string{"1"},
		NodeGroup:     "",
	}

	_, err = options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodeGroup")
	require.Contains(t, err.Error(), "NodesToRemove")

	options = clusters.ResizeOpts{
		NodeCount: 1,
		NodeGroup: "",
	}

	_, err = options.ToClusterResizeMap()
	require.NoError(t, err)

}
