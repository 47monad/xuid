package xuid

import (
	"encoding/json"
	"uuid"
)

// MarshalText implements the encoding.TextMarshaler interface.
//
// It returns the canonical XUID string, so XUID works as a JSON map
// key and with text-based encoders such as YAML, TOML and encoding/xml.
//
// A zero-value XUID (nil UUID) marshals to the empty string, mirroring
// the null value produced by MarshalJSON and the NULL value produced by
// Value for SQL storage.
func (x XUID) MarshalText() ([]byte, error) {
	if x.uuid == uuid.Nil() {
		return []byte{}, nil
	}
	return []byte(x.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
//
// An empty string maps to the zero-value XUID, mirroring the behavior
// of Scan and UnmarshalJSON. Any other value must be a valid XUID string.
func (x *XUID) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		x.uuid = uuid.Nil()
		x.prefix = ""
		return nil
	}
	xid, err := Parse(string(data))
	if err != nil {
		return err
	}
	x.uuid = xid.uuid
	x.prefix = xid.prefix
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
//
// It delegates to MarshalText and encodes the result as a JSON string,
// except for a zero-value XUID (empty text), which marshals to null.
func (x XUID) MarshalJSON() ([]byte, error) {
	text, err := x.MarshalText()
	if err != nil {
		return nil, err
	}
	if len(text) == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(string(text))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// null is accepted and maps to the zero-value XUID, mirroring the
// behavior of Scan for SQL storage. Any other value must be a string
// holding a valid XUID. It delegates to UnmarshalText so that
// null/empty handling lives in one place.
func (x *XUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return x.UnmarshalText(nil)
	}
	var res string
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	return x.UnmarshalText([]byte(res))
}
