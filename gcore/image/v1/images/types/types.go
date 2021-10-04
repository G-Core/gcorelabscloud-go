package types

import (
	"encoding/json"
	"fmt"
)

type Visibility string

// HwMachineType virtual chipset type.
type HwMachineType string

// SshKeyType whether the image supports SSH key or not
type SshKeyType string

// OSType the operating system installed on the image.
type OSType string

//HwFirmwareType specifies the type of firmware with which to boot the guest.
type HwFirmwareType string

//ImageSourceType
type ImageSourceType string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPublic  Visibility = "public"

	HwMachineI440 HwMachineType = "i440"
	HwMachineQ35  HwMachineType = "q35"

	SshKeyAllow    SshKeyType = "allow"
	SshKeyDeny     SshKeyType = "deny"
	SshKeyRequired SshKeyType = "required"

	OsLinux   OSType = "linux"
	OsWindows OSType = "windows"

	HwFirmwareBIOS HwFirmwareType = "bios"
	HwFirmwareUEFI HwFirmwareType = "uefi"

	ImageSourceVolume ImageSourceType = "volume"
)

func (v Visibility) IsValid() error {
	switch v {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return nil
	}
	return fmt.Errorf("invalid Visibility type: %v", v)
}

func (v Visibility) ValidOrNil() (*Visibility, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v Visibility) String() string {
	return string(v)
}

func (v Visibility) List() []Visibility {
	return []Visibility{VisibilityPrivate, VisibilityShared, VisibilityPublic}
}

func (v Visibility) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for Visibility
func (v *Visibility) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := Visibility(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for Visibility
func (v *Visibility) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v HwMachineType) IsValid() error {
	switch v {
	case HwMachineI440, HwMachineQ35:
		return nil
	}
	return fmt.Errorf("invalid HwMachineType type: %v", v)
}

func (v HwMachineType) ValidOrNil() (*HwMachineType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v HwMachineType) String() string {
	return string(v)
}

func (v HwMachineType) List() []HwMachineType {
	return []HwMachineType{HwMachineI440, HwMachineQ35}
}

func (v HwMachineType) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for HwMachineType
func (v *HwMachineType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := HwMachineType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for HwMachineType
func (v *HwMachineType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v SshKeyType) IsValid() error {
	switch v {
	case SshKeyAllow, SshKeyDeny, SshKeyRequired:
		return nil
	}
	return fmt.Errorf("invalid SshKeyType type: %v", v)
}

func (v SshKeyType) ValidOrNil() (*SshKeyType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v SshKeyType) String() string {
	return string(v)
}

func (v SshKeyType) List() []SshKeyType {
	return []SshKeyType{SshKeyAllow, SshKeyDeny, SshKeyRequired}
}

func (v SshKeyType) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for SshKeyType
func (v *SshKeyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := SshKeyType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for SshKeyType
func (v *SshKeyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v OSType) IsValid() error {
	switch v {
	case OsLinux, OsWindows:
		return nil
	}
	return fmt.Errorf("invalid OSType type: %v", v)
}

func (v OSType) ValidOrNil() (*OSType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v OSType) String() string {
	return string(v)
}

func (v OSType) List() []OSType {
	return []OSType{OsLinux, OsWindows}
}

func (v OSType) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for OSType
func (v *OSType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := OSType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for OSType
func (v *OSType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v HwFirmwareType) IsValid() error {
	switch v {
	case HwFirmwareBIOS, HwFirmwareUEFI:
		return nil
	}
	return fmt.Errorf("invalid HwFirmwareType type: %v", v)
}

func (v HwFirmwareType) ValidOrNil() (*HwFirmwareType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v HwFirmwareType) String() string {
	return string(v)
}

func (v HwFirmwareType) List() []HwFirmwareType {
	return []HwFirmwareType{HwFirmwareBIOS, HwFirmwareUEFI}
}

func (v HwFirmwareType) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for HwFirmwareType
func (v *HwFirmwareType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := HwFirmwareType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for HwFirmwareType
func (v *HwFirmwareType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v ImageSourceType) IsValid() error {
	switch v {
	case ImageSourceVolume:
		return nil
	}
	return fmt.Errorf("invalid ImageSourceType type: %v", v)
}

func (v ImageSourceType) ValidOrNil() (*ImageSourceType, error) {
	if v.String() == "" {
		return nil, nil
	}
	err := v.IsValid()
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (v ImageSourceType) String() string {
	return string(v)
}

func (v ImageSourceType) List() []ImageSourceType {
	return []ImageSourceType{ImageSourceVolume}
}

func (v ImageSourceType) StringList() []string {
	var s []string
	for _, v := range v.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for ImageSourceType
func (v *ImageSourceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	vt := ImageSourceType(s)
	err := vt.IsValid()
	if err != nil {
		return err
	}
	*v = vt
	return nil
}

// MarshalJSON - implements Marshaler interface for ImageSourceType
func (v *ImageSourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}
