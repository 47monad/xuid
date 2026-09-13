package xuid

import (
	"database/sql/driver"
	"fmt"
	"uuid"
)

// Value implements the driver.Valuer interface.
// This allows XUID to be stored in SQL databases as UUID.
func (x XUID) Value() (driver.Value, error) {
	if x.uuid == uuid.Nil() {
		return nil, nil
	}
	return x.uuid.String(), nil
}

// Scan implements the sql.Scanner interface.
// This allows XUID to be loaded from SQL databases.
// Note: The prefix information is lost when loading from database.
// Restore prefixes with WithPrefix after loading, based on the table
// or column the value was read from:
//
//	var loaded xuid.XUID
//	loaded.Scan(value)
//	restored := loaded.WithPrefix("user")
func (x *XUID) Scan(value interface{}) error {
	if value == nil {
		x.uuid = uuid.Nil()
		x.prefix = ""
		return nil
	}

	switch d := value.(type) {
	case string:
		id, err := uuid.Parse(d)
		if err != nil {
			return fmt.Errorf("%w: invalid UUID string %q", ErrScan, d)
		}
		x.uuid = id
		x.prefix = ""
		return nil
	case []byte:
		if len(d) != 16 {
			return fmt.Errorf("%w: invalid UUID byte length %d, want 16", ErrScan, len(d))
		}
		copy(x.uuid[:], d)
		x.prefix = "" // Prefix is lost when loading from database
		return nil
	}

	return fmt.Errorf("%w: unsupported type %T", ErrScan, value)
}
