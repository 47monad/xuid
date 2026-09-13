package xuid

import (
	"database/sql/driver"
	"fmt"
	"uuid"
)

// Value implements the driver.Valuer interface.
//
// It returns the 16 raw UUID bytes as a []byte so drivers store XUIDs in
// binary UUID columns (BYTEA in PostgreSQL, BINARY(16) in MySQL) rather
// than a text column. The prefix is not stored; restore it with WithPrefix
// after scanning.
//
// A nil UUID is encoded as SQL NULL.
func (x XUID) Value() (driver.Value, error) {
	if x.uuid == uuid.Nil() {
		return nil, nil
	}
	return x.uuid[:], nil
}

// Scan implements the sql.Scanner interface.
// This allows XUID to be loaded from SQL databases.
//
// It accepts the shapes drivers actually deliver UUID columns in:
//
//   - string in UUID format
//   - []byte of exactly 16 raw bytes
//   - []byte containing a UUID string (e.g. lib/pq)
//   - [16]byte or uuid.UUID (e.g. pgx's native UUID type)
//
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
		var id uuid.UUID
		if len(d) == 16 {
			copy(id[:], d)
		} else {
			// Drivers such as lib/pq deliver UUID columns as a []byte
			// holding the textual form rather than the 16 raw bytes.
			var err error
			id, err = uuid.Parse(string(d))
			if err != nil {
				return fmt.Errorf("%w: invalid UUID bytes %q", ErrScan, d)
			}
		}
		x.uuid = id
		x.prefix = "" // Prefix is lost when loading from database
		return nil
	case [16]byte:
		x.uuid = uuid.UUID(d)
		x.prefix = ""
		return nil
	case uuid.UUID:
		x.uuid = d
		x.prefix = ""
		return nil
	}

	return fmt.Errorf("%w: unsupported type %T", ErrScan, value)
}
