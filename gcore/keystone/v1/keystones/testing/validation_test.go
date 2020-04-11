package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/keystone/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/keystone/v1/keystones"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := keystones.UpdateOpts{}
	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	opts = keystones.UpdateOpts{
		State: types.KeystoneStateDeleted,
	}
	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)
}
