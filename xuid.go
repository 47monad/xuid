// Package xuid provides a compact, type-safe identifier system built on UUIDs.
//
// XUID combines the robustness of UUIDs with practical enhancements:
// - Sortable identifiers using UUIDv7 for chronological ordering
// - Optional string prefixes for human-readable context (e.g., "user_", "order_")
// - Base58 encoding for shorter, URL-safe representations
// - Built-in JSON marshaling and unmarshaling support
//
// Example usage:
//
//	// Create a sortable identifier with prefix
//	userID := xuid.MustNewSortable("user")
//	fmt.Println(userID.String()) // user_8M7Qq2vR3kGbF9wN5pL2xA
//
//	// Parse from string
//	parsed, err := xuid.Parse("user_8M7Qq2vR3kGbF9wN5pL2xA")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(parsed.GetPrefix()) // user
package xuid

import (
	"fmt"
	"strings"
	"uuid"
)

type XUID struct {
	uuid   uuid.UUID
	prefix string
}

func NewWith(id uuid.UUID, prefix string) (XUID, error) {
	return XUID{
		uuid:   id,
		prefix: prefix,
	}, nil
}

func NewSortable(prefix string) (XUID, error) {
	return XUID{
		uuid:   uuid.NewV7(),
		prefix: prefix,
	}, nil
}

func MustNewSortable(prefix string) XUID {
	return Must(NewSortable(prefix))
}

func NewRandom(prefix string) (XUID, error) {
	return XUID{
		uuid:   uuid.NewV4(),
		prefix: prefix,
	}, nil
}

func MustNewRandom(prefix string) XUID {
	return Must(NewRandom(prefix))
}

func NilUUID() (XUID, error) {
	return NewWith(uuid.Nil(), "")
}

func (x XUID) GetUUID() uuid.UUID {
	return x.uuid
}

func (x XUID) IsSortable() bool {
	return x.uuid[6]>>4 == 7
}

func (x XUID) IsRandom() bool {
	return x.uuid[6]>>4 == 4
}

func (x XUID) GetPrefix() string {
	return x.prefix
}

// WithPrefix returns a copy of x with the prefix set to prefix.
// It is the supported way to restore a prefix after loading the
// underlying UUID from a database column, since Scan discards prefixes:
//
//	restored := xuid.MustParse(s).WithPrefix("user")
//
// Being immutable, it chains off any value, including non-addressable
// ones such as function results.
func (x XUID) WithPrefix(prefix string) XUID {
	x.prefix = prefix
	return x
}

// SetPrefix sets the prefix field to the specified prefix.
// This is useful when loading XUIDs from database and need to restore the prefix.
//
// Deprecated: use WithPrefix instead. SetPrefix mixes mutation with
// chaining semantics and cannot be chained off non-addressable values,
// such as function results.
func (x *XUID) SetPrefix(prefix string) *XUID {
	x.prefix = prefix
	return x
}

func (x XUID) String() string {
	if x.prefix == "" {
		return encodeBase58(x.uuid[:])
	}
	return x.prefix + "_" + encodeBase58(x.uuid[:])
}

// Equal reports whether x and y identify the same XUID.
//
// Two XUIDs are equal only if both their UUID and prefix match.
// A different prefix means the identifiers are not equal, even when
// they carry the same UUID, because the prefix is part of the
// identifier's identity. To compare only the underlying UUIDs,
// use EqualUUID.
func (x XUID) Equal(y XUID) bool {
	return x.uuid == y.uuid && x.prefix == y.prefix
}

// EqualUUID reports whether x and y carry the same UUID, ignoring
// their prefixes. It is useful when the same identifier is stored
// under different prefixes (e.g. after restoring a prefix from a
// database column).
func (x XUID) EqualUUID(y XUID) bool {
	return x.uuid == y.uuid
}

// Compare returns an integer comparing two XUIDs.
//
// The result is -1 if x sorts before y, 0 if x and y are ordered
// identically, and +1 if x sorts after y.
//
// XUIDs are ordered by prefix first (an empty prefix sorts before
// any non-empty prefix), then by UUID bytes. For UUIDv7 identifiers
// sharing a prefix, the UUID bytes preserve chronological order, so
// Compare yields time-ordered results within each prefix.
//
// Note that Compare treats XUIDs as ordered values, not identical
// ones: Compare returns 0 only when both prefix and UUID match.
func Compare(x, y XUID) int {
	if c := strings.Compare(x.prefix, y.prefix); c != 0 {
		return c
	}
	return x.uuid.Compare(y.uuid)
}

// Less reports whether x sorts before y. It is equivalent to
// Compare(x, y) < 0.
func Less(x, y XUID) bool {
	return Compare(x, y) < 0
}

// Parse decodes an XUID string. The final underscore separates the
// optional prefix from the base58-encoded UUID, so any underscores
// inside the prefix are preserved:
//
//	user_admin_8M7Qq2vR3kGbF9wN5pL2xA -> prefix "user_admin"
//
// Every failure wraps ErrParse with the underlying cause, so callers
// can detect them with errors.Is(err, ErrParse) while still logging a
// message that explains why parsing failed.
func Parse(idstr string) (XUID, error) {
	prefix := ""
	uuidstr := idstr
	if i := strings.LastIndex(idstr, "_"); i >= 0 {
		prefix = idstr[:i]
		uuidstr = idstr[i+1:]
	}

	// A 16-byte UUID never encodes to more than maxEncodedLen base58
	// digits, so anything longer (or empty) is invalid before decoding.
	// This also bounds decode work on adversarially long input.
	if len(uuidstr) == 0 {
		return XUID{}, fmt.Errorf("%w: empty identifier", ErrParse)
	}
	if len(uuidstr) > maxEncodedLen {
		return XUID{}, fmt.Errorf("%w: identifier length %d exceeds %d", ErrParse, len(uuidstr), maxEncodedLen)
	}

	decoded, err := decodeBase58(uuidstr)
	if err != nil {
		return XUID{}, fmt.Errorf("%w: %v", ErrParse, err)
	}
	var id uuid.UUID
	if len(decoded) != len(id) {
		return XUID{}, fmt.Errorf("%w: decoded %d bytes, want %d", ErrParse, len(decoded), len(id))
	}
	copy(id[:], decoded)
	return NewWith(id, prefix)
}

func MustParse(idstr string) XUID {
	xid, err := Parse(idstr)
	if err != nil {
		panic(err)
	}
	return xid
}

func IsValid(idstr string) bool {
	_, err := Parse(idstr)
	return err == nil
}

func Must(xid XUID, err error) XUID {
	if err != nil {
		panic(err)
	}
	return xid
}

func IsEmpty(xid XUID) bool {
	return xid.uuid == uuid.Nil()
}
