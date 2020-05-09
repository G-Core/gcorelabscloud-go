package testing

import (
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/types"

	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := regions.UpdateOpts{}
	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	opts = regions.UpdateOpts{
		State: types.RegionStateDeleted,
	}
	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)
}
