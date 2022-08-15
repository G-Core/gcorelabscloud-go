package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	fake "github.com/G-Core/gcorelabscloud-go/testhelper/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func prepareListProfileTemplatesTestURLParams(regionID int) string {
	return fmt.Sprintf("/v1/ddos/profile-templates/%d", regionID)
}

func prepareListProfilesTestURLParams(projectID, regionID int) string {
	return fmt.Sprintf("/v1/ddos/profiles/%d/%d", projectID, regionID)
}

func prepareProfileTestURLParams(projectID, regionID, profileID int) string {
	return fmt.Sprintf("/v1/ddos/profiles/%d/%d/%d", projectID, regionID, profileID)
}

func prepareActivateProfileTestURLParams(projectID, regionID, profileID int) string {
	return fmt.Sprintf("/v1/ddos/profiles/%d/%d/%d/action", projectID, regionID, profileID)
}

func prepareListProfileTemplatesTestURL() string {
	return prepareListProfileTemplatesTestURLParams(fake.RegionID)
}

func prepareAccessibilityTestURL() string {
	return fmt.Sprintf("/v1/ddos/accessibility")
}

func prepareListProfilesTestURL() string {
	return prepareListProfilesTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareCreateProfileTestURL() string {
	return prepareListProfilesTestURLParams(fake.ProjectID, fake.RegionID)
}

func prepareUpdateProfileTestURL(profileID int) string {
	return prepareProfileTestURLParams(fake.ProjectID, fake.RegionID, profileID)
}

func prepareDeleteProfileTestURL(profileID int) string {
	return prepareProfileTestURLParams(fake.ProjectID, fake.RegionID, profileID)
}

func prepareRegionCoverageCheckTestURL() string {
	return fmt.Sprintf("/v1/ddos/region_coverage/%d", fake.RegionID)
}

func prepareActivateProfileTestURL(profileID int) string {
	return prepareActivateProfileTestURLParams(fake.ProjectID, fake.RegionID, profileID)
}

func TestAccessibility(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareAccessibilityTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, accessibilityResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos", "v1")
	ct, err := ddos.GetAccessibility(client).Extract()

	require.NoError(t, err)
	require.Equal(t, accessStatus, *ct)
}

func TestRegionCoverageCheck(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareRegionCoverageCheckTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := fmt.Fprint(w, regionCoverageResponse); err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos/region_coverage", "v1")
	coverage, err := ddos.CheckRegionCoverage(client).Extract()

	require.NoError(t, err)
	require.Equal(t, regionCoverage, *coverage)
}

func TestListAllProfileTemplates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListProfileTemplatesTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, profileTemplatesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos", "v1")
	templates, err := ddos.ListAllProfileTemplates(client)

	require.NoError(t, err)
	require.Equal(t, profileTemplates, templates[0])
}

func TestListAllProfiles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListProfilesTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, profilesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos/profiles", "v1")
	profiles, err := ddos.ListAllProfiles(client)

	require.NoError(t, err)
	require.Equal(t, profile, profiles[0])
}

func TestCreateProfile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareCreateProfileTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createProfileRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, tasksList)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos/profiles", "v1")
	opts := ddos.CreateProfileOpts{
		IPAddress:           "123.123.123.1",
		ProfileTemplate:     1,
		BaremetalInstanceID: "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		Fields: []ddos.ProfileField{
			{
				BaseField: 1,
				Value:     "string",
			},
		},
	}

	tasks, err := ddos.CreateProfile(client, opts).Extract()
	require.NoError(t, err)
	require.Equal(t, task, *tasks)
}

func TestDeleteProfile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareDeleteProfileTestURL(profile.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, tasksList)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("ddos/profiles", "v1")
	tasks, err := ddos.DeleteProfile(client, profile.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, task, *tasks)
}

func TestUpdateProfile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareUpdateProfileTestURL(profile.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, tasksList)
		if err != nil {
			log.Error(err)
		}
	})

	opts := ddos.UpdateProfileOpts{
		ProfileTemplate:     7,
		BaremetalInstanceID: "9f310fe7-baa2-47a3-b6a6-63c5d78becc2",
		IPAddress:           "123.123.123.1",
	}

	client := fake.ServiceTokenClient("ddos/profiles", "v1")
	tasks, err := ddos.UpdateProfile(client, profile.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, task, *tasks)
}

func TestActivateProfile(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareActivateProfileTestURL(profile.ID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, tasksList)
		if err != nil {
			log.Error(err)
		}
	})

	opts := ddos.ActivateProfileOpts{
		BGP:    true,
		Active: true,
	}

	client := fake.ServiceTokenClient("ddos/profiles", "v1")
	tasks, err := ddos.ActivateProfile(client, profile.ID, opts).Extract()

	require.NoError(t, err)
	require.Equal(t, task, *tasks)
}
