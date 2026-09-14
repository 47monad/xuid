# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-09-14

Initial release. Nothing was previously tagged, so the entire public API is
new. Because the module was still under development, the API reflects a number
of deliberate design decisions (prefix validation, the nil-UUID/prefix
invariant, base58 storage format, and SQL `Value` returning raw bytes). See the
notes below and the README for the behavioral contracts.

### Added

- Sortable identifiers built on UUIDv7 (`NewSortable`, `MustNewSortable`) for
  chronological ordering, and random UUIDv4 identifiers (`NewRandom`,
  `MustNewRandom`).
- Optional string prefixes with validation (`ValidatePrefix`, `MaxPrefixLen`),
  including the immutable `WithPrefix`/`MustWithPrefix` and `GetPrefix`. A
  non-empty prefix on the nil UUID is rejected (`ErrNilUUIDWithPrefix`).
- Compact, URL-safe Bitcoin base58 encoding implemented internally (no external
  base58 dependency), with `Parse`, `MustParse` and `IsValid`.
- Comparison helpers: `Equal`, `EqualUUID`, `Compare` and `Less`.
- `Time` to recover the creation instant from a UUIDv7, and `IsSortable` /
  `IsRandom` version-and-variant checks (`ErrNotSortable`).
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler` support for text
  encoders, templates and JSON map keys.
- JSON support: `MarshalJSON` / `UnmarshalJSON`, with the nil UUID encoded as
  `null`.
- gob support: `GobEncode` / `GobDecode`.
- SQL support: `Value` / `Scan` for `driver.Valuer` / `sql.Scanner`, plus
  `NullXUID` for nullable columns.
- Sentinel errors (`ErrParse`, `ErrScan`, `ErrInvalidPrefix`,
  `ErrNilUUIDWithPrefix`, `ErrNotSortable`) that can be matched with
  `errors.Is`.
- CI (build, vet, race tests, fuzz smoke tests, lint) and fuzz targets for
  parsing and base58 round trips.

### Notes

- Requires Go 1.27 or newer: UUID generation uses the standard library `uuid`
  package, which was added in Go 1.27.
- The nil UUID (the zero value) renders as the empty string via `String` and
  `MarshalText`, and as `null` / `NULL` in JSON and SQL. `Parse` rejects the
  empty string by design; use `IsEmpty` to detect it, or
  `MarshalText`/`UnmarshalText` to round-trip it.
- `XUID.Value` stores only the 16 raw UUID bytes for SQL. Prefixes are not
  persisted and must be restored per column with `WithPrefix`.
- `SetPrefix` is deprecated in favour of `WithPrefix` / `MustWithPrefix`.

[Unreleased]: https://github.com/47monad/xuid/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/47monad/xuid/releases/tag/v0.1.0
