package gcorecloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/ladydascalie/currency"

	log "github.com/sirupsen/logrus"
)

/*
Result is an internal type to be used by individual resource packages, but its
methods will be available on a wide variety of user-facing embedding types.

It acts as a base struct that other Result types, returned from request
functions, can embed for convenience. All Results capture basic information
from the HTTP transaction that was performed, including the response body,
HTTP headers, and any errors that happened.

Generally, each Result type will have an Extract method that can be used to
further interpret the result's payload in a specific context. Extensions or
providers can then provide additional extraction functions to pull out
provider- or extension-specific information as well.
*/
type Result struct {
	// Body is the payload of the HTTP response from the server. In most cases,
	// this will be the deserialized JSON structure.
	Body interface{}

	// Header contains the HTTP header structure from the original response.
	Header http.Header

	// Err is an error that occurred during the operation. It's deferred until
	// extraction to make it easier to chain the Extract call.
	Err error
}

// ExtractInto allows users to provide an object into which `Extract` will extract
// the `Result.Body`. This would be useful for OpenStack providers that have
// different fields in the response object than OpenStack proper.
func (r Result) ExtractInto(to interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	if reader, ok := r.Body.(io.Reader); ok {
		if readCloser, ok := reader.(io.Closer); ok {
			defer func() {
				err := readCloser.Close()
				if err != nil {
					log.Error(err)
				}
			}()
		}
		return json.NewDecoder(reader).Decode(to)
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)

	return err
}

func (r Result) extractIntoPtr(to interface{}, label string) error {
	if label == "" {
		return r.ExtractInto(&to)
	}

	var m map[string]interface{}
	err := r.ExtractInto(&m)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m[label])
	if err != nil {
		return err
	}
	toValue := reflect.ValueOf(to)
	if toValue.Kind() == reflect.Ptr {
		toValue = toValue.Elem()
	}

	switch toValue.Kind() {
	case reflect.Slice:
		typeOfV := toValue.Type().Elem()
		if typeOfV.Kind() == reflect.Struct {
			if typeOfV.NumField() > 0 && typeOfV.Field(0).Anonymous {
				newSlice := reflect.MakeSlice(reflect.SliceOf(typeOfV), 0, 0)

				if mSlice, ok := m[label].([]interface{}); ok {
					for _, v := range mSlice {
						// For each iteration of the slice, we create a new struct.
						// This is to work around a bug where elements of a slice
						// are reused and not overwritten when the same copy of the
						// struct is used:
						//
						// https://github.com/golang/go/issues/21092
						// https://github.com/golang/go/issues/24155
						// https://play.golang.org/p/NHo3ywlPZli
						newType := reflect.New(typeOfV).Elem()

						b, err := json.Marshal(v)
						if err != nil {
							return err
						}

						// This is needed for structs with an UnmarshalJSON method.
						// Technically this is just unmarshalling the response into
						// a struct that is never used, but it's good enough to
						// trigger the UnmarshalJSON method.
						for i := 0; i < newType.NumField(); i++ {
							s := newType.Field(i).Addr().Interface()

							// Unmarshal is used rather than NewDecoder to also work
							// around the above-mentioned bug.
							err = json.Unmarshal(b, s)
							if err != nil {
								return err
							}
						}

						newSlice = reflect.Append(newSlice, newType)
					}
				}

				// "to" should now be properly modeled to receive the
				// JSON response body and unmarshal into all the correct
				// fields of the struct or composed extension struct
				// at the end of this method.
				toValue.Set(newSlice)
			}
		}
	case reflect.Struct:
		typeOfV := toValue.Type()
		if typeOfV.NumField() > 0 && typeOfV.Field(0).Anonymous {
			for i := 0; i < toValue.NumField(); i++ {
				toField := toValue.Field(i)
				if toField.Kind() == reflect.Struct {
					s := toField.Addr().Interface()
					err = json.NewDecoder(bytes.NewReader(b)).Decode(s)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	err = json.Unmarshal(b, &to)
	return err
}

// ExtractIntoStructPtr will unmarshal the Result (r) into the provided
// interface{} (to).
//
// NOTE: For internal use only
//
// `to` must be a pointer to an underlying struct type
//
// If provided, `label` will be filtered out of the response
// body prior to `r` being unmarshalled into `to`.
func (r Result) ExtractIntoStructPtr(to interface{}, label string) error {
	if r.Err != nil {
		return r.Err
	}

	t := reflect.TypeOf(to)
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Struct:
		return r.extractIntoPtr(to, label)
	default:
		return fmt.Errorf("expected pointer to struct, got: %v", t)
	}
}

// ExtractIntoMapPtr will unmarshal the Result (r) into the provided
// interface{} (to).
//
// NOTE: For internal use only
//
// `to` must be a pointer to an underlying map type
//
// If provided, `label` will be filtered out of the response
// body prior to `r` being unmarshalled into `to`.
func (r Result) ExtractIntoMapPtr(to interface{}, label string) error {
	if r.Err != nil {
		return r.Err
	}

	t := reflect.TypeOf(to)
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Map:
		return r.extractIntoPtr(to, label)
	default:
		return fmt.Errorf("expected pointer to map, got: %v", t)
	}
}

// ExtractIntoSlicePtr will unmarshal the Result (r) into the provided
// interface{} (to).
//
// NOTE: For internal use only
//
// `to` must be a pointer to an underlying slice type
//
// If provided, `label` will be filtered out of the response
// body prior to `r` being unmarshalled into `to`.
func (r Result) ExtractIntoSlicePtr(to interface{}, label string) error {
	if r.Err != nil {
		return r.Err
	}

	t := reflect.TypeOf(to)
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Slice:
		return r.extractIntoPtr(to, label)
	default:
		return fmt.Errorf("expected pointer to slice, got: %v", t)
	}
}

// PrettyPrintJSON creates a string containing the full response body as
// pretty-printed JSON. It's useful for capturing test fixtures and for
// debugging extraction bugs. If you include its output in an issue related to
// a buggy extraction function, we will all love you forever.
func (r Result) PrettyPrintJSON() string {
	pretty, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return string(pretty)
}

// ErrResult is an internal type to be used by individual resource packages, but
// its methods will be available on a wide variety of user-facing embedding
// types.
//
// It represents results that only contain a potential error and
// nothing else. Usually, if the operation executed successfully, the Err field
// will be nil; otherwise it will be stocked with a relevant error. Use the
// ExtractErr method
// to cleanly pull it out.
type ErrResult struct {
	Result
}

// ExtractErr is a function that extracts error information, or nil, from a result.
func (r ErrResult) ExtractErr() error {
	return r.Err
}

/*
HeaderResult is an internal type to be used by individual resource packages, but
its methods will be available on a wide variety of user-facing embedding types.

It represents a result that only contains an error (possibly nil) and an
http.Header. This is used, for example, by the objectstorage packages in
openstack, because most of the operations don't return response bodies, but do
have relevant information in headers.
*/
type HeaderResult struct {
	Result
}

// ExtractInto allows users to provide an object into which `Extract` will
// extract the http.Header headers of the result.
func (r HeaderResult) ExtractInto(to interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	tmpHeaderMap := map[string]string{}
	for k, v := range r.Header {
		if len(v) > 0 {
			tmpHeaderMap[k] = v[0]
		}
	}

	b, err := json.Marshal(tmpHeaderMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)

	return err
}

// RFC3339Milli describes a common time format used by some API responses.
const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

// JSONRFC3339Milli describes time.Time in RFC3339Milli format.
type JSONRFC3339Milli time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339Milli
func (jt *JSONRFC3339Milli) UnmarshalJSON(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(RFC3339Milli, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

// RFC3339MilliNoZ is the time format with millis.
const RFC3339MilliNoZ = "2006-01-02T15:04:05.999999"

// JSONRFC3339MilliNoZ describes time.Time in RFC3339MilliNoZ format.
type JSONRFC3339MilliNoZ time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339MilliNoZ
func (jt *JSONRFC3339MilliNoZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339MilliNoZ, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339MilliNoZ(t)
	return nil
}

// JSONRFC1123 describes time.Time in time.RFC1123 format.
type JSONRFC1123 time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC1123
func (jt *JSONRFC1123) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC1123, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC1123(t)
	return nil
}

// JSONUnix describes time.Time in unix format.
type JSONUnix time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONUnix
func (jt *JSONUnix) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	unix, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	t = time.Unix(unix, 0)
	*jt = JSONUnix(t)
	return nil
}

// RFC3339NoZ is the time format used in Heat (Orchestration).
const RFC3339NoZ = "2006-01-02T15:04:05"

// JSONRFC3339NoZ describes time.Time in RFC3339NoZ format.
type JSONRFC3339NoZ struct {
	time.Time
}

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339NoZ
func (jt *JSONRFC3339NoZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339NoZ, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

func (jt *JSONRFC3339NoZ) String() string {
	return jt.Format(RFC3339NoZ)
}

// MarshalJSON - implements Marshaler interface for JSONRFC3339NoZ
func (jt *JSONRFC3339NoZ) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Format(RFC3339NoZ))
}

// RFC3339Z is the time format used in Heat (Orchestration).
const RFC3339Z = "2006-01-02T15:04:05-0700"

// JSONRFC3339Z describes time.Time in RFC3339Z format.
type JSONRFC3339Z struct {
	time.Time
}

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339Z
func (jt *JSONRFC3339Z) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339Z, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

// MarshalJSON - implements Marshaler interface for JSONRFC3339Z
func (jt *JSONRFC3339Z) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Format(RFC3339Z))
}

// RFC3339ZColon is the time format used in secrets.
const RFC3339ZColon = "2006-01-02T15:04:05-07:00"

// JSONRFC3339ZColon describes time.Time in RFC3339ZColon format.
type JSONRFC3339ZColon struct {
	time.Time
}

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339ZColon
func (jt *JSONRFC3339ZColon) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339ZColon, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

// MarshalJSON - implements Marshaler interface for JSONRFC3339ZColon
func (jt *JSONRFC3339ZColon) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Format(RFC3339ZColon))
}

// RFC3339ZZ describes a common time format used by some API responses.
const RFC3339ZZ = "2006-01-02T15:04:05Z"

// JSONRFC3339ZZ describes time.Time in RFC3339ZZ format
type JSONRFC3339ZZ struct {
	time.Time
}

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339ZZ
func (jt *JSONRFC3339ZZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339ZZ, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

// MarshalJSON - implements Marshaler interface for JSONRFC3339ZZ
func (jt *JSONRFC3339ZZ) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Format(RFC3339ZZ))
}

// RFC3339ZNoT is the time format used in Zun (Containers Service).
const RFC3339ZNoT = "2006-01-02 15:04:05-07:00"

// JSONRFC3339ZNoT describes time.Time in RFC3339ZNoT format.
type JSONRFC3339ZNoT time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339ZNoT
func (jt *JSONRFC3339ZNoT) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339ZNoT, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339ZNoT(t)
	return nil
}

// RFC3339ZNoTNoZ is another time format used in Zun (Containers Service).
const RFC3339ZNoTNoZ = "2006-01-02 15:04:05"

// JSONRFC3339ZNoTNoZ describes time.Time in RFC3339ZNoTNoZ format.
type JSONRFC3339ZNoTNoZ time.Time

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339ZNoTNoZ
func (jt *JSONRFC3339ZNoTNoZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339ZNoTNoZ, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339ZNoTNoZ(t)
	return nil
}

// RFC3339Date describes a common time format used by some API responses.
const RFC3339Date = "2006-01-02"

// JSONRFC3339Date describes time.Time in RFC3339Date format
type JSONRFC3339Date struct {
	time.Time
}

// UnmarshalJSON - implements Unmarshaler interface for JSONRFC3339Date
func (jt *JSONRFC3339Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339Date, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

// MarshalJSON - implements Marshaler interface for JSONRFC3339Date
func (jt *JSONRFC3339Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(jt.Format(RFC3339Date))
}

/*
Link is an internal type to be used in packages of collection resources that are
paginated in a certain way.

It's a response substructure common to many paginated collection results that is
used to point to related pages. Usually, the one we care about is the one with
Rel field set to "next".
*/
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

/*
ExtractNextURL is an internal function useful for packages of collection
resources that are paginated in a certain way.

It attempts to extract the "next" URL from slice of Link structs, or
"" if no such URL is present.
*/
func ExtractNextURL(links []Link) (string, error) {
	var nextURL string

	for _, l := range links {
		if l.Rel == "next" {
			nextURL = l.Href
		}
	}

	if nextURL == "" {
		return "", nil
	}

	return nextURL, nil
}

type CIDR struct {
	net.IPNet
}

func ParseCIDRString(s string) (*CIDR, error) {
	_, nt, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &CIDR{IPNet: *nt}, nil
}

func ParseCIDRStringOrNil(s string) (*CIDR, error) {
	if s == "" {
		return nil, nil
	}
	return ParseCIDRString(s)
}

// UnmarshalJSON - implements Unmarshaler interface for CIDR
func (c *CIDR) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cd, err := ParseCIDRString(s)
	if err != nil {
		return err
	}
	*c = *cd
	return nil
}

// MarshalJSON - implements Marshaler interface for CIDR
func (c CIDR) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// String - implements Stringer
func (c CIDR) String() string {
	return c.IPNet.String()
}

type Currency struct {
	*currency.Currency
}

func ParseCurrency(s string) (*Currency, error) {
	c, err := currency.Get(s)
	if err != nil {
		return nil, err
	}
	return &Currency{Currency: c}, nil
}

// String - implements Stringer
func (c Currency) String() string {
	return c.Currency.Code()
}

// UnmarshalJSON - implements Unmarshaler interface for Currency
func (c *Currency) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := ParseCurrency(s)
	if err != nil {
		return err
	}
	*c = *v
	return nil
}

// MarshalJSON - implements Marshaler interface for Currency
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type GcoreErrorType struct {
	ExceptionClass string `json:"exception_class"`
	Message        string `json:"message"`
	Traceback      string `json:"traceback"`
}

type MAC struct {
	net.HardwareAddr
}

func ParseMacString(s string) (*MAC, error) {
	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	return &MAC{HardwareAddr: mac}, nil
}

// UnmarshalJSON - implements Unmarshaler interface for MAC
func (m *MAC) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	mc, err := ParseMacString(s)
	if err != nil {
		return err
	}
	*m = *mc
	return nil
}

// MarshalJSON - implements Marshaler interface for MAC
func (m MAC) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// String - implements Stringer
func (m MAC) String() string {
	return m.HardwareAddr.String()
}

type URL struct {
	*url.URL
}

func ParseURL(s string) (*URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return &URL{URL: u}, nil
}

func MustParseURL(s string) *URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return &URL{URL: u}
}

func ParseURLNonMandatory(s string) (*URL, error) {
	if s == "" {
		return nil, nil
	}
	return ParseURL(s)
}

// UnmarshalJSON - implements Unmarshaler interface
func (u *URL) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ur, err := ParseURL(s)
	if err != nil {
		return err
	}
	*u = *ur
	return nil
}

// MarshalJSON - implements Marshaler interface
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// String - implements Stringer
func (u URL) String() string {
	return u.URL.String()
}

type ItemName struct {
	Name string `json:"name"`
}

type ItemIDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemID struct {
	ID string `json:"id"`
}
