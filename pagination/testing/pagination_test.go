package testing

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/testhelper"
)

func createClient() *gcorecloud.ServiceClient {
	return &gcorecloud.ServiceClient{
		ProviderClient: &gcorecloud.ProviderClient{AccessTokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}
