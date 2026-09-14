package xuid_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"uuid"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXUIDValue(t *testing.T) {
	t.Run("returns raw UUID bytes for valid XUID", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, _ := xuid.NewWith(testUUID, "user")

		value, err := id.Value()

		require.NoError(t, err)
		assert.Equal(t, testUUID[:], value)
		assert.IsType(t, []byte{}, value)
		assert.Len(t, value, 16)
	})

	t.Run("returns nil for nil UUID", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		value, err := id.Value()

		require.NoError(t, err)
		assert.Nil(t, value)
	})

	t.Run("ignores prefix when converting to value", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id1, _ := xuid.NewWith(testUUID, "user")
		id2, _ := xuid.NewWith(testUUID, "different_prefix")

		value1, err1 := id1.Value()
		value2, err2 := id2.Value()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, value1, value2)
	})

	t.Run("implements driver.Valuer interface", func(t *testing.T) {
		var _ driver.Valuer = xuid.XUID{}
	})
}

func TestXUIDScan(t *testing.T) {
	t.Run("scans UUID string successfully", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		var id xuid.XUID

		err := id.Scan(testUUID.String())

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix()) // Prefix is lost when scanning
	})

	t.Run("scans UUID byte slice successfully", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		var id xuid.XUID

		err := id.Scan(testUUID[:])

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix()) // Prefix is lost when scanning
	})

	t.Run("scans UUID string delivered as byte slice", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		var id xuid.XUID

		// lib/pq returns UUID columns as a []byte holding the textual form.
		err := id.Scan([]byte(testUUID.String()))

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix()) // Prefix is lost when scanning
	})

	t.Run("scans the package's own XUID string format", func(t *testing.T) {
		original := xuid.MustNewSortable("user")
		var id xuid.XUID

		err := id.Scan(original.String())

		require.NoError(t, err)
		assert.Equal(t, original.GetUUID(), id.GetUUID())
		assert.Equal(t, "user", id.GetPrefix()) // XUID text preserves the prefix
		assert.True(t, original.Equal(id))
	})

	t.Run("scans an unprefixed XUID string", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")
		var id xuid.XUID

		err := id.Scan(original.String())

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("scans an XUID string delivered as a byte slice", func(t *testing.T) {
		original := xuid.MustNewSortable("order")
		var id xuid.XUID

		err := id.Scan([]byte(original.String()))

		require.NoError(t, err)
		assert.True(t, original.Equal(id))
	})

	t.Run("scans 16-byte XUID text instead of raw bytes", func(t *testing.T) {
		// Regression test: "1111111111111111" is the base58 string form
		// of the nil XUID and is exactly 16 bytes long. It must be parsed
		// as that (empty) identifier, not as the raw bytes 0x31..., which
		// would silently yield a different, valid-looking UUID.
		var id xuid.XUID

		err := id.Scan([]byte("1111111111111111"))

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("scans a 16-byte non-nil XUID text", func(t *testing.T) {
		// A UUID with 15 leading zero bytes encodes to the 16-character
		// XUID string "1111111111111112".
		testUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
		original, _ := xuid.NewWith(testUUID, "")
		require.Len(t, original.String(), 16, "regression fixture must be 16 bytes")

		var id xuid.XUID
		err := id.Scan([]byte(original.String()))

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("scans 16-byte array successfully", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		var raw [16]byte
		copy(raw[:], testUUID[:])
		var id xuid.XUID

		// pgx surfaces its native UUID type as [16]byte.
		err := id.Scan(raw)

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix()) // Prefix is lost when scanning
	})

	t.Run("scans uuid.UUID successfully", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		var id xuid.XUID

		err := id.Scan(testUUID)

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix()) // Prefix is lost when scanning
	})

	t.Run("scans nil value successfully", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan(nil)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns ErrScan for invalid XUID format", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan("invalid-xuid-format")

		assert.ErrorIs(t, err, xuid.ErrScan)
	})

	t.Run("returns ErrScan for invalid UUID byte length", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan([]byte{0x01, 0x02, 0x03})

		assert.ErrorIs(t, err, xuid.ErrScan)
	})

	t.Run("returns ErrScan for unsupported type", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan(123)

		assert.ErrorIs(t, err, xuid.ErrScan)
	})

	t.Run("implements sql.Scanner interface", func(t *testing.T) {
		var _ sql.Scanner = (*xuid.XUID)(nil)
	})
}

func TestNullXUID(t *testing.T) {
	t.Run("Value returns nil when invalid", func(t *testing.T) {
		var n xuid.NullXUID

		value, err := n.Value()

		require.NoError(t, err)
		assert.Nil(t, value)
	})

	t.Run("Value returns raw UUID bytes when valid", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, _ := xuid.NewWith(testUUID, "user")
		n := xuid.NullXUID{XUID: id, Valid: true}

		value, err := n.Value()

		require.NoError(t, err)
		assert.Equal(t, testUUID[:], value)
		assert.IsType(t, []byte{}, value)
		assert.Len(t, value, 16)
	})

	t.Run("Value keeps a valid nil UUID distinct from NULL", func(t *testing.T) {
		id, _ := xuid.NilUUID()
		n := xuid.NullXUID{XUID: id, Valid: true}

		value, err := n.Value()

		require.NoError(t, err)
		assert.Equal(t, make([]byte, 16), value)
	})

	t.Run("Scan of NULL clears the value and Valid", func(t *testing.T) {
		id := xuid.MustNewSortable("user")
		n := xuid.NullXUID{XUID: id, Valid: true}

		err := n.Scan(nil)

		require.NoError(t, err)
		assert.False(t, n.Valid)
		assert.True(t, xuid.IsEmpty(n.XUID))
	})

	t.Run("Scan accepts the shapes XUID.Scan accepts", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		raw := testUUID[:]
		unprefixed, _ := xuid.NewWith(testUUID, "")

		values := map[string]interface{}{
			"uuid string":   testUUID.String(),
			"xuid string":   unprefixed.String(),
			"raw bytes":     raw,
			"16-byte array": [16]byte(raw),
			"uuid.UUID":     testUUID,
		}

		for name, value := range values {
			t.Run(name, func(t *testing.T) {
				var n xuid.NullXUID

				err := n.Scan(value)

				require.NoError(t, err)
				assert.True(t, n.Valid)
				assert.Equal(t, testUUID, n.XUID.GetUUID())
				assert.Equal(t, "", n.XUID.GetPrefix())
			})
		}
	})

	t.Run("Scan returns ErrScan and stays invalid for bad input", func(t *testing.T) {
		var n xuid.NullXUID

		err := n.Scan("not-a-uuid")

		assert.ErrorIs(t, err, xuid.ErrScan)
		assert.False(t, n.Valid)
		assert.True(t, xuid.IsEmpty(n.XUID))
	})

	t.Run("round trips NULL", func(t *testing.T) {
		var original xuid.NullXUID

		value, err := original.Value()
		require.NoError(t, err)
		assert.Nil(t, value)

		var loaded xuid.NullXUID
		require.NoError(t, loaded.Scan(value))
		assert.False(t, loaded.Valid)
	})

	t.Run("round trips a value", func(t *testing.T) {
		original := xuid.NullXUID{XUID: xuid.MustNewSortable("user"), Valid: true}

		value, err := original.Value()
		require.NoError(t, err)

		var loaded xuid.NullXUID
		require.NoError(t, loaded.Scan(value))

		assert.True(t, loaded.Valid)
		assert.True(t, original.XUID.EqualUUID(loaded.XUID))
	})

	t.Run("implements driver.Valuer and sql.Scanner", func(t *testing.T) {
		var _ driver.Valuer = xuid.NullXUID{}
		var _ sql.Scanner = (*xuid.NullXUID)(nil)
	})
}

func TestXUIDWithPrefix(t *testing.T) {
	t.Run("adds prefix to existing XUID", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")

		withPrefix := original.MustWithPrefix("user")

		assert.Equal(t, testUUID, withPrefix.GetUUID())
		assert.Equal(t, "user", withPrefix.GetPrefix())
	})

	t.Run("replaces existing prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "old_prefix")

		withNewPrefix := original.MustWithPrefix("new_prefix")

		assert.Equal(t, testUUID, withNewPrefix.GetUUID())
		assert.Equal(t, "new_prefix", withNewPrefix.GetPrefix())
	})

	t.Run("handles empty prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "old_prefix")

		withoutPrefix := original.MustWithPrefix("")

		assert.Equal(t, testUUID, withoutPrefix.GetUUID())
		assert.Equal(t, "", withoutPrefix.GetPrefix())
	})
}

func TestSQLRoundTrip(t *testing.T) {
	t.Run("complete SQL round trip preserves UUID", func(t *testing.T) {
		original, _ := xuid.NewSortable("user")

		// Simulate saving to database
		value, err := original.Value()
		require.NoError(t, err)

		// Simulate loading from database
		var loaded xuid.XUID
		err = loaded.Scan(value)
		require.NoError(t, err)

		// UUID should be preserved
		assert.Equal(t, original.GetUUID(), loaded.GetUUID())

		// Prefix is lost (expected behavior)
		assert.Equal(t, "", loaded.GetPrefix())

		// Restore prefix after loading
		restored := loaded.MustWithPrefix("user")
		assert.Equal(t, original.GetUUID(), restored.GetUUID())
		assert.Equal(t, original.GetPrefix(), restored.GetPrefix())
		assert.True(t, original.Equal(restored))
	})

	t.Run("handles nil UUID round trip", func(t *testing.T) {
		original, _ := xuid.NilUUID()

		// Simulate saving to database
		value, err := original.Value()
		require.NoError(t, err)
		assert.Nil(t, value)

		// Simulate loading from database
		var loaded xuid.XUID
		err = loaded.Scan(value)
		require.NoError(t, err)

		assert.Equal(t, original.GetUUID(), loaded.GetUUID())
		assert.True(t, xuid.IsEmpty(loaded))
	})
}

// Example usage patterns for SQL integration
func TestSQLUsagePatterns(t *testing.T) {
	t.Run("demonstrates typical database workflow", func(t *testing.T) {
		// 1. Create XUID in application
		userID := xuid.MustNewSortable("user")
		orderID := xuid.MustNewSortable("order")

		// 2. Store in database (only UUID is stored)
		userValue, _ := userID.Value()
		orderValue, _ := orderID.Value()

		// These would be stored as binary UUID columns in PostgreSQL
		assert.IsType(t, []byte{}, userValue)
		assert.IsType(t, []byte{}, orderValue)

		// 3. Load from database
		var loadedUserID, loadedOrderID xuid.XUID
		require.NoError(t, loadedUserID.Scan(userValue))
		require.NoError(t, loadedOrderID.Scan(orderValue))

		// 4. Restore prefixes based on context (table/column knowledge)
		restoredUserID := loadedUserID.MustWithPrefix("user")
		restoredOrderID := loadedOrderID.MustWithPrefix("order")

		// 5. Verify restoration
		assert.True(t, userID.Equal(restoredUserID))
		assert.True(t, orderID.Equal(restoredOrderID))
	})

	t.Run("demonstrates prefix restoration strategies", func(t *testing.T) {
		// Strategy 1: Restore based on table/column context
		userID := xuid.MustNewSortable("user")
		value, _ := userID.Value()

		var loaded xuid.XUID
		require.NoError(t, loaded.Scan(value))

		// In your repository/DAO layer:
		restoredFromUsersTable := loaded.MustWithPrefix("user")
		assert.Equal(t, "user", restoredFromUsersTable.GetPrefix())

		// Strategy 2: Store prefix separately if needed
		prefix := userID.GetPrefix()
		assert.Equal(t, "user", prefix)
		// You could store this in a separate column if prefix variety is needed

		// Strategy 3: Use WithPrefix("") for prefix-agnostic operations
		withoutPrefix := userID.MustWithPrefix("")
		assert.Equal(t, "", withoutPrefix.GetPrefix())
		assert.Equal(t, userID.GetUUID(), withoutPrefix.GetUUID())
	})
}

func BenchmarkSQLOperations(b *testing.B) {
	id := xuid.MustNewSortable("bench")

	b.Run("Value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = id.Value()
		}
	})

	b.Run("Scan", func(b *testing.B) {
		value, _ := id.Value()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var xid xuid.XUID
			_ = xid.Scan(value)
		}
	})
}
