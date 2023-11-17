package testing

import (
	"encoding/json"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/stretchr/testify/require"
)

func TestCreateProfileOpts(t *testing.T) {
	options := ddos.CreateProfileOpts{
		ProfileTemplate: 1,
		ResourceID:      "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		ResourceType:    ddos.ResourceTypeInstance,
		IPAddress:       "123.123.123.1",
		Fields: []ddos.ProfileField{
			{
				Value:     "string",
				BaseField: 1,
			},
		},
	}

	mp, err := options.ToProfileCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, createProfileRequest, string(s))
}

func TestCreateProfileOptsLegacy(t *testing.T) {
	options := ddos.CreateProfileOpts{
		ProfileTemplate:     1,
		BaremetalInstanceID: "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		IPAddress:           "123.123.123.1",
		Fields: []ddos.ProfileField{
			{
				Value:     "string",
				BaseField: 1,
			},
		},
	}

	mp, err := options.ToProfileCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, createProfileRequestLegacy, string(s))
}

func TestUpdateProfileOpts(t *testing.T) {
	options := ddos.UpdateProfileOpts{
		ResourceID:      "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		ResourceType:    ddos.ResourceTypeLoadBalancer,
		ProfileTemplate: 1,
		IPAddress:       "123.123.123.1",
		Fields: []ddos.ProfileField{
			{
				Value:     "string",
				BaseField: 1,
			},
		},
	}

	mp, err := options.ToProfileUpdateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, updateProfileRequest, string(s))
}

func TestUpdateProfileOptsLegacy(t *testing.T) {
	options := ddos.UpdateProfileOpts{
		BaremetalInstanceID: "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		ProfileTemplate:     1,
		IPAddress:           "123.123.123.1",
		Fields: []ddos.ProfileField{
			{
				Value:     "string",
				BaseField: 1,
			},
		},
	}

	mp, err := options.ToProfileUpdateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, updateProfileRequestLegacy, string(s))
}

func TestActivateProfileOpts(t *testing.T) {
	options := ddos.ActivateProfileOpts{
		BGP:    true,
		Active: true,
	}

	mp, err := options.ToActivateProfileMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, activateProfileRequest, string(s))
}
