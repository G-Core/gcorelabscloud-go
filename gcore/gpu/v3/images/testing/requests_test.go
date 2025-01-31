package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/images"
	th "github.com/G-Core/gcorelabscloud-go/testhelper"
	thclient "github.com/G-Core/gcorelabscloud-go/testhelper/client"
)

func TestUploadBaremetalImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(fmt.Sprintf("/v3/gpu/baremetal/%d/%d/images", thclient.ProjectID, thclient.RegionID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", thclient.TokenID))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{
			"url": "http://example.com/image.img",
			"name": "test-image",
			"ssh_key": "allow",
			"cow_format": true,
			"architecture": "x86_64",
			"os_distro": "ubuntu",
			"os_type": "linux",
			"os_version": "20.04",
			"hw_firmware_type": "bios",
			"metadata": {
				"key": "value"
			}
		}`)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, TestUploadBaremetalImageResponse)
		if err != nil {
			t.Error(err)
		}
	})

	sshKey := images.SshKeyAllow
	cowFormat := true
	arch := "x86_64"
	osType := images.OsTypeLinux
	hwType := images.HwFirmwareTypeBios

	opts := images.UploadBaremetalImageOpts{
		URL:            "http://example.com/image.img",
		Name:           "test-image",
		SshKey:         &sshKey,
		CowFormat:      &cowFormat,
		Architecture:   &arch,
		OsDistro:       stringPtr("ubuntu"),
		OsType:         &osType,
		OsVersion:      stringPtr("20.04"),
		HwFirmwareType: &hwType,
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	sc := &images.ServiceClient{
		ServiceClient: thclient.ServiceClient(),
	}
	taskResults, err := sc.UploadBaremetalImage(opts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, ExpectedTaskResults, taskResults)
}

func TestUploadVirtualImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(fmt.Sprintf("/v3/gpu/virtual/%d/%d/images", thclient.ProjectID, thclient.RegionID), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", thclient.TokenID))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{
			"url": "http://example.com/image.img",
			"name": "test-image",
			"ssh_key": "allow",
			"cow_format": true,
			"architecture": "x86_64",
			"os_distro": "ubuntu",
			"os_type": "linux",
			"os_version": "20.04",
			"hw_firmware_type": "bios",
			"metadata": {
				"key": "value"
			}
		}`)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, TestUploadVirtualImageResponse)
		if err != nil {
			t.Error(err)
		}
	})

	sshKey := images.SshKeyAllow
	cowFormat := true
	arch := "x86_64"
	osType := images.OsTypeLinux
	hwType := images.HwFirmwareTypeBios

	opts := images.UploadVirtualImageOpts{
		URL:            "http://example.com/image.img",
		Name:           "test-image",
		SshKey:         &sshKey,
		CowFormat:      &cowFormat,
		Architecture:   &arch,
		OsDistro:       stringPtr("ubuntu"),
		OsType:         &osType,
		OsVersion:      stringPtr("20.04"),
		HwFirmwareType: &hwType,
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	sc := &images.ServiceClient{
		ServiceClient: thclient.ServiceClient(),
	}
	taskResults, err := sc.UploadVirtualImage(opts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, ExpectedTaskResults, taskResults)
}

func stringPtr(s string) *string {
	return &s
}
