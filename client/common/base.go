package common

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/gcore"

	"github.com/urfave/cli/v2"
)

func buildTokenClient(c *cli.Context, endpointName, endpointType string, version string) (*gcorecloud.ServiceClient, error) {
	settings, err := gcore.NewGCloudTokenAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	accessToken := c.String("access")
	if accessToken != "" {
		settings.AccessToken = accessToken
	}

	refreshToken := c.String("refresh")
	if refreshToken != "" {
		settings.RefreshToken = refreshToken
	}

	if version == "" {
		version = c.String("api-version")
	}
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")
	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToTokenOptions()
	eo := settings.ToEndpointOptions()
	return gcore.TokenClientServiceWithDebug(options, eo, settings.Debug)
}

func buildPlatformClient(c *cli.Context, endpointName, endpointType string, version string) (*gcorecloud.ServiceClient, error) {
	settings, err := gcore.NewGCloudPlatformAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	username := c.String("username")
	if username != "" {
		settings.Username = username
	}

	password := c.String("password")
	if password != "" {
		settings.Password = password
	}

	if version == "" {
		version = c.String("api-version")
	}
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")

	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToAuthOptions()
	eo := settings.ToEndpointOptions()
	return gcore.AuthClientServiceWithDebug(options, eo, settings.Debug)
}

func BuildClient(c *cli.Context, endpointName, version string) (*gcorecloud.ServiceClient, error) {
	clientType := flags.ClientType
	if clientType == "" {
		clientType = c.String("client-type")
	}
	if clientType == "token" {
		return buildTokenClient(c, endpointName, "", version)
	}
	return buildPlatformClient(c, endpointName, "", version)
}
