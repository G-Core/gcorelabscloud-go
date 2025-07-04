package clusters

import (
	"encoding/json"
	"fmt"
)

type IPFamilyType string
type VolumeSource string
type VolumeType string
type ClusterStatusType string
type FloatingIPSource string

type InterfaceType string

type ClusterAction string

const (
	IPv4IPFamilyType      IPFamilyType = "ipv4"
	IPv6IPFamilyType      IPFamilyType = "ipv6"
	DualStackIPFamilyType IPFamilyType = "dual"

	Standard      VolumeType = "standard"
	SsdHiIops     VolumeType = "ssd_hiiops"
	SsdLowLatency VolumeType = "ssd_lowlatency"
	SsdLocal      VolumeType = "ssd_local"
	Cold          VolumeType = "cold"
	Ultra         VolumeType = "ultra"

	New       ClusterStatusType = "new"
	Active    ClusterStatusType = "active"
	Resizing  ClusterStatusType = "resizing"
	Suspended ClusterStatusType = "suspended"
	Deleting  ClusterStatusType = "deleting"
	Error     ClusterStatusType = "error"

	NewVolume VolumeSource = "new"
	Image     VolumeSource = "image"
	Snapshot  VolumeSource = "snapshot"

	External  InterfaceType = "external"
	Subnet    InterfaceType = "subnet"
	AnySubnet InterfaceType = "any_subnet"

	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"

	StartClusterAction      ClusterAction = "start"
	StopClusterAction       ClusterAction = "stop"
	HardRebootClusterAction ClusterAction = "hard_reboot"
	SoftRebootClusterAction ClusterAction = "soft_reboot"
	ResizeClusterAction     ClusterAction = "resize"
	UpdateTagsClusterAction ClusterAction = "update_tags"
)

func (it *IPFamilyType) IsValid() error {
	switch *it {
	case IPv6IPFamilyType, IPv4IPFamilyType, DualStackIPFamilyType:
		return nil
	}
	return fmt.Errorf("invalid IPFamilyType type: %v", it)
}

func (it *IPFamilyType) ValidOrNil() (*IPFamilyType, error) {
	if it.String() == "" {
		return nil, nil
	}
	err := it.IsValid()
	if err != nil {
		return it, err
	}
	return it, nil
}

func InterfaceTypeList() []InterfaceType {
	return []InterfaceType{
		External,
		Subnet,
		AnySubnet,
	}
}

func InterfaceTypeStringList() []string {
	var s []string
	for _, v := range InterfaceTypeList() {
		s = append(s, v.String())
	}
	return s
}

func (it *InterfaceType) String() string {
	return string(*it)
}

func (fis *FloatingIPSource) String() string {
	return string(*fis)
}

func FloatingIPSourceList() []FloatingIPSource {
	return []FloatingIPSource{
		NewFloatingIP,
		ExistingFloatingIP,
	}
}

func FloatingIPSourceStringList() []string {
	var s []string
	for _, v := range FloatingIPSourceList() {
		s = append(s, v.String())
	}
	return s
}

func (vs *VolumeSource) String() string {
	return string(*vs)
}

func VolumeSourcesList() []VolumeSource {
	return []VolumeSource{
		NewVolume,
		Image,
		Snapshot,
	}
}

func VolumeSourcesStringList() []string {
	var s []string
	for _, v := range VolumeSourcesList() {
		s = append(s, v.String())
	}
	return s
}

func (it *IPFamilyType) String() string {
	return string(*it)
}

func IPFamilyTypeList() []IPFamilyType {
	return []IPFamilyType{
		IPv6IPFamilyType,
		IPv4IPFamilyType,
		DualStackIPFamilyType,
	}
}

func IPFamilyTypeListStringList() []string {
	var s []string
	for _, v := range IPFamilyTypeList() {
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
	*it = v
	return nil
}

// MarshalJSON - implements Marshaller interface
func (it *IPFamilyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.String())
}

func (vt *VolumeType) IsValid() error {
	switch *vt {
	case Standard, SsdHiIops, SsdLowLatency, SsdLocal, Cold, Ultra:
		return nil
	}
	return fmt.Errorf("invalid VolumeType type: %v", vt)
}

func (vt *VolumeType) ValidOrNil() (*VolumeType, error) {
	if vt.String() == "" {
		return nil, nil
	}
	err := vt.IsValid()
	if err != nil {
		return vt, err
	}
	return vt, nil
}

func (vt *VolumeType) String() string {
	return string(*vt)
}

func VolumeTypesList() []VolumeType {
	return []VolumeType{
		Standard,
		SsdHiIops,
		SsdLowLatency,
		SsdLocal,
		Cold,
		Ultra,
	}
}

func VolumeTypesStringList() []string {
	var s []string
	for _, v := range VolumeTypesList() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (vt *VolumeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := VolumeType(s)
	*vt = v
	return nil
}

// MarshalJSON - implements Marshaller interface
func (vt *VolumeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(vt.String())
}

func (ct *ClusterStatusType) IsValid() error {
	switch *ct {
	case New, Active, Suspended, Error:
		return nil
	}
	return fmt.Errorf("invalid ClusterStatusType type: %v", ct)
}

func (ct *ClusterStatusType) ValidOrNil() (*ClusterStatusType, error) {
	if ct.String() == "" {
		return nil, nil
	}
	err := ct.IsValid()
	if err != nil {
		return ct, err
	}
	return ct, nil
}

func (ct *ClusterStatusType) String() string {
	return string(*ct)
}

func (ct *ClusterStatusType) List() []ClusterStatusType {
	return []ClusterStatusType{
		New,
		Active,
		Suspended,
		Error,
	}
}

func (ct *ClusterStatusType) StringList() []string {
	var s []string
	for _, v := range ct.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ct *ClusterStatusType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ClusterStatusType(s)
	*ct = v
	return nil
}

// MarshalJSON - implements Marshaller interface
func (ct *ClusterStatusType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}

func (ca ClusterAction) IsValid() error {
	switch ca {
	case StartClusterAction,
		StopClusterAction,
		SoftRebootClusterAction,
		HardRebootClusterAction,
		ResizeClusterAction,
		UpdateTagsClusterAction:
		return nil
	}
	return fmt.Errorf("invalid ClusterAction: %v", ca)
}

func (ca ClusterAction) ValidOrNil() (*ClusterAction, error) {
	if ca.String() == "" {
		return nil, nil
	}
	err := ca.IsValid()
	if err != nil {
		return &ca, err
	}
	return &ca, nil
}

func (ca ClusterAction) String() string {
	return string(ca)
}

func (ca ClusterAction) List() []ClusterAction {
	return []ClusterAction{
		StartClusterAction,
		StopClusterAction,
		SoftRebootClusterAction,
		HardRebootClusterAction,
		ResizeClusterAction,
		UpdateTagsClusterAction,
	}
}

func (ca ClusterAction) StringList() []string {
	var s []string
	for _, v := range ca.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ca *ClusterAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ClusterAction(s)
	*ca = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ca *ClusterAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}
