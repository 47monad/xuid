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
	"errors"
	"strings"
	"uuid"

	"github.com/btcsuite/btcd/btcutil/base58"
)

type XUID struct {
	uuid   uuid.UUID
	prefix string
}

func New() (XUID, error) {
	var xid XUID
	return xid, errors.New("method not supported")
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

// SetPrefix sets the prefix field to the specified prefix.
// This is useful when loading XUIDs from database and need to restore the prefix.
func (x *XUID) SetPrefix(prefix string) *XUID {
	x.prefix = prefix
	return x
}

func (x XUID) String() string {
	if x.prefix == "" {
		return base58.Encode(x.uuid[:])
	}
	return x.prefix + "_" + base58.Encode(x.uuid[:])
}

func (x XUID) Equal(y XUID) bool {
	return x.String() == y.String()
}

func Parse(idstr string) (XUID, error) {
	underscoreIndex := strings.LastIndex(idstr, "_")
	uuidstr := idstr[underscoreIndex+1:]
	prefix := ""
	if underscoreIndex >= 0 {
		prefix = idstr[:underscoreIndex]
	}
	_str := base58.Decode(uuidstr)
	var _uuid uuid.UUID
	if len(_str) != len(_uuid) {
		return XUID{}, ErrParse
	}
	copy(_uuid[:], _str)
	return NewWith(_uuid, prefix)
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
