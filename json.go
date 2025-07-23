package xuid

import (
	"encoding/json"
)

func (x XUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}

func (x *XUID) UnmarshalJSON(data []byte) error {
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
