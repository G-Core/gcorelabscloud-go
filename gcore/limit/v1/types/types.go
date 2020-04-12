package types

import (
	"encoding/json"
	"fmt"
)

type LimitRequestStatus string

const (
	LimitRequestInProgress LimitRequestStatus = "in progress"
	LimitRequestRejected   LimitRequestStatus = "rejected"
	LimitRequestDone       LimitRequestStatus = "done"
)

func (ks LimitRequestStatus) IsValid() error {
	switch ks {
	case LimitRequestRejected,
		LimitRequestDone,
		LimitRequestInProgress:
		return nil
	}
	return fmt.Errorf("invalid LimitRequestStatus type: %v", ks)
}

func (ks LimitRequestStatus) ValidOrNil() (*LimitRequestStatus, error) {
	if ks.String() == "" {
		return nil, nil
	}
	err := ks.IsValid()
	if err != nil {
		return &ks, err
	}
	return &ks, nil
}

func (ks LimitRequestStatus) String() string {
	return string(ks)
}

func (ks LimitRequestStatus) List() []LimitRequestStatus {
	return []LimitRequestStatus{
		LimitRequestRejected,
		LimitRequestDone,
		LimitRequestInProgress,
	}
}

func (ks LimitRequestStatus) StringList() []string {
	var s []string
	for _, v := range ks.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for LimitRequestStatus
func (ks *LimitRequestStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := LimitRequestStatus(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ks = v
	return nil
}

// MarshalJSON - implements Marshaler interface for LimitRequestStatus
func (ks *LimitRequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ks.String())
}
