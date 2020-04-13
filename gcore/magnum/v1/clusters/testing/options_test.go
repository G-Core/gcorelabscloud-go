package testing

import (
	"testing"

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
