package testing

import (
	"strings"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"

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

func TestListQueryWrongLimitParam(t *testing.T) {
	opts := instances.ListOpts{
		ExcludeSecGroup:   "",
		AvailableFloating: true,
		Limit:             -2,
	}
	_, err := opts.ToInstanceListQuery()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Limit must be greater")
}

func TestListQueryFlavorID(t *testing.T) {
	opts := instances.ListOpts{
		AvailableFloating: true,
		FlavorID:          "g1-standard-1-2",
	}
	res, err := opts.ToInstanceListQuery()
	require.NoError(t, err)
	require.Equal(t, res, "?available_floating=true&flavor_id=g1-standard-1-2")
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
		Source:        types.NewVolume,
		BootIndex:     0,
		Size:          0,
		TypeName:      volumes.Standard,
		Name:          "name",
		ImageID:       "",
		SnapshotID:    "",
		VolumeID:      "",
		AttachmentTag: "",
	}

	err := gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "Size")

	opts = instances.CreateVolumeOpts{
		Source:        types.NewVolume,
		BootIndex:     0,
		Size:          1,
		TypeName:      volumes.Standard,
		Name:          "name",
		ImageID:       "",
		SnapshotID:    "",
		VolumeID:      "",
		AttachmentTag: "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)

	opts = instances.CreateVolumeOpts{
		Source:        types.NewVolume,
		BootIndex:     0,
		Size:          1,
		TypeName:      volumes.Standard,
		Name:          "name",
		ImageID:       "28bfe198-a003-4283-8dca-ab5da4a71b62",
		SnapshotID:    "",
		VolumeID:      "",
		AttachmentTag: "",
	}

	err = gcorecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	require.Contains(t, err.Error(), "ImageID")

	opts = instances.CreateVolumeOpts{
		Source:        types.ExistingVolume,
		BootIndex:     0,
		Size:          1,
		TypeName:      volumes.Standard,
		Name:          "name",
		ImageID:       "28bfe198-a003-4283-8dca-ab5da4a71b62",
		SnapshotID:    "28bfe198-a003-4283-8dca-ab5da4a71b62",
		VolumeID:      "",
		AttachmentTag: "",
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
		Interfaces: []instances.InterfaceOpts{{
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

func TestMetadataOpts(t *testing.T) {
	opts := instances.MetadataOpts{
		Key:   "test",
		Value: strings.Repeat("test", 100),
	}
	err := opts.Validate()
	require.Error(t, err)
}

func TestMetadataSetOpts(t *testing.T) {
	opts := instances.MetadataSetOpts{
		Metadata: []instances.MetadataOpts{
			{
				Key:   "test",
				Value: strings.Repeat("test", 100),
			},
			{
				Key:   "test",
				Value: strings.Repeat("test", 100),
			},
		},
	}
	err := opts.Validate()
	require.Error(t, err)

	blankOpts := instances.MetadataSetOpts{}
	err = blankOpts.Validate()
	require.Error(t, err)

	_, err = opts.ToMetadataMap()
	require.Error(t, err)

	opts = instances.MetadataSetOpts{
		Metadata: []instances.MetadataOpts{
			{
				Key:   "test1",
				Value: "test1",
			},
			{
				Key:   "test2",
				Value: "test2",
			},
		},
	}

	data, err := opts.ToMetadataMap()
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"test1": "test1", "test2": "test2"}, data)
}

func TestInterfaceOpts(t *testing.T) {
	opts := instances.InterfaceOpts{
		Type:      types.SubnetInterfaceType,
		SubnetID:  "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		NetworkID: "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
	}
	err := opts.Validate()
	require.NoError(t, err)

	opts = instances.InterfaceOpts{
		Type: types.SubnetInterfaceType,
	}
	err = opts.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "SubnetID")

	opts = instances.InterfaceOpts{
		Type: types.AnySubnetInterfaceType,
	}
	err = opts.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NetworkID")

	opts = instances.InterfaceOpts{
		Type: types.ReservedFixedIpType,
	}
	err = opts.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "PortID")

	opts = instances.InterfaceOpts{
		Type:      types.ReservedFixedIpType,
		SubnetID:  "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		NetworkID: "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		PortID:    "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
	}
	err = opts.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NetworkID SubnetID FloatingIP")

	opts = instances.InterfaceOpts{
		Type:      types.SubnetInterfaceType,
		SubnetID:  "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		NetworkID: "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		PortID:    "9bc36cf6-407c-4a74-bc83-ce3aa3854c3d",
		IpAddress: "192.168.100.2",
	}
	err = opts.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Type NetworkID SubnetID FloatingIP")
}
