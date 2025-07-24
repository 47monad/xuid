# XUID

A Go package for generating compact, sortable UUID-based identifiers with optional string prefixes and base58 encoding.

## Features

- 🎯 **Sortable UUIDs**: Generate UUIDv7 identifiers that maintain chronological order
- 🎲 **Random UUIDs**: Generate UUIDv4 identifiers for non-sortable use cases
- 🏷️ **Optional Prefixes**: Add human-readable prefixes to your identifiers (e.g., `user_`, `order_`)
- 📦 **Compact Encoding**: Uses base58 encoding for shorter, URL-safe strings
- 🔄 **JSON Support**: Built-in JSON marshaling and unmarshaling
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
import "github.com/google/uuid"

existingUUID := uuid.New()
id, err := xuid.NewWith(existingUUID, "custom")
```

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

// Check if empty
if xuid.IsEmpty(id) {
    fmt.Println("Empty XUID")
}
```

#### Comparison

```go
id1 := xuid.MustNewSortable("user")
id2 := xuid.MustNewSortable("user")

if id1.Equal(id2) {
    fmt.Println("XUIDs are equal")
}
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

### SQL Support

XUIDs integrate seamlessly with SQL databases such as PostgreSQL and MySQL. However, there are a few caveats to keep in mind:

- **Only the UUID bytes are stored** — The 16-byte UUID is stored in the database as a []byte (e.g., BYTEA in PostgreSQL or BINARY(16) in MySQL). This ensures efficient storage and indexing.
- **Prefixes are not stored** — If your application relies on the XUID prefix (e.g., "file_", "user_") for querying or categorization, you’ll need to store the prefix in a separate column.

## Format

XUIDs follow this format:

- **Without prefix**: `8M7Qq2vR3kGbF9wN5pL2xA`
- **With prefix**: `prefix_8M7Qq2vR3kGbF9wN5pL2xA`

The identifier part is a base58-encoded UUID, making it:

- **Shorter** than standard UUID strings (22 characters vs 36)
- **URL-safe** (no special characters that need encoding)
- **Case-sensitive** but avoids confusing characters (0, O, I, l)

## Error Handling

The package defines specific error types:

```go
var (
    ErrInvalidUUIDString = errors.New("UUID string is invalid")
    ErrParse             = errors.New("XUID string cannot be parsed")
)
```

## Dependencies

- `github.com/google/uuid` - UUID generation and manipulation
- `github.com/btcsuite/btcd/btcutil/base58` - Base58 encoding

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
