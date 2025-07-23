package xuid

import (
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
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
	id, err := uuid.NewV7()
	if err != nil {
		return XUID{}, err
	}
	return XUID{
		uuid:   id,
		prefix: prefix,
	}, nil
}

func MustNewSortable(prefix string) XUID {
	return Must(NewSortable(prefix))
}

func NewRandom(prefix string) (XUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return XUID{}, err
	}
	return XUID{
		uuid:   id,
		prefix: prefix,
	}, nil
}

func NilUUID() (XUID, error) {
	return NewWith(uuid.Nil, "")
}

func (x XUID) GetUUID() uuid.UUID {
	return x.uuid
}

func (x XUID) IsSortable() bool {
	return x.GetUUID().Version().String() == "VERSION_7"
}

func (x XUID) IsRandom() bool {
	return x.GetUUID().Version().String() == "VERSION_4"
}

func (x XUID) GetPrefix() string {
	return x.prefix
}

func (x XUID) String() string {
	b, _ := x.uuid.MarshalBinary()
	if x.prefix == "" {
		return base58.Encode(b)
	}
	return x.prefix + "_" + base58.Encode(b)
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
	_uuid, err := uuid.FromBytes(_str)
	if err != nil {
		return XUID{}, ErrParse
	}
	return NewWith(_uuid, prefix)
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
	return xid.uuid == uuid.Nil
}
