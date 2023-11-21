package types

import (
	"encoding/json"
	"fmt"
)

type AddressType string
type VolumeSource string
type FloatingIPSource string
type InterfaceType string
type IPFamilyType string
type MetricsTimeUnit string

const (
	AddressTypeFixed    AddressType = "fixed"
	AddressTypeFloating AddressType = "floating"

	NewVolume      VolumeSource = "new-volume"
	ExistingVolume VolumeSource = "existing-volume"
	Image          VolumeSource = "image"
	Snapshot       VolumeSource = "snapshot"

	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"

	SubnetInterfaceType    InterfaceType = "subnet"
	AnySubnetInterfaceType InterfaceType = "any_subnet"
	ExternalInterfaceType  InterfaceType = "external"
	ReservedFixedIpType    InterfaceType = "reserved_fixed_ip"

	IPv4IPFamilyType      IPFamilyType = "ipv4"
	IPv6IPFamilyType      IPFamilyType = "ipv6"
	DualStackIPFamilyType IPFamilyType = "dual"

	HourMetricsTimeUnit MetricsTimeUnit = "hour"
	DayMetricsTimeUnit  MetricsTimeUnit = "day"
)

func (vs VolumeSource) IsValid() error {
	switch vs {
	case NewVolume, Image, Snapshot, ExistingVolume:
		return nil
	}
	return fmt.Errorf("invalid VolumeSource type: %v", vs)
}

func (vs VolumeSource) ValidOrNil() (*VolumeSource, error) {
	if vs.String() == "" {
		return nil, nil
	}
	err := vs.IsValid()
	if err != nil {
		return &vs, err
	}
	return &vs, nil
}

func (vs VolumeSource) String() string {
	return string(vs)
}

func (vs VolumeSource) List() []VolumeSource {
	return []VolumeSource{NewVolume, Image, Snapshot, ExistingVolume}
}

func (vs VolumeSource) Bootable() bool {
	switch vs {
	case Image, Snapshot:
		return true
	}
	return false
}

func (vs VolumeSource) StringList() []string {
	var s []string
	for _, v := range vs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for VolumeSource
func (vs *VolumeSource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := VolumeSource(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*vs = v
	return nil
}

// MarshalJSON - implements Marshaler interface for VolumeSource
func (vs *VolumeSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(vs.String())
}

func (fip FloatingIPSource) IsValid() error {
	switch fip {
	case NewFloatingIP, ExistingFloatingIP:
		return nil
	}
	return fmt.Errorf("invalid FloatingIPSource type: %v", fip)
}

func (fip FloatingIPSource) ValidOrNil() (*FloatingIPSource, error) {
	if fip.String() == "" {
		return nil, nil
	}
	err := fip.IsValid()
	if err != nil {
		return &fip, err
	}
	return &fip, nil
}

func (fip FloatingIPSource) String() string {
	return string(fip)
}

func (fip FloatingIPSource) List() []FloatingIPSource {
	return []FloatingIPSource{
		NewFloatingIP,
		ExistingFloatingIP,
	}
}

func (fip FloatingIPSource) StringList() []string {
	var s []string
	for _, v := range fip.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (fip *FloatingIPSource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := FloatingIPSource(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*fip = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (fip *FloatingIPSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(fip.String())
}

func (at AddressType) IsValid() error {
	switch at {
	case AddressTypeFixed, AddressTypeFloating:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", at)
}

func (at AddressType) ValidOrNil() (*AddressType, error) {
	if at.String() == "" {
		return nil, nil
	}
	err := at.IsValid()
	if err != nil {
		return &at, err
	}
	return &at, nil
}

func (at AddressType) String() string {
	return string(at)
}

func (at AddressType) List() []AddressType {
	return []AddressType{
		AddressTypeFixed,
		AddressTypeFloating,
	}
}

func (at AddressType) StringList() []string {
	var s []string
	for _, v := range at.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (at *AddressType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := AddressType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*at = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (at *AddressType) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.String())
}

func (it InterfaceType) IsValid() error {
	switch it {
	case ExternalInterfaceType, SubnetInterfaceType, AnySubnetInterfaceType, ReservedFixedIpType:
		return nil
	}
	return fmt.Errorf("invalid InterfaceType type: %v", it)
}

func (it InterfaceType) ValidOrNil() (*InterfaceType, error) {
	if it.String() == "" {
		return nil, nil
	}
	err := it.IsValid()
	if err != nil {
		return &it, err
	}
	return &it, nil
}

func (it InterfaceType) String() string {
	return string(it)
}

func (it InterfaceType) List() []InterfaceType {
	return []InterfaceType{
		ExternalInterfaceType,
		SubnetInterfaceType,
		AnySubnetInterfaceType,
		ReservedFixedIpType,
	}
}

func (it InterfaceType) StringList() []string {
	var s []string
	for _, v := range it.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (it *InterfaceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := InterfaceType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*it = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (it *InterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.String())
}

func (u MetricsTimeUnit) IsValid() error {
	switch u {
	case HourMetricsTimeUnit, DayMetricsTimeUnit:
		return nil
	}
	return fmt.Errorf("invalid MetricsTypeUnit type: %v", u)
}

func (u MetricsTimeUnit) ValidOrNil() (*MetricsTimeUnit, error) {
	if u.String() == "" {
		return nil, nil
	}
	err := u.IsValid()
	if err != nil {
		return &u, err
	}
	return &u, nil
}

func (u MetricsTimeUnit) String() string {
	return string(u)
}

func (u MetricsTimeUnit) List() []MetricsTimeUnit {
	return []MetricsTimeUnit{
		HourMetricsTimeUnit,
		DayMetricsTimeUnit,
	}
}

func (u MetricsTimeUnit) StringList() []string {
	var s []string
	for _, v := range u.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (u *MetricsTimeUnit) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := MetricsTimeUnit(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*u = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (u *MetricsTimeUnit) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (it IPFamilyType) IsValid() error {
	switch it {
	case IPv6IPFamilyType, IPv4IPFamilyType, DualStackIPFamilyType:
		return nil
	}
	return fmt.Errorf("invalid IPFamilyType type: %v", it)
}

func (it IPFamilyType) ValidOrNil() (*IPFamilyType, error) {
	if it.String() == "" {
		return nil, nil
	}
	err := it.IsValid()
	if err != nil {
		return &it, err
	}
	return &it, nil
}

func (it IPFamilyType) String() string {
	return string(it)
}

func (it IPFamilyType) List() []IPFamilyType {
	return []IPFamilyType{
		IPv6IPFamilyType,
		IPv4IPFamilyType,
		DualStackIPFamilyType,
	}
}

func (it IPFamilyType) StringList() []string {
	var s []string
	for _, v := range it.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (it *IPFamilyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := IPFamilyType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*it = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (it *IPFamilyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.String())
}
