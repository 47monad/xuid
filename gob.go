package xuid

// GobEncode implements the gob.GobEncoder interface.
//
// XUID stores its state in unexported fields, so encoding/gob cannot
// encode it by reflection; without these methods it fails with
// "type xuid.XUID has no exported fields". Implementing GobEncode and
// GobDecode lets XUIDs be transmitted and stored with gob, including
// their prefixes.
//
// The encoding is the canonical text form returned by MarshalText: the
// empty slice for the nil UUID, and prefix_base58 otherwise. It is
// stable for a given XUID, so encoded values remain decodable as the
// package evolves.
func (x XUID) GobEncode() ([]byte, error) {
	return x.MarshalText()
}

// GobDecode implements the gob.GobDecoder interface.
//
// It accepts exactly the bytes produced by GobEncode: an empty slice
// decodes to the nil UUID, and anything else must be a valid XUID
// string. Because the text form carries the prefix, gob round-trips an
// XUID with its prefix intact.
func (x *XUID) GobDecode(data []byte) error {
	return x.UnmarshalText(data)
}
