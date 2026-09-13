package xuid

import "errors"

var (
	// ErrParse is returned by Parse when a string is not a valid XUID.
	// Parse wraps it with the underlying cause, so failures can be
	// detected with errors.Is(err, ErrParse).
	ErrParse = errors.New("XUID string cannot be parsed")

	// ErrScan is returned by XUID.Scan when a database value cannot be
	// converted to an XUID. Scan wraps it with context, so failures can be
	// detected with errors.Is(err, ErrScan).
	ErrScan = errors.New("XUID cannot be scanned from a SQL value")

	// ErrInvalidPrefix is returned when a prefix does not satisfy the
	// rules enforced by ValidatePrefix. Constructors and Parse wrap it
	// with context, so failures can be detected with
	// errors.Is(err, ErrInvalidPrefix).
	ErrInvalidPrefix = errors.New("XUID prefix is invalid")

	// ErrNilUUIDWithPrefix is returned when a non-empty prefix is combined
	// with the nil UUID. A nil UUID represents an absent identifier and
	// never carries a prefix: JSON, text and SQL all encode it as
	// null/empty, so a prefixed nil UUID cannot round-trip. Constructors,
	// Parse and WithPrefix wrap it with context, so failures can be
	// detected with errors.Is(err, ErrNilUUIDWithPrefix).
	ErrNilUUIDWithPrefix = errors.New("XUID cannot combine the nil UUID with a prefix")

	// ErrNotSortable is returned by Time when the XUID does not embed a
	// UUIDv7 timestamp, such as a UUIDv4 or the nil UUID. Time wraps it
	// with context, so failures can be detected with
	// errors.Is(err, ErrNotSortable), and IsSortable can be used to check
	// first.
	ErrNotSortable = errors.New("XUID does not embed a sortable timestamp")
)
