# XUID

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/47monad/xuid.svg)](https://pkg.go.dev/github.com/47monad/xuid)
![Status: Under Development](https://img.shields.io/badge/Status-Under%20Development-orange)

A Go package for generating compact, sortable UUID-based identifiers with optional string prefixes and base58 encoding.

## Features

- 🎯 **Sortable UUIDs**: Generate UUIDv7 identifiers that maintain chronological order
- 🎲 **Random UUIDs**: Generate UUIDv4 identifiers for non-sortable use cases
- 🏷️ **Optional Prefixes**: Add human-readable prefixes to your identifiers (e.g., `user_`, `order_`)
- 📦 **Compact Encoding**: Uses base58 encoding for shorter, URL-safe strings
- 🔄 **JSON Support**: Built-in JSON marshaling and unmarshaling
- 🔤 **Text Support**: Implements `encoding.TextMarshaler`/`TextUnmarshaler` for YAML, TOML, XML and map keys
- 🗄️ **SQL Database Support**: Seamless integration with SQL databases (PostgreSQL, MySQL, etc.)
- ✅ **Type Safety**: Strong typing with validation and parsing utilities

## Installation

```bash
go get github.com/47monad/xuid
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/47monad/xuid"
)

func main() {
    // Generate a sortable XUID with prefix. The UUID is random, so the
    // exact output varies; it is always "user_" followed by at most 22
    // base58 characters.
    id := xuid.MustNewSortable("user")
    fmt.Println(id.String()) // e.g. user_Cf1k9VmUZGg55baoJFXnT

    // Generate a random XUID.
    randomID, _ := xuid.NewRandom("session")
    fmt.Println(randomID.String()) // e.g. session_8QMBv8hxcm3BpPJ9wYgzNt

    // Parse an existing XUID string.
    parsed, _ := xuid.Parse("user_Cf1k9VmUZGg55baoJFXnT")
    fmt.Println(parsed.GetPrefix()) // user
}
```

## Usage

### Creating XUIDs

#### Sortable UUIDs (UUIDv7)

```go
// With error handling
id, err := xuid.NewSortable("order")
if err != nil {
    log.Fatal(err)
}

// Without error handling (panics on error)
id := xuid.MustNewSortable("order")
```

#### Random UUIDs (UUIDv4)

```go
id, err := xuid.NewRandom("session")
if err != nil {
    log.Fatal(err)
}

// Without error handling (panics on error)
id := xuid.MustNewRandom("order")
```

#### From Existing UUID

```go
existingUUID := uuid.New()
id, err := xuid.NewWith(existingUUID, "custom")
```

> Requires Go 1.27+ for the standard library `uuid` package.

#### Nil UUID

```go
nilID, err := xuid.NilUUID()
```

A nil UUID represents an absent identifier and never carries a prefix. `NewWith`
rejects a non-empty prefix on the nil UUID (wrapping `ErrNilUUIDWithPrefix`),
and `Parse` rejects strings such as `user_1111111111111111` whose decoded UUID
is the nil UUID. This keeps the nil UUID encodable as `null`/empty by JSON, text
and SQL without losing information.

### Working with XUIDs

#### String Representation

```go
id := xuid.MustNewSortable("user")
fmt.Println(id.String()) // e.g. user_Cf1k9VmUZGg55baoJFXnT
```

#### Access Properties

```go
id := xuid.MustNewSortable("user")

// Get the underlying UUID
uuid := id.GetUUID()

// Get the prefix
prefix := id.GetPrefix() // "user"

// Check UUID version and variant (both must be RFC 9562)
isSortable := id.IsSortable() // true for well-formed UUIDv7
isRandom := id.IsRandom()     // true for well-formed UUIDv4

// Extract the creation instant from a UUIDv7 (errors for anything else)
created, err := id.Time()
```

#### Parsing and Validation

```go
// Parse a XUID string
id, err := xuid.Parse("user_Cf1k9VmUZGg55baoJFXnT")
if err != nil {
    log.Fatal(err)
}

// Validate a XUID string
if xuid.IsValid("user_Cf1k9VmUZGg55baoJFXnT") {
    fmt.Println("Valid XUID")
}

// Validate a prefix on its own (used by the constructors and Parse)
if err := xuid.ValidatePrefix("user"); err != nil {
    fmt.Println("Invalid prefix:", err)
}

// Check if empty
if xuid.IsEmpty(id) {
    fmt.Println("Empty XUID")
}
```

#### Comparison

```go
id1 := xuid.MustNewSortable("user")
id2 := xuid.MustNewSortable("user")

// Equal compares both the UUID and the prefix: a different prefix
// means the identifiers are not equal, even with the same UUID.
if id1.Equal(id2) {
    fmt.Println("XUIDs are equal")
}

// EqualUUID compares only the underlying UUIDs, ignoring prefixes.
if id1.EqualUUID(id2) {
    fmt.Println("Same identifier, prefix ignored")
}

// Compare orders XUIDs by prefix first (empty prefix sorts first),
// then by UUID bytes. For UUIDv7 IDs sharing a prefix, this is
// chronological order.
if xuid.Compare(id1, id2) < 0 {
    fmt.Println("id1 sorts before id2")
}

// Sort a slice of XUIDs.
ids := []xuid.XUID{id2, id1}
slices.SortFunc(ids, xuid.Compare)
```

### JSON Support

XUIDs can be seamlessly marshaled to and from JSON:

```go
type User struct {
    ID   xuid.XUID `json:"id"`
    Name string    `json:"name"`
}

user := User{
    ID:   xuid.MustNewSortable("user"),
    Name: "John Doe",
}

// Marshal to JSON
data, _ := json.Marshal(user)
// {"id":"user_Cf1k9VmUZGg55baoJFXnT","name":"John Doe"}

// Unmarshal from JSON
var parsed User
json.Unmarshal(data, &parsed)
```

### Text Support

`XUID` implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler`, so it
also works with text-based encoders (YAML, TOML, `encoding/xml`), templates, and
as a JSON map key:

```go
id := xuid.MustNewSortable("user")

// Use as a JSON map key.
scores := map[xuid.XUID]int{id: 10}
data, _ := json.Marshal(scores)
// {"user_Cf1k9VmUZGg55baoJFXnT":10}

// Or marshal/unmarshal the text form directly.
text, _ := id.MarshalText()
var loaded xuid.XUID
loaded.UnmarshalText(text)
```

A zero-value XUID marshals to an empty string, and an empty string unmarshals to
the zero value, mirroring how JSON uses `null` and SQL uses `NULL`.

### SQL Support

XUIDs integrate seamlessly with SQL databases such as PostgreSQL and MySQL. However, there are a few caveats to keep in mind:

- **Only the UUID bytes are stored** — `Value` writes the 16-byte UUID as a `[]byte` (e.g., BYTEA in PostgreSQL or BINARY(16) in MySQL). This ensures efficient storage and indexing.
- **Prefixes are not stored** — A binary or UUID-format value carries no prefix, so scanning one yields a XUID with an empty prefix. Restore the prefix with `WithPrefix`/`MustWithPrefix` in your repository/DAO layer, based on the table or column the value was read from:

```go
// Load from database
var loaded xuid.XUID
loaded.Scan(value)

// Restore prefix (the repo layer knows this column is a user ID)
restored := loaded.MustWithPrefix("user")

// WithPrefix validates and returns an error instead of panicking:
if restored, err := loaded.WithPrefix("user"); err != nil {
    // invalid prefix
}
```

`Scan` accepts the shapes drivers deliver UUID columns in:

- a `string` in UUID format or in the package's own XUID format (`user_Cf1k9VmUZGg55baoJFXnT`),
- a `[]byte` holding either of those textual forms (e.g., lib/pq),
- a `[]byte` of exactly 16 raw UUID bytes,
- a `[16]byte` or `uuid.UUID` (e.g., pgx's native UUID type).

A `[]byte` of exactly 16 bytes is ambiguous: it can be the raw UUID bytes written by `Value` or a 16-character textual identifier — most notably `"1111111111111111"`, the canonical string form of the nil UUID. `Scan` prefers the text interpretation when the bytes are valid UUID/XUID text and otherwise treats them as raw UUID bytes, so a 16-character text column is never silently misread. The reverse case — a raw 16-byte UUID whose bytes happen to form valid XUID text — is astronomically rare; pass it as a `[16]byte` or `uuid.UUID` to force the raw interpretation.

Scanning an XUID-format value (or a `[]byte` holding one) preserves the prefix it encodes. UUID-format and binary values cannot: the prefix is not stored in a binary column.

`WithPrefix` validates the prefix exactly like the constructors and `Parse`, and returns `(XUID, error)`. It is immutable: it returns a copy. Like the constructors, it rejects a non-empty prefix on the nil UUID (`ErrNilUUIDWithPrefix`). `MustWithPrefix` is the panic-on-error variant, so it chains off any value, including non-addressable ones such as `xuid.MustParse(s).MustWithPrefix("user")`. The older `SetPrefix` method is deprecated and does not validate (and does not enforce the nil-UUID/prefix invariant).

#### Nullable Columns

Because `XUID.Value` maps a nil UUID to SQL `NULL`, a plain `XUID` cannot
distinguish a `NULL` column from an all-zero UUID. Use `NullXUID` for
nullable columns:

```go
type User struct {
    ID xuid.NullXUID `db:"id"`
}

// Scan (NULL sets Valid to false)
var id xuid.NullXUID
if err := id.Scan(dbValue); err != nil {
    log.Fatal(err)
}
if id.Valid {
    fmt.Println("user ID:", id.XUID)
}

// Value (NULL when Valid is false)
value, err := id.Value()
```

`NullXUID` implements `driver.Valuer` and `sql.Scanner`, and mirrors the
standard library's `sql.Null*` types. `Valid` is authoritative: a
`NullXUID` with `Valid` set is stored as non-`NULL`, even when it holds the
nil UUID. Assign its `XUID` field to get at the underlying identifier;
prefixes are still lost on scan and can be restored with `WithPrefix`.

## Format

XUIDs follow this format:

- **Without prefix**: `Cf1k9VmUZGg55baoJFXnT`
- **With prefix**: `prefix_Cf1k9VmUZGg55baoJFXnT`

Prefixes are validated by every constructor and by `Parse`:

- At most 32 bytes long (`xuid.MaxPrefixLen`).
- Only ASCII letters, digits, and underscores: `[a-zA-Z0-9_]`. Hyphens are not allowed.
- An empty prefix means the identifier has no prefix.
- A non-empty prefix is only allowed on a non-nil UUID. The nil UUID (the empty
  XUID) always has an empty prefix, so `NewWith(uuid.Nil(), "user")` and
  `Parse("user_1111111111111111")` are rejected.

Use `xuid.ValidatePrefix(s)` to check a prefix on its own without
constructing an XUID. Because `Parse` applies the same rules, `IsValid`
enforces them too.

The identifier part is a base58-encoded UUID, making it:

- **Shorter** than standard UUID strings (at most 22 characters vs 36)
- **URL-safe** (no special characters that need encoding)
- **Case-sensitive** but avoids confusing characters (0, O, I, l)

The base58 encoding is not zero-padded, so its length varies: a UUIDv7
generated today renders in 21 characters, while a random UUIDv4 usually
renders in 22.

## Error Handling

The package defines sentinel errors:

```go
var (
    ErrParse             = errors.New("XUID string cannot be parsed")
    ErrScan              = errors.New("XUID cannot be scanned from a SQL value")
    ErrInvalidPrefix     = errors.New("XUID prefix is invalid")
    ErrNilUUIDWithPrefix = errors.New("XUID cannot combine the nil UUID with a prefix")
    ErrNotSortable       = errors.New("XUID does not embed a sortable timestamp")
)
```

`Parse` wraps `ErrParse` with the underlying cause for malformed XUID
strings, so failures can be detected with `errors.Is` while still carrying
a message that explains what went wrong:

```go
if _, err := xuid.Parse(s); errors.Is(err, xuid.ErrParse) {
    // s was not a valid XUID
}
```

`Scan` likewise wraps `ErrScan` with context:

```go
var id xuid.XUID
if err := id.Scan(value); errors.Is(err, xuid.ErrScan) {
    // the database value was not a valid UUID
}
```

## Dependencies

- UUID generation uses the standard library `uuid` package (Go 1.27+)
- Base58 encoding is implemented internally using the Bitcoin alphabet, which excludes the visually ambiguous characters `0`, `O`, `I`, and `l` for readability — there is no external base58 dependency

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
