package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/project/v1/projects"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := projects.UpdateOpts{}
	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	opts = projects.UpdateOpts{
		Name: "test",
	}
	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)
}
