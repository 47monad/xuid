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

	t.Run("unmarshals null to zero-value XUID", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte("null"), &id)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("unmarshals null over existing value resets to zero-value XUID", func(t *testing.T) {
		id := xuid.MustNewSortable("user")

		err := json.Unmarshal([]byte("null"), &id)

		require.NoError(t, err)
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns error for invalid XUID string", func(t *testing.T) {
		var id xuid.XUID

		err := json.Unmarshal([]byte(`"not-an-xuid"`), &id)

		assert.Error(t, err)
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
