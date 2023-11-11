package xtreamcodes

// Package base64 provides a byte slice type that marshals into json as
// a raw (no padding) base64url value.
// Originally from https://github.com/manifoldco/go-base64

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// Base64Value is a base64url encoded json object,
type Base64Value []byte

// New returns a pointer to a Base64Value, cast from the given byte slice.
// This is a convenience function, handling the address creation that a direct
// cast would not allow.
func New(b []byte) *Base64Value {
	v := Base64Value(b)
	return &v
}

// NewFromString returns a Base64Value containing the decoded data in encoded.
func NewFromString(encoded string) (*Base64Value, error) {
	out, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	return New(out), nil
}

// MarshalJSON returns the ba64url encoding of bv for JSON representation.
func (bv *Base64Value) MarshalJSON() ([]byte, error) {
	return []byte("\"" + base64.RawURLEncoding.EncodeToString(*bv) + "\""), nil
}

func (bv *Base64Value) String() string {
	return base64.RawURLEncoding.EncodeToString(*bv)
}

// UnmarshalJSON sets bv to the bytes represented in the base64url encoding b.
func (bv *Base64Value) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return errors.New("value is not a string")
	}

	// Unmarshal the JSON string
	var jsonString string
	if err := json.Unmarshal(b, &jsonString); err != nil {
		return err
	}

	// Determine the base64 encoding type and decode accordingly
	var decodedData []byte
	var err error
	if strings.ContainsAny(jsonString, "-_") {
		// Base64 URL encoded data
		decodedData, err = base64.RawURLEncoding.DecodeString(jsonString)
	} else {
		// Standard base64 encoded data
		decodedData, err = base64.StdEncoding.DecodeString(jsonString)
	}

	if err != nil {
		return err
	}

	v := reflect.ValueOf(bv).Elem()
	v.SetBytes(decodedData)
	return nil
}
