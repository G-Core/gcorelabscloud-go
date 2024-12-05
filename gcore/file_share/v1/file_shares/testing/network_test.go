package testing

import (
	"encoding/json"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	sfs "github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/stretchr/testify/require"
)

func TestFileShareSubnetNotUUIDValidate(t *testing.T) {
	opts := sfs.FileShareNetworkOpts{
		NetworkID: "9b17dd07-1281-4fe0-8c13-d80c5725e297",
		SubnetID:  "0:0:0:0:0:ffff:192.1.56.10",
	}
	err := gcorecloud.Validate.Struct(opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Field validation for 'SubnetID' failed on the 'uuid4'")
}

func TestFileShareSubnetOmittedOnEmptyValidate(t *testing.T) {
	opts := sfs.FileShareNetworkOpts{
		NetworkID: "9b17dd07-1281-4fe0-8c13-d80c5725e297",
	}
	err := gcorecloud.Validate.Struct(opts)
	require.NoError(t, err)
	s, err := json.Marshal(opts)
	require.NoError(t, err)
	require.JSONEq(t, fileShareNetworkConfigWithoutSubnet, string(s))
}
