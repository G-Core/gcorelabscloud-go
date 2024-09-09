package clusters

import (
	"encoding/json"
	"fmt"
)

type CNIProvider string
type TunnelType string
type LBModeType string
type RoutingModeType string

const (
	CalicoProvider    CNIProvider     = "calico"
	CiliumProvider    CNIProvider     = "cilium"
	VXLANTunnel       TunnelType      = "vxlan"
	GeneveTunnel      TunnelType      = "geneve"
	EmptyTunnel       TunnelType      = ""
	SNATLBMode        LBModeType      = "snat"
	DSRLBMode         LBModeType      = "dsr"
	HybridLBMode      LBModeType      = "hybrid"
	NativeRoutingMode RoutingModeType = "native"
	TunnelRoutingMode RoutingModeType = "tunnel"
)

func (cn CNIProvider) IsValid() error {
	switch cn {
	case CalicoProvider, CiliumProvider:
		return nil
	}
	return fmt.Errorf("invalid CNIProvider type: %v", cn)
}

func (cn CNIProvider) ValidOrNil() (*CNIProvider, error) {
	if cn.String() == "" {
		return nil, nil
	}
	err := cn.IsValid()
	if err != nil {
		return &cn, err
	}
	return &cn, nil
}

func (cn CNIProvider) String() string {
	return string(cn)
}

func (cn CNIProvider) List() []CNIProvider {
	return []CNIProvider{
		CalicoProvider,
		CiliumProvider,
	}
}

func (cn CNIProvider) StringList() []string {
	var s []string
	for _, v := range cn.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for CNIProvider
func (cn *CNIProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := CNIProvider(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*cn = v
	return nil
}

// MarshalJSON - implements Marshaler interface for CNIProvider
func (cn *CNIProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(cn.String())
}

func (t TunnelType) IsValid() error {
	switch t {
	case VXLANTunnel, GeneveTunnel, EmptyTunnel:
		return nil
	}
	return fmt.Errorf("invalid TunnelType type: %v", t)
}

func (t TunnelType) ValidOrNil() (*TunnelType, error) {
	if t.String() == "" {
		return nil, nil
	}
	err := t.IsValid()
	if err != nil {
		return &t, err
	}
	return &t, nil
}

func (t TunnelType) String() string {
	return string(t)
}

func (t TunnelType) List() []TunnelType {
	return []TunnelType{
		VXLANTunnel,
		GeneveTunnel,
		EmptyTunnel,
	}
}

func (t TunnelType) StringList() []string {
	var s []string
	for _, v := range t.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for TunnelType
func (t *TunnelType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := TunnelType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*t = v
	return nil
}

// MarshalJSON - implements Marshaler interface for TunnelType
func (t *TunnelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t LBModeType) IsValid() error {
	switch t {
	case DSRLBMode, SNATLBMode, HybridLBMode:
		return nil
	}
	return fmt.Errorf("invalid LBModeType type: %v", t)
}

func (t LBModeType) ValidOrNil() (*LBModeType, error) {
	if t.String() == "" {
		return nil, nil
	}
	err := t.IsValid()
	if err != nil {
		return &t, err
	}
	return &t, nil
}

func (t LBModeType) String() string {
	return string(t)
}

func (t LBModeType) List() []LBModeType {
	return []LBModeType{
		DSRLBMode,
		SNATLBMode,
		HybridLBMode,
	}
}

func (t LBModeType) StringList() []string {
	var s []string
	for _, v := range t.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for LBModeType
func (t *LBModeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := LBModeType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*t = v
	return nil
}

// MarshalJSON - implements Marshaler interface for LBModeType
func (t *LBModeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t RoutingModeType) IsValid() error {
	switch t {
	case NativeRoutingMode, TunnelRoutingMode:
		return nil
	}
	return fmt.Errorf("invalid RoutingModeType type: %v", t)
}

func (t RoutingModeType) ValidOrNil() (*RoutingModeType, error) {
	if t.String() == "" {
		return nil, nil
	}
	err := t.IsValid()
	if err != nil {
		return &t, err
	}
	return &t, nil
}

func (t RoutingModeType) String() string {
	return string(t)
}

func (t RoutingModeType) List() []RoutingModeType {
	return []RoutingModeType{
		NativeRoutingMode,
		TunnelRoutingMode,
	}
}

func (t RoutingModeType) StringList() []string {
	var s []string
	for _, v := range t.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for RoutingModeType
func (t *RoutingModeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := RoutingModeType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*t = v
	return nil
}

// MarshalJSON - implements Marshaler interface for RoutingModeType
func (t *RoutingModeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
