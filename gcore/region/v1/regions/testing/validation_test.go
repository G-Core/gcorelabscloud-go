package testing

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/types"

	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := regions.UpdateOpts{}
	err := gcorecloud.ValidateStruct(opts)
	require.Error(t, err)

	opts = regions.UpdateOpts{
		State: types.RegionStateDeleted,
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		State: types.RegionStateActive,
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		DisplayName: "test",
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		SpiceProxyURL: gcorecloud.MustParseURL("http://test.com"),
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		EndpointType: types.EndpointTypePublic,
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		ExternalNetworkID: uuid.NewV4().String(),
	}
	err = gcorecloud.ValidateStruct(opts)
	require.NoError(t, err)

}
