package testing

import (
	"github.com/G-Core/gcorelabscloud-go/client/limits/v2/limits"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestGlobalUpdateStructFromString(t *testing.T) {
	gl := &limits.GlobalLimitsFlag{}
	err := gl.Set("keypair_count_limit=1;project_count_limit=2")
	require.NoError(t, err)
	require.Equal(t, *gl, limits.GlobalLimitsFlag{KeypairCountLimit: 1, ProjectCountLimit: 2})
}

func TestRegionalUpdateStructFromString(t *testing.T) {
	gl := &limits.RegionLimitsFlag{}
	testParam := "region_id=1;secret_count_limit=1,region_id=2;cpu_count_limit=1"
	err := gl.Set(testParam)
	expectedFlagValue := limits.RegionLimitsFlag{1: {RegionID: 1, SecretCountLimit: 1}, 2: {RegionID: 2, CPUCountLimit: 1}}
	require.NoError(t, err)
	require.Equal(t, *gl, expectedFlagValue)
}

func TestGenericOpts(t *testing.T) {
	createCommand := limits.Commands.Subcommands[2]
	app := &cli.App{Writer: ioutil.Discard, Description: "test"}

	os.Args = []string{
		"limit", "create", "--description=desc test",
		"--global=keypair_count_limit=1;project_count_limit=2",
		"--regions=region_id=1;cpu_count_limit=2",
	}
	app.Commands = append(app.Commands, createCommand)

	err := app.Run(os.Args)
	require.NoError(t, err)
}
