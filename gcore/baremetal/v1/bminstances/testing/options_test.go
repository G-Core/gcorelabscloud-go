package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bminstances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
)

func TestValidateCreateInstanceInterfaceOpts(t *testing.T) {
	opts := bminstances.CreateOpts{
		Flavor:        "",
		Names:         []string{"name"},
		NameTemplates: nil,
		ImageID:       "28bfe198-a003-4283-8dca-ab5da4a71b62",
		Interfaces: []bminstances.InterfaceOpts{{
			Type:      types.SubnetInterfaceType,
			NetworkID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
			FloatingIP: &bminstances.CreateNewInterfaceFloatingIPOpts{
				Source: types.NewFloatingIP,
			},
		}},
		Keypair:  "",
		Password: "",
		Username: "",
		UserData: "",
	}

	err := opts.Validate()
	require.Error(t, err)
	errTranslated := gcorecloud.TranslateValidationError(err)
	require.Error(t, errTranslated)
	require.Contains(t, errTranslated.Error(), "SubnetID")

	opts = bminstances.CreateOpts{
		Flavor:        "",
		Names:         []string{"name"},
		NameTemplates: nil,
		ImageID:       "28bfe198-a003-4283-8dca-ab5da4a71b62",
		Interfaces: []bminstances.InterfaceOpts{{
			Type: types.ExternalInterfaceType,
		}},
		Keypair:  "",
		Password: "",
		Username: "",
		UserData: "",
	}

	err = opts.Validate()
	require.NoError(t, err)

	opts = bminstances.CreateOpts{
		Flavor:        "",
		Names:         []string{"name"},
		NameTemplates: nil,
		AppTemplateID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
		Interfaces: []bminstances.InterfaceOpts{{
			Type:      types.SubnetInterfaceType,
			NetworkID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
			SubnetID:  "28bfe198-a003-4283-8dca-ab5da4a71b62",
			FloatingIP: &bminstances.CreateNewInterfaceFloatingIPOpts{
				Source: types.NewFloatingIP,
			},
		}},
		Keypair:  "",
		Password: "",
		Username: "",
		UserData: "",
	}

	err = opts.Validate()
	require.NoError(t, err)
}
