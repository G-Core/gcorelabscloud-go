package testing

import (
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/limit/v1/limits"

	"github.com/stretchr/testify/require"
)

func TestUpdateOptsToMap(t *testing.T) {
	opts := limits.NewUpdateOpts()
	_, err := opts.ToLimitUpdateMap()
	require.Error(t, err)
	opts = limits.NewUpdateOpts()
	opts.ProjectCountLimit = 0
	opts.VMCountLimit = 0
	opts.CPUCountLimit = 0
	opts.RAMLimit = -2
	_, err = opts.ToLimitUpdateMap()
	require.Error(t, err)
	opts = limits.NewUpdateOpts()
	opts.ProjectCountLimit = 0
	m, err := opts.ToLimitUpdateMap()
	require.NoError(t, err)
	require.Len(t, m, 1)
}

func TestCreateOptsToMap(t *testing.T) {
	opts := limits.NewCreateOpts("test")
	_, err := opts.ToLimitCreateMap()
	require.Error(t, err)
	opts = limits.NewCreateOpts("test")
	opts.RequestedQuotas.RAMLimit = 1
	m, err := opts.ToLimitCreateMap()
	require.NoError(t, err)
	require.Len(t, m["requested_quotas"], 1)
	opts = limits.NewCreateOpts("test")
	opts.RequestedQuotas.RAMLimit = 1
	opts.RequestedQuotas.CPUCountLimit = -2
	_, err = opts.ToLimitCreateMap()
	require.Error(t, err)
}
