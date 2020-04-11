package types

import (
	"encoding/json"
	"fmt"
)

type KeystoneState string

const (
	KeystoneStateNew                  KeystoneState = "NEW"
	KeystoneStateInitializationFailed KeystoneState = "INITIALIZATION_FAILED"
	KeystoneStateInitialized          KeystoneState = "INITIALIZED"
	KeystoneStateDeleted              KeystoneState = "DELETED"
)

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
