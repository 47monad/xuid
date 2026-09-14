package xuid_test

import (
	"bytes"
	"encoding/gob"
	"testing"
	"uuid"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXUIDGobEncodeDecode(t *testing.T) {
	t.Run("round trips a prefixed XUID", func(t *testing.T) {
		original := xuid.MustNewSortable("user")

		data, err := original.GobEncode()
		require.NoError(t, err)
		assert.Equal(t, []byte(original.String()), data)

		var loaded xuid.XUID
		require.NoError(t, loaded.GobDecode(data))
		assert.True(t, original.Equal(loaded))
		assert.Equal(t, "user", loaded.GetPrefix())
	})

	t.Run("round trips an unprefixed XUID", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")

		data, err := original.GobEncode()
		require.NoError(t, err)

		var loaded xuid.XUID
		require.NoError(t, loaded.GobDecode(data))
		assert.True(t, original.Equal(loaded))
	})

	t.Run("round trips the nil UUID", func(t *testing.T) {
		var original xuid.XUID

		data, err := original.GobEncode()
		require.NoError(t, err)
		assert.Empty(t, data)

		// Start from a non-empty value to prove GobDecode resets it.
		loaded := xuid.MustNewSortable("user")
		require.NoError(t, loaded.GobDecode(data))
		assert.True(t, xuid.IsEmpty(loaded))
	})

	t.Run("returns ErrParse for invalid data", func(t *testing.T) {
		var id xuid.XUID

		err := id.GobDecode([]byte("not-an-xuid"))

		assert.ErrorIs(t, err, xuid.ErrParse)
		assert.True(t, xuid.IsEmpty(id))
	})
}

func TestXUIDGobRoundTrip(t *testing.T) {
	t.Run("encodes and decodes through encoding/gob", func(t *testing.T) {
		original := xuid.MustNewSortable("user")

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(original))

		var loaded xuid.XUID
		require.NoError(t, gob.NewDecoder(&buf).Decode(&loaded))

		assert.True(t, original.Equal(loaded))
		assert.Equal(t, "user", loaded.GetPrefix())
	})

	t.Run("round trips the nil UUID through encoding/gob", func(t *testing.T) {
		var original xuid.XUID

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(original))

		loaded := xuid.MustNewSortable("user")
		require.NoError(t, gob.NewDecoder(&buf).Decode(&loaded))

		assert.True(t, xuid.IsEmpty(loaded))
	})

	t.Run("round trips an XUID struct field", func(t *testing.T) {
		type record struct {
			ID   xuid.XUID
			Name string
		}
		original := record{ID: xuid.MustNewSortable("order"), Name: "widget"}

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(original))

		var loaded record
		require.NoError(t, gob.NewDecoder(&buf).Decode(&loaded))

		assert.True(t, original.ID.Equal(loaded.ID))
		assert.Equal(t, original.Name, loaded.Name)
	})

	t.Run("round trips a NullXUID", func(t *testing.T) {
		original := xuid.NullXUID{XUID: xuid.MustNewSortable("user"), Valid: true}

		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(original))

		var loaded xuid.NullXUID
		require.NoError(t, gob.NewDecoder(&buf).Decode(&loaded))

		assert.True(t, loaded.Valid)
		assert.True(t, original.XUID.Equal(loaded.XUID))
	})
}

func TestXUIDGobInterfaces(t *testing.T) {
	var _ gob.GobEncoder = xuid.XUID{}
	var _ gob.GobDecoder = (*xuid.XUID)(nil)
}
