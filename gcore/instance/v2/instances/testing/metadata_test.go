package testing

import (
	"testing"

	instancesV2 "github.com/G-Core/gcorelabscloud-go/gcore/instance/v2/instances"
	"github.com/stretchr/testify/require"
)

func TestMetadataItemEmptyParam(t *testing.T) {
	opts := instancesV2.MetadataItemOpts{
		Key: "",
	}
	_, err := opts.ToMetadataItemQuery()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Key is a required field")
}
