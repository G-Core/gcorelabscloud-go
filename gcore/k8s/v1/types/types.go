package types

import (
	"encoding/json"
	"fmt"
)

type PoolRole string
type HealthStatus string
type IngressController string

const (
	NodegroupMasterRole      PoolRole          = "master"
	NodegroupWorkerRole      PoolRole          = "worker"
	HealthStatusUnknown      HealthStatus      = "UNKNOWN"
	HealthStatusHealthy      HealthStatus      = "HEALTHY"
	HealthStatusUnHealthy    HealthStatus      = "UNHEALTHY"
	IngressControllerOctavia IngressController = "octavia"
	IngressControllerNginx   IngressController = "nginx"
	IngressControllerTraefik IngressController = "traefik"
)

func (ng PoolRole) IsValid() error {
	switch ng {
	case NodegroupMasterRole,
		NodegroupWorkerRole:
		return nil
	}
	return fmt.Errorf("invalid PoolRole type: %v", ng)
}

func (ng PoolRole) ValidOrNil() (*PoolRole, error) {
	if ng.String() == "" {
		return nil, nil
	}
	err := ng.IsValid()
	if err != nil {
		return &ng, err
	}
	return &ng, nil
}

func (ng PoolRole) String() string {
	return string(ng)
}

func (ng PoolRole) List() []PoolRole {
	return []PoolRole{
		NodegroupMasterRole,
		NodegroupWorkerRole,
	}
}

func (ng PoolRole) StringList() []string {
	var s []string
	for _, v := range ng.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ng *PoolRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := PoolRole(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ng = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ng *PoolRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ng.String())
}

func (hs HealthStatus) IsValid() error {
	switch hs {
	case HealthStatusHealthy,
		HealthStatusUnHealthy,
		HealthStatusUnknown:
		return nil
	}
	return fmt.Errorf("invalid HealthStatus type: %v", hs)
}

func (hs HealthStatus) ValidOrNil() (*HealthStatus, error) {
	if hs.String() == "" {
		return nil, nil
	}
	err := hs.IsValid()
	if err != nil {
		return &hs, err
	}
	return &hs, nil
}

func (hs HealthStatus) String() string {
	return string(hs)
}

func (hs HealthStatus) List() []HealthStatus {
	return []HealthStatus{
		HealthStatusUnHealthy,
		HealthStatusUnknown,
		HealthStatusUnHealthy,
	}
}

func (hs HealthStatus) StringList() []string {
	var s []string
	for _, v := range hs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (hs *HealthStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := HealthStatus(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*hs = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (hs *HealthStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(hs.String())
}

func (ic IngressController) IsValid() error {
	switch ic {
	case IngressControllerOctavia,
		IngressControllerNginx,
		IngressControllerTraefik:
		return nil
	}
	return fmt.Errorf("invalid IngressController type: %v", ic)
}

func (ic IngressController) ValidOrNil() (*IngressController, error) {
	if ic.String() == "" {
		return nil, nil
	}
	err := ic.IsValid()
	if err != nil {
		return &ic, err
	}
	return &ic, nil
}

func (ic IngressController) String() string {
	return string(ic)
}

func (ic IngressController) List() []IngressController {
	return []IngressController{
		IngressControllerOctavia,
		IngressControllerNginx,
		IngressControllerTraefik,
	}
}

func (ic IngressController) StringList() []string {
	var s []string
	for _, v := range ic.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ic *IngressController) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := IngressController(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ic = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ic *IngressController) MarshalJSON() ([]byte, error) {
	return json.Marshal(ic.String())
}
