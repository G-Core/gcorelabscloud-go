package stacks

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

type ErrInvalidEnvironment struct {
	gcorecloud.BaseError
	Section string
}

func (e ErrInvalidEnvironment) Error() string {
	return fmt.Sprintf("Environment has wrong section: %s", e.Section)
}

type ErrInvalidDataFormat struct {
	gcorecloud.BaseError
}

func (e ErrInvalidDataFormat) Error() string {
	return fmt.Sprintf("Data in neither json nor yaml format.")
}

type ErrInvalidTemplateFormatVersion struct {
	gcorecloud.BaseError
	Version string
}

func (e ErrInvalidTemplateFormatVersion) Error() string {
	return fmt.Sprintf("Template format version not found.")
}

type ErrTemplateRequired struct {
	gcorecloud.BaseError
}

func (e ErrTemplateRequired) Error() string {
	return fmt.Sprintf("Template required for this function.")
}
