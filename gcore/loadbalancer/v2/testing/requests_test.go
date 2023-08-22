package testing

import (
	metadataV2Testing "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v2/metadata/testing"
	"testing"
)

func TestMetadataCreate(t *testing.T) {
	metadataV2Testing.BuildTestMetadataCreate("loadbalancers", LoadBalancer1.ID)(t)
}

func TestMetadataUpdate(t *testing.T) {
	metadataV2Testing.BuildTestMetadataUpdate("loadbalancers", LoadBalancer1.ID)(t)
}

func TestMetadataDelete(t *testing.T) {
	metadataV2Testing.BuildTestMetadataDelete("loadbalancers", LoadBalancer1.ID)(t)
}
