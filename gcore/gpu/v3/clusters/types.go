package clusters

import (
	"encoding/json"
	"fmt"
)

type IPFamilyType string
type VolumeType string
type ClusterStatusType string

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
	Suspended ClusterStatusType = "suspended"
	Error     ClusterStatusType = "error"
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

func (it *IPFamilyType) String() string {
	return string(*it)
}

func (it *IPFamilyType) List() []IPFamilyType {
	return []IPFamilyType{
		IPv6IPFamilyType,
		IPv4IPFamilyType,
		DualStackIPFamilyType,
	}
}

func (it *IPFamilyType) StringList() []string {
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

func (vt *VolumeType) List() []VolumeType {
	return []VolumeType{
		Standard,
		SsdHiIops,
		SsdLowLatency,
		SsdLocal,
		Cold,
		Ultra,
	}
}

func (vt *VolumeType) StringList() []string {
	var s []string
	for _, v := range vt.List() {
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
	err := v.IsValid()
	if err != nil {
		return err
	}
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
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ct = v
	return nil
}

// MarshalJSON - implements Marshaller interface
func (ct *ClusterStatusType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}
