package xuid

import (
	"encoding/json"
	"uuid"
)

// MarshalJSON implements the json.Marshaler interface.
// A zero-value XUID (nil UUID, empty prefix) marshals to null,
// mirroring the behavior of Value for SQL storage.
func (x XUID) MarshalJSON() ([]byte, error) {
	if x.uuid == uuid.Nil() {
		return []byte("null"), nil
	}
	return json.Marshal(x.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// null is accepted and maps to the zero-value XUID, mirroring
// the behavior of Scan for SQL storage. Any other value must
// be a valid XUID string.
func (x *XUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		x.uuid = uuid.Nil()
		x.prefix = ""
		return nil
	}
	var res string
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	xid, err := Parse(res)
	if err != nil {
		return err
	}
	x.uuid = xid.uuid
	x.prefix = xid.prefix
	return nil
}
