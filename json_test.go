package xuid_test

import (
	"encoding/json"
	"testing"
	"uuid"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXUIDMarshalJSON(t *testing.T) {
	t.Run("marshals valid XUID to string", func(t *testing.T) {
		id := xuid.MustNewSortable("user")

		data, err := json.Marshal(id)

		require.NoError(t, err)
		assert.Equal(t, `"`+id.String()+`"`, string(data))
	})

	t.Run("marshals XUID without prefix to string", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, _ := xuid.NewWith(testUUID, "")

		data, err := json.Marshal(id)

		require.NoError(t, err)
		assert.Equal(t, `"`+id.String()+`"`, string(data))
	})

	t.Run("marshals zero-value XUID to null", func(t *testing.T) {
		var id xuid.XUID

		data, err := json.Marshal(id)

		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("marshals nil UUID XUID to null", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		data, err := json.Marshal(id)

		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("implements json.Marshaler interface", func(t *testing.T) {
		var _ json.Marshaler = xuid.XUID{}
	})
}

func TestXUIDUnmarshalJSON(t *testing.T) {
	t.Run("unmarshals prefixed XUID string", func(t *testing.T) {
		original := xuid.MustNewSortable("user")
		var id xuid.XUID

		err := json.Unmarshal([]byte(`"`+original.String()+`"`), &id)

		require.NoError(t, err)
		assert.True(t, original.Equal(id))
		assert.Equal(t, "user", id.GetPrefix())
	})

	t.Run("unmarshals XUID string without prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")
		var id xuid.XUID

		err := json.Unmarshal([]byte(`"`+original.String()+`"`), &id)

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("unmarshalling null into the zero value leaves it empty", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte("null"), &id)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("unmarshalling null leaves an existing value unchanged", func(t *testing.T) {
		original := xuid.MustNewSortable("user")
		id := original

		err := json.Unmarshal([]byte("null"), &id)

		require.NoError(t, err)
		assert.True(t, original.Equal(id))
		assert.Equal(t, "user", id.GetPrefix())
	})

	t.Run("unmarshalling null leaves an existing struct field unchanged", func(t *testing.T) {
		type user struct {
			ID xuid.XUID `json:"id"`
		}
		original := user{ID: xuid.MustNewSortable("user")}
		got := original

		err := json.Unmarshal([]byte(`{"id":null}`), &got)

		require.NoError(t, err)
		assert.True(t, original.ID.Equal(got.ID))
	})

	t.Run("unmarshalling an empty string clears the value", func(t *testing.T) {
		id := xuid.MustNewSortable("user")

		err := json.Unmarshal([]byte(`""`), &id)

		require.NoError(t, err)
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns error for invalid XUID string", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte(`"not-an-xuid"`), &id)

		assert.Error(t, err)
	})

	t.Run("rejects a prefixed nil UUID string", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte(`"user_1111111111111111"`), &id)

		assert.ErrorIs(t, err, xuid.ErrParse)
		assert.ErrorIs(t, err, xuid.ErrNilUUIDWithPrefix)
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns error for non-string JSON value", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte("123"), &id)

		assert.Error(t, err)
	})

	t.Run("implements json.Unmarshaler interface", func(t *testing.T) {
		var _ json.Unmarshaler = (*xuid.XUID)(nil)
	})
}

func TestJSONRoundTrip(t *testing.T) {
	t.Run("complete JSON round trip preserves UUID and prefix", func(t *testing.T) {
		original := xuid.MustNewSortable("user")

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var loaded xuid.XUID
		err = json.Unmarshal(data, &loaded)
		require.NoError(t, err)

		assert.True(t, original.Equal(loaded))
		assert.Equal(t, original.GetPrefix(), loaded.GetPrefix())
		assert.Equal(t, original.GetUUID(), loaded.GetUUID())
	})

	t.Run("handles nil UUID round trip", func(t *testing.T) {
		original, _ := xuid.NilUUID()

		data, err := json.Marshal(original)
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))

		var loaded xuid.XUID
		err = json.Unmarshal(data, &loaded)
		require.NoError(t, err)

		assert.True(t, xuid.IsEmpty(loaded))
	})
}
