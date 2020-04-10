package types

import (
	"encoding/json"
	"fmt"
)

type RegionState string
type KeystoneState string
type EndpointType string

const (
	EndpointTypePublic                EndpointType  = "public"
	EndpointTypeInternal              EndpointType  = "internal"
	EndpointTypeAdmin                 EndpointType  = "admin"
	RegionStateActive                 RegionState   = "ACTIVE"
	RegionStateDeleted                RegionState   = "DELETED"
	RegionStateDeletionFailed         RegionState   = "DELETION_FAILED"
	RegionStateMaintenance            RegionState   = "MAINTENANCE"
	RegionStateInactive               RegionState   = "INACTIVE"
	RegionStateDeleting               RegionState   = "DELETING"
	RegionStateNew                    RegionState   = "NEW"
	KeystoneStateNew                  KeystoneState = "NEW"
	KeystoneStateInitializationFailed KeystoneState = "INITIALIZATION_FAILED"
	KeystoneStateInitialized          KeystoneState = "INITIALIZED"
	KeystoneStateDeleted              KeystoneState = "DELETED"
)

func (et EndpointType) IsValid() error {
	switch et {
	case EndpointTypeAdmin, EndpointTypeInternal, EndpointTypePublic:
		return nil
	}
	return fmt.Errorf("invalid EndpointType type: %v", et)
}

func (et EndpointType) ValidOrNil() (*EndpointType, error) {
	if et.String() == "" {
		return nil, nil
	}
	err := et.IsValid()
	if err != nil {
		return &et, err
	}
	return &et, nil
}

func (et EndpointType) String() string {
	return string(et)
}

func (et EndpointType) List() []EndpointType {
	return []EndpointType{EndpointTypeAdmin, EndpointTypeInternal, EndpointTypePublic}
}

func (et EndpointType) StringList() []string {
	var s []string
	for _, v := range et.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for EndpointType
func (et *EndpointType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := EndpointType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*et = v
	return nil
}

// MarshalJSON - implements Marshaler interface for EndpointType
func (et *EndpointType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (rs RegionState) IsValid() error {
	switch rs {
	case RegionStateActive,
		RegionStateDeleted,
		RegionStateDeleting,
		RegionStateDeletionFailed,
		RegionStateInactive,
		RegionStateMaintenance,
		RegionStateNew:
		return nil
	}
	return fmt.Errorf("invalid RegionState type: %v", rs)
}

func (rs RegionState) ValidOrNil() (*RegionState, error) {
	if rs.String() == "" {
		return nil, nil
	}
	err := rs.IsValid()
	if err != nil {
		return &rs, err
	}
	return &rs, nil
}

func (rs RegionState) String() string {
	return string(rs)
}

func (rs RegionState) List() []RegionState {
	return []RegionState{
		RegionStateActive,
		RegionStateDeleted,
		RegionStateDeleting,
		RegionStateDeletionFailed,
		RegionStateInactive,
		RegionStateMaintenance,
		RegionStateNew,
	}
}

func (rs RegionState) StringList() []string {
	var s []string
	for _, v := range rs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for RegionState
func (rs *RegionState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RegionState(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rs = v
	return nil
}

// MarshalJSON - implements Marshaler interface for RegionState
func (rs *RegionState) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

func (ks KeystoneState) IsValid() error {
	switch ks {
	case KeystoneStateInitializationFailed,
		KeystoneStateDeleted,
		KeystoneStateInitialized,
		KeystoneStateNew:
		return nil
	}
	return fmt.Errorf("invalid KeystoneState type: %v", ks)
}

func (ks KeystoneState) ValidOrNil() (*KeystoneState, error) {
	if ks.String() == "" {
		return nil, nil
	}
	err := ks.IsValid()
	if err != nil {
		return &ks, err
	}
	return &ks, nil
}

func (ks KeystoneState) String() string {
	return string(ks)
}

func (ks KeystoneState) List() []KeystoneState {
	return []KeystoneState{
		KeystoneStateInitializationFailed,
		KeystoneStateDeleted,
		KeystoneStateInitialized,
		KeystoneStateNew,
	}
}

func (ks KeystoneState) StringList() []string {
	var s []string
	for _, v := range ks.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for KeystoneState
func (ks *KeystoneState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := KeystoneState(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ks = v
	return nil
}

// MarshalJSON - implements Marshaler interface for KeystoneState
func (ks *KeystoneState) MarshalJSON() ([]byte, error) {
	return json.Marshal(ks.String())
}
