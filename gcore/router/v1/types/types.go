package types

import (
	"encoding/json"
	"fmt"
)

type GatewayType string
type InterfaceType string

const (
	DefaultGateway      GatewayType   = "default"
	ManualGateway       GatewayType   = "manual"
	SubnetInterfaceType InterfaceType = "subnet"
)

func (gw GatewayType) String() string {
	return string(gw)
}

func (gw GatewayType) List() []GatewayType {
	return []GatewayType{DefaultGateway, ManualGateway}
}

func (gw GatewayType) StringList() []string {
	var s []string
	for _, v := range gw.List() {
		s = append(s, v.String())
	}
	return s
}

func (gw GatewayType) IsValid() error {
	switch gw {
	case DefaultGateway, ManualGateway:
		return nil
	}
	return fmt.Errorf("invalid GatewayType type: %v", gw)
}

func (gw GatewayType) ValidOrNil() (*GatewayType, error) {
	if gw.String() == "" {
		return nil, nil
	}
	err := gw.IsValid()
	if err != nil {
		return &gw, err
	}
	return &gw, nil
}

// UnmarshalJSON - implements Unmarshaler interface for GatewayType
func (gw *GatewayType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := GatewayType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*gw = v
	return nil
}

// MarshalJSON - implements Marshaler interface for GatewayType
func (gw *GatewayType) MarshalJSON() ([]byte, error) {
	return json.Marshal(gw.String())
}

func (it InterfaceType) String() string {
	return string(it)
}

func (it InterfaceType) List() []GatewayType {
	return []GatewayType{DefaultGateway, ManualGateway}
}

func (it InterfaceType) StringList() []string {
	var s []string
	for _, v := range it.List() {
		s = append(s, v.String())
	}
	return s
}

func (it InterfaceType) IsValid() error {
	switch it {
	case SubnetInterfaceType:
		return nil
	}
	return fmt.Errorf("invalid GatewayType type: %v", it)
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

// UnmarshalJSON - implements Unmarshaler interface for InterfaceType
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

// MarshalJSON - implements Marshaler interface for InterfaceType
func (it *InterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(it.String())
}
