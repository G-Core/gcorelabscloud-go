package types

import (
	"encoding/json"
	"fmt"
)

type InstanceActionType string

const (
	InstanceActionTypeStart      InstanceActionType = "start"
	InstanceActionTypeStop       InstanceActionType = "stop"
	InstanceActionTypeReboot     InstanceActionType = "reboot"
	InstanceActionTypeRebootHard InstanceActionType = "reboot_hard"
	InstanceActionTypeSuspend    InstanceActionType = "suspend"
	InstanceActionTypeResume     InstanceActionType = "resume"
)

func (iat InstanceActionType) IsValid() error {
	switch iat {
	case InstanceActionTypeStart,
		InstanceActionTypeStop,
		InstanceActionTypeReboot,
		InstanceActionTypeRebootHard,
		InstanceActionTypeSuspend,
		InstanceActionTypeResume:
		return nil
	}
	return fmt.Errorf("invalid InstanceActionType type: %v", iat)
}

func (iat InstanceActionType) ValidOrNil() (*InstanceActionType, error) {
	if iat.String() == "" {
		return nil, nil
	}
	err := iat.IsValid()
	if err != nil {
		return &iat, err
	}
	return &iat, nil
}

func (iat InstanceActionType) String() string {
	return string(iat)
}

func (iat InstanceActionType) List() []InstanceActionType {
	return []InstanceActionType{
		InstanceActionTypeStart,
		InstanceActionTypeStop,
		InstanceActionTypeReboot,
		InstanceActionTypeRebootHard,
		InstanceActionTypeSuspend,
		InstanceActionTypeResume,
	}
}

func (iat InstanceActionType) StringList() []string {
	var s []string
	for _, v := range iat.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (iat *InstanceActionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := InstanceActionType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*iat = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (iat *InstanceActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(iat.String())
}
