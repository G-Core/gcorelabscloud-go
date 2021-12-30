package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v1/limits"

	"github.com/stretchr/testify/require"
)

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
