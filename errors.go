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
)
