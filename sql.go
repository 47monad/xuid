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
// or MustWithPrefix after scanning.
//
// A nil UUID is encoded as SQL NULL.
func (x XUID) Value() (driver.Value, error) {
	if x.uuid == uuid.Nil() {
		return nil, nil
	}
	return x.uuid[:], nil
}

// NullXUID represents an XUID that may be NULL in a SQL database. It is
// the XUID counterpart to sql.NullString and implements driver.Valuer and
// sql.Scanner, so it can be used both as a query argument and as a scan
// destination for a nullable UUID column.
//
// Valid reports whether the value is not NULL, mirroring the other
// sql.Null* types. Unlike XUID.Value, which maps a nil UUID to NULL,
// NullXUID treats Valid as the single source of truth: with Valid set, its
// Value is non-NULL even when it holds the nil UUID, so a non-NULL all-zero
// UUID round-trips without collapsing to NULL.
type NullXUID struct {
	XUID  XUID
	Valid bool
}

// Value implements the driver.Valuer interface. It returns SQL NULL when
// Valid is false, and the 16 raw UUID bytes otherwise.
func (n NullXUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	id := n.XUID.GetUUID()
	return id[:], nil
}

// Scan implements the sql.Scanner interface. A NULL value sets Valid to
// false and leaves the XUID at its zero value. Any supported non-NULL value
// is scanned with XUID.Scan and sets Valid to true.
//
// As with XUID.Scan, the prefix is not stored in the database and is lost
// when scanning; restore it with WithPrefix or MustWithPrefix based on the
// table or column the value was read from.
func (n *NullXUID) Scan(value interface{}) error {
	if value == nil {
		n.XUID = XUID{}
		n.Valid = false
		return nil
	}

	var id XUID
	if err := id.Scan(value); err != nil {
		n.XUID = XUID{}
		n.Valid = false
		return err
	}
	n.XUID = id
	n.Valid = true
	return nil
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
// Restore prefixes with WithPrefix or MustWithPrefix after loading, based
// on the table or column the value was read from:
//
//	var loaded xuid.XUID
//	loaded.Scan(value)
//	restored := loaded.MustWithPrefix("user")
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
