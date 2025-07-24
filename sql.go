package xuid

import (
	"database/sql/driver"
	"errors"

	"github.com/google/uuid"
)

// Value implements the driver.Valuer interface.
// This allows XUID to be stored in SQL databases as UUID.
func (x XUID) Value() (driver.Value, error) {
	if x.uuid == uuid.Nil {
		return nil, nil
	}
	return x.uuid[:], nil
}

// Scan implements the sql.Scanner interface.
// This allows XUID to be loaded from SQL databases.
// Note: The prefix information is lost when loading from database.
// You should reconstruct XUIDs with their appropriate prefixes after loading.
func (x *XUID) Scan(value interface{}) error {
	if value == nil {
		x.uuid = uuid.Nil
		x.prefix = ""
		return nil
	}

	b, ok := value.([]byte)
	if !ok || len(b) != 16 {
		return errors.New("failed to scan from database. Invalid XUID bytes")
	}
	copy(x.uuid[:], b)
	x.prefix = "" // Prefix is lost when loading from database
	return nil
}
