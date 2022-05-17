package utils

import (
	"net/url"
	"regexp"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func BaseRootEndpoint(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// NormalizeURLPath removes duplicated slashes
func NormalizeURLPath(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	path := u.Path
	r := regexp.MustCompile(`//+`)
	u.Path = r.ReplaceAllLiteralString(path, "/")
	return gcorecloud.NormalizeURL(u.String()), nil
}
