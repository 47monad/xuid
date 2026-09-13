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
    // Generate a sortable XUID with prefix
    id := xuid.MustNewSortable("user")
    fmt.Println(id.String()) // Output: user_8M7Qq2vR3kGbF9wN5pL2xA

    // Generate a random XUID
    randomID, _ := xuid.NewRandom("session")
    fmt.Println(randomID.String()) // Output: session_5K9Mm1nP7jDcE8vL3qR6yB

    // Parse an existing XUID string
    parsed, _ := xuid.Parse("user_8M7Qq2vR3kGbF9wN5pL2xA")
    fmt.Println(parsed.GetPrefix()) // Output: user
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

### Working with XUIDs

#### String Representation

```go
id := xuid.MustNewSortable("user")
fmt.Println(id.String()) // user_8M7Qq2vR3kGbF9wN5pL2xA
```

#### Access Properties

```go
id := xuid.MustNewSortable("user")

// Get the underlying UUID
uuid := id.GetUUID()

// Get the prefix
prefix := id.GetPrefix() // "user"

// Check UUID version
isSortable := id.IsSortable() // true for UUIDv7
isRandom := id.IsRandom()     // true for UUIDv4

// Extract the creation instant from a UUIDv7 (errors for other versions)
created, err := id.Time()
```

#### Parsing and Validation

```go
// Parse a XUID string
id, err := xuid.Parse("user_8M7Qq2vR3kGbF9wN5pL2xA")
if err != nil {
    log.Fatal(err)
}

// Validate a XUID string
if xuid.IsValid("user_8M7Qq2vR3kGbF9wN5pL2xA") {
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
// {"id":"user_8M7Qq2vR3kGbF9wN5pL2xA","name":"John Doe"}

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
// {"user_8M7Qq2vR3kGbF9wN5pL2xA":10}

// Or marshal/unmarshal the text form directly.
text, _ := id.MarshalText()
var loaded xuid.XUID
loaded.UnmarshalText(text)
```

A zero-value XUID marshals to an empty string, and an empty string unmarshals to
the zero value, mirroring how JSON uses `null` and SQL uses `NULL`.

### SQL Support

XUIDs integrate seamlessly with SQL databases such as PostgreSQL and MySQL. However, there are a few caveats to keep in mind:

- **Only the UUID bytes are stored** — The 16-byte UUID is stored in the database as a []byte (e.g., BYTEA in PostgreSQL or BINARY(16) in MySQL). This ensures efficient storage and indexing.
- **Prefixes are not stored** — Scanning a database value yields a XUID with an empty prefix. Restore the prefix with `WithPrefix`/`MustWithPrefix` in your repository/DAO layer, based on the table or column the value was read from:

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

`WithPrefix` validates the prefix exactly like the constructors and `Parse`, and returns `(XUID, error)`. It is immutable: it returns a copy. `MustWithPrefix` is the panic-on-error variant, so it chains off any value, including non-addressable ones such as `xuid.MustParse(s).MustWithPrefix("user")`. The older `SetPrefix` method is deprecated and does not validate.

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

- **Without prefix**: `8M7Qq2vR3kGbF9wN5pL2xA`
- **With prefix**: `prefix_8M7Qq2vR3kGbF9wN5pL2xA`

Prefixes are validated by every constructor and by `Parse`:

- At most 32 bytes long (`xuid.MaxPrefixLen`).
- Only ASCII letters, digits, and underscores: `[a-zA-Z0-9_]`. Hyphens are not allowed.
- An empty prefix means the identifier has no prefix.

Use `xuid.ValidatePrefix(s)` to check a prefix on its own without
constructing an XUID. Because `Parse` applies the same rules, `IsValid`
enforces them too.

The identifier part is a base58-encoded UUID, making it:

- **Shorter** than standard UUID strings (22 characters vs 36)
- **URL-safe** (no special characters that need encoding)
- **Case-sensitive** but avoids confusing characters (0, O, I, l)

## Error Handling

The package defines sentinel errors:

```go
var (
    ErrParse         = errors.New("XUID string cannot be parsed")
    ErrScan          = errors.New("XUID cannot be scanned from a SQL value")
    ErrInvalidPrefix = errors.New("XUID prefix is invalid")
    ErrNotSortable   = errors.New("XUID does not embed a sortable timestamp")
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
