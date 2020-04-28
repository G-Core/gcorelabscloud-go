package testing

import (
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"

	"github.com/G-Core/gcorelabscloud-go"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"

	"github.com/stretchr/testify/require"
)

func TestListQueryNullParams(t *testing.T) {
	opts := instances.ListOpts{
		ExcludeSecGroup:   "",
		AvailableFloating: false,
	}
	res, err := opts.ToInstanceListQuery()
	require.NoError(t, err)
	require.Equal(t, "", res)
}

func TestListQueryBoolParam(t *testing.T) {
	opts := instances.ListOpts{
		ExcludeSecGroup:   "",
		AvailableFloating: true,
	}
	res, err := opts.ToInstanceListQuery()
	require.NoError(t, err)
	require.Contains(t, res, "available_floating")
}

func TestValidateCreateVolumeBlankVolumeIDOpts(t *testing.T) {
	opts := instances.CreateVolumeOpts{
		Source:     types.ExistingVolume,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "",
	}

	err := opts.Validate()
	require.Error(t, err)
	errTranslated := gcorecloud.TranslateValidationError(err)
	require.Error(t, errTranslated)
	require.Contains(t, errTranslated.Error(), "VolumeID")

	opts = instances.CreateVolumeOpts{
		Source:     types.ExistingVolume,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
	}

	err = opts.Validate()
	require.NoError(t, err)

}

func TestValidateCreateVolumeBlankSnapshotIDOpts(t *testing.T) {

	opts := instances.CreateVolumeOpts{
		Source:     types.Snapshot,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "",
	}

	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "Snapshot")

	opts = instances.CreateVolumeOpts{
		Source:     types.Snapshot,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "2bf3a5d7-9072-40aa-8ac0-a64e39427a2c",
		VolumeID:   "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)

}

func TestValidateCreateVolumeOpts(t *testing.T) {

	opts := instances.CreateVolumeOpts{
		Source:     types.NewVolume,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "",
	}

	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "Size")

	opts = instances.CreateVolumeOpts{
		Source:     types.NewVolume,
		BootIndex:  0,
		Size:       1,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)

	opts = instances.CreateVolumeOpts{
		Source:     types.NewVolume,
		BootIndex:  0,
		Size:       1,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "28bfe198-a003-4283-8dca-ab5da4a71b62",
		SnapshotID: "",
		VolumeID:   "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "ImageID")

	opts = instances.CreateVolumeOpts{
		Source:     types.ExistingVolume,
		BootIndex:  0,
		Size:       1,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "28bfe198-a003-4283-8dca-ab5da4a71b62",
		SnapshotID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
		VolumeID:   "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "ImageID")
	require.Contains(t, err.Error(), "VolumeID")
	require.Contains(t, err.Error(), "SnapshotID")
	require.Contains(t, err.Error(), "Size")

	opts = instances.CreateVolumeOpts{
		Source:     "cc",
		BootIndex:  0,
		Size:       1,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "28bfe198-a003-4283-8dca-ab5da4a71b62",
		SnapshotID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
		VolumeID:   "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "Source")

}

func TestValidateCreateInstanceBlankSnapshotIDOpts(t *testing.T) {

	volumeOpts := []instances.CreateVolumeOpts{{
		Source:     types.Snapshot,
		BootIndex:  0,
		Size:       0,
		TypeName:   volumes.Standard,
		Name:       "name",
		ImageID:    "",
		SnapshotID: "",
		VolumeID:   "",
	}}

	opts := instances.CreateOpts{
		Flavor:        "",
		Names:         []string{"name"},
		NameTemplates: nil,
		Volumes:       volumeOpts,
		Interfaces: []instances.CreateInterfaceOpts{{
			Type:      types.SubnetInterfaceType,
			NetworkID: "28bfe198-a003-4283-8dca-ab5da4a71b62",
			SubnetID:  "28bfe198-a003-4283-8dca-ab5da4a71b62",
			FloatingIP: &instances.CreateNewInterfaceFloatingIPOpts{
				Source:             types.NewFloatingIP,
				ExistingFloatingID: "",
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
	require.Contains(t, errTranslated.Error(), "Snapshot")

}

func TestValidateDeleteOpts(t *testing.T) {

	opts := instances.DeleteOpts{
		Volumes:         []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
		DeleteFloatings: true,
		FloatingIPs:     []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
	}

	err := opts.Validate()
	require.Error(t, err)
	errTranslated := gcorecloud.TranslateValidationError(err)
	require.Error(t, errTranslated)
	require.Contains(t, errTranslated.Error(), "DeleteFloatings")

	opts = instances.DeleteOpts{
		Volumes:         []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
		DeleteFloatings: false,
		FloatingIPs:     []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
	}

	err = opts.Validate()
	require.NoError(t, err)

	opts = instances.DeleteOpts{
		Volumes:         []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
		DeleteFloatings: false,
		FloatingIPs:     nil,
	}

	err = opts.Validate()
	require.NoError(t, err)

	opts = instances.DeleteOpts{
		Volumes:         []string{"28bfe198-a003-4283-8dca-ab5da4a71b62"},
		DeleteFloatings: true,
		FloatingIPs:     nil,
	}

	err = opts.Validate()
	require.NoError(t, err)

}

func TestDeleteOpts(t *testing.T) {

	opts := instances.DeleteOpts{
		Volumes:         nil,
		DeleteFloatings: true,
		FloatingIPs:     nil,
	}
	query, err := opts.ToInstanceDeleteQuery()
	require.NoError(t, err)
	require.NotEmpty(t, query)

	opts = instances.DeleteOpts{
		Volumes:         []string{"28bfe198-a003-4283-8dca-ab5da4a71b62", "29bfe198-a003-4283-8dca-ab5da4a71b62"},
		DeleteFloatings: false,
		FloatingIPs:     nil,
	}
	query, err = opts.ToInstanceDeleteQuery()
	require.NoError(t, err)
	require.Equal(t, "?volumes=28bfe198-a003-4283-8dca-ab5da4a71b62%2C29bfe198-a003-4283-8dca-ab5da4a71b62", query)

}
