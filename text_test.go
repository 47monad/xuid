package xuid_test

import (
	"encoding"
	"encoding/json"
	"testing"
	"uuid"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXUIDMarshalText(t *testing.T) {
	t.Run("marshals valid XUID to canonical string", func(t *testing.T) {
		id := xuid.MustNewSortable("user")

		data, err := id.MarshalText()

		require.NoError(t, err)
		assert.Equal(t, id.String(), string(data))
	})

	t.Run("marshals XUID without prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, _ := xuid.NewWith(testUUID, "")

		data, err := id.MarshalText()

		require.NoError(t, err)
		assert.Equal(t, id.String(), string(data))
	})

	t.Run("marshals zero-value XUID to empty text", func(t *testing.T) {
		var id xuid.XUID

		data, err := id.MarshalText()

		require.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("marshals nil UUID XUID to empty text", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		data, err := id.MarshalText()

		require.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("implements encoding.TextMarshaler interface", func(t *testing.T) {
		var _ encoding.TextMarshaler = xuid.XUID{}
	})
}

func TestXUIDUnmarshalText(t *testing.T) {
	t.Run("unmarshals prefixed XUID string", func(t *testing.T) {
		original := xuid.MustNewSortable("user")
		var id xuid.XUID

		err := id.UnmarshalText([]byte(original.String()))

		require.NoError(t, err)
		assert.True(t, original.Equal(id))
		assert.Equal(t, "user", id.GetPrefix())
	})

	t.Run("unmarshals XUID string without prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")
		var id xuid.XUID

		err := id.UnmarshalText([]byte(original.String()))

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("unmarshals empty text to zero-value XUID", func(t *testing.T) {
		var id xuid.XUID

		err := id.UnmarshalText(nil)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("unmarshals empty text over existing value resets to zero-value XUID", func(t *testing.T) {
		id := xuid.MustNewSortable("user")

		err := id.UnmarshalText([]byte{})

		require.NoError(t, err)
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns ErrParse for invalid XUID string", func(t *testing.T) {
		var id xuid.XUID

		err := id.UnmarshalText([]byte("not-an-xuid"))

		require.Error(t, err)
		assert.ErrorIs(t, err, xuid.ErrParse)
	})

	t.Run("implements encoding.TextUnmarshaler interface", func(t *testing.T) {
		var _ encoding.TextUnmarshaler = (*xuid.XUID)(nil)
	})
}

func TestTextRoundTrip(t *testing.T) {
	t.Run("text round trip preserves UUID and prefix", func(t *testing.T) {
		original := xuid.MustNewSortable("user")

		data, err := original.MarshalText()
		require.NoError(t, err)

		var loaded xuid.XUID
		err = loaded.UnmarshalText(data)
		require.NoError(t, err)

		assert.True(t, original.Equal(loaded))
	})

	t.Run("nil UUID text round trip stays empty", func(t *testing.T) {
		original, _ := xuid.NilUUID()

		data, err := original.MarshalText()
		require.NoError(t, err)
		assert.Empty(t, data)

		var loaded xuid.XUID
		err = loaded.UnmarshalText(data)
		require.NoError(t, err)

		assert.True(t, xuid.IsEmpty(loaded))
	})
}

func TestXUIDJSONMapKey(t *testing.T) {
	t.Run("marshals XUID map keys via TextMarshaler", func(t *testing.T) {
		id := xuid.MustNewSortable("user")
		input := map[xuid.XUID]string{id: "value"}

		data, err := json.Marshal(input)
		require.NoError(t, err)
		assert.JSONEq(t, `{"`+id.String()+`":"value"}`, string(data))

		var output map[xuid.XUID]string
		err = json.Unmarshal(data, &output)
		require.NoError(t, err)
		assert.Equal(t, input, output)
	})

	t.Run("marshals nil UUID map key as empty string", func(t *testing.T) {
		var id xuid.XUID
		input := map[xuid.XUID]string{id: "value"}

		data, err := json.Marshal(input)
		require.NoError(t, err)
		assert.JSONEq(t, `{"":"value"}`, string(data))

		var output map[xuid.XUID]string
		err = json.Unmarshal(data, &output)
		require.NoError(t, err)
		assert.Equal(t, input, output)
	})
}
