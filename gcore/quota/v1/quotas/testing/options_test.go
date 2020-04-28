package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/quota/v1/quotas"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptsToMap(t *testing.T) {
	opts := quotas.UpdateOpts{Quota: quotas.NewQuota()}
	_, err := opts.ToQuotaUpdateMap()
	require.Error(t, err)
	opts = quotas.UpdateOpts{Quota: quotas.NewQuota()}
	opts.ProjectCountLimit = 0
	opts.VMCountLimit = 0
	opts.CPUCountLimit = 0
	opts.CPUCountUsage = -1
	m, err := opts.ToQuotaUpdateMap()
	require.NoError(t, err)
	require.Len(t, m, 3)

	opts = quotas.UpdateOpts{Quota: quotas.NewQuota()}
	opts.ProjectCountLimit = -2
	_, err = opts.ToQuotaUpdateMap()
	require.Error(t, err)

}

func TestReplaceOptsToMap(t *testing.T) {
	opts := quotas.ReplaceOpts{Quota: Quota1}
	_, err := opts.ToQuotaReplaceMap()
	require.NoError(t, err)
	opts = quotas.ReplaceOpts{Quota: quotas.NewQuota()}
	opts.ProjectCountLimit = 0
	opts.VMCountLimit = 0
	opts.CPUCountLimit = 0
	_, err = opts.ToQuotaReplaceMap()
	require.Error(t, err)

	opts = quotas.ReplaceOpts{Quota: quotas.NewQuota()}
	opts.ProjectCountLimit = -2
	_, err = opts.ToQuotaReplaceMap()
	require.Error(t, err)

}
