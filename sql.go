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
// As with XUID.Scan, a UUID-format or binary value carries no prefix and
// scanning it yields an empty prefix; an XUID-format value preserves its
// prefix. Restore a lost prefix with WithPrefix or MustWithPrefix based on
// the table or column the value was read from.
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
//   - string in UUID format or in the package's own XUID format
//   - []byte of exactly 16 raw bytes
//   - []byte containing a UUID or XUID string (e.g. lib/pq)
//   - [16]byte or uuid.UUID (e.g. pgx's native UUID type)
//
// A []byte of exactly 16 bytes is ambiguous: it may be the raw UUID bytes
// written by Value, or a 16-byte textual identifier such as the canonical
// form of the nil UUID, "1111111111111111". Scan prefers the text
// interpretation when the bytes are valid UUID/XUID text and otherwise
// treats them as raw UUID bytes, so text columns are read correctly.
// The trade-off is that the astronomically rare raw 16-byte UUID whose
// bytes also form valid XUID text is read as that text; pass such a value
// as a [16]byte or uuid.UUID to force the raw interpretation.
//
// A UUID-format value never carries a prefix, so scanning one yields an
// empty prefix. An XUID-format value preserves the prefix it encodes.
// Binary columns store no prefix, so restore one with WithPrefix or
// MustWithPrefix based on the table or column the value was read from:
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
		if x.scanText(d) {
			return nil
		}
		return fmt.Errorf("%w: invalid UUID string %q", ErrScan, d)
	case []byte:
		// A []byte of exactly 16 bytes is ambiguous between the raw UUID
		// bytes written by Value and 16-byte XUID text, so try the text
		// form first and fall back to raw bytes. Drivers such as lib/pq
		// deliver text columns as a []byte rather than a string.
		if x.scanText(string(d)) {
			return nil
		}
		if len(d) == 16 {
			copy(x.uuid[:], d)
			x.prefix = ""
			return nil
		}
		return fmt.Errorf("%w: invalid UUID bytes %q", ErrScan, d)
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

// scanText sets x from s if s is either a UUID-format string or one of the
// package's own XUID strings, and reports whether it succeeded. The UUID
// format is tried first: it is unambiguous and cannot collide with an XUID
// string, since XUID strings never contain hyphens and never exceed
// maxEncodedLen base58 digits (plus an optional prefix).
func (x *XUID) scanText(s string) bool {
	if id, err := uuid.Parse(s); err == nil {
		x.uuid = id
		x.prefix = ""
		return true
	}
	if id, err := Parse(s); err == nil {
		x.uuid = id.uuid
		x.prefix = id.prefix
		return true
	}
	return false
}
