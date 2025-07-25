package xuid_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/47monad/xuid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXUIDValue(t *testing.T) {
	t.Run("returns UUID string for valid XUID", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, _ := xuid.NewWith(testUUID, "user")

		value, err := id.Value()

		require.NoError(t, err)
		assert.Equal(t, testUUID.String(), value)
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

	t.Run("scans nil value successfully", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan(nil)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns error for invalid XUID format", func(t *testing.T) {
		var id xuid.XUID

		err := id.Scan("invalid-xuid-format")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid XUID string")
	})

	t.Run("implements sql.Scanner interface", func(t *testing.T) {
		var _ sql.Scanner = (*xuid.XUID)(nil)
	})
}

func TestXUIDSetPrefix(t *testing.T) {
	t.Run("adds prefix to existing XUID", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "")

		withPrefix := original.SetPrefix("user")

		assert.Equal(t, testUUID, withPrefix.GetUUID())
		assert.Equal(t, "user", withPrefix.GetPrefix())
	})

	t.Run("replaces existing prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "old_prefix")

		withNewPrefix := original.SetPrefix("new_prefix")

		assert.Equal(t, testUUID, withNewPrefix.GetUUID())
		assert.Equal(t, "new_prefix", withNewPrefix.GetPrefix())
	})

	t.Run("handles empty prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		original, _ := xuid.NewWith(testUUID, "old_prefix")

		withoutPrefix := original.SetPrefix("")

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
		restored := loaded.SetPrefix("user")
		assert.Equal(t, original.GetUUID(), restored.GetUUID())
		assert.Equal(t, original.GetPrefix(), restored.GetPrefix())
		assert.True(t, original.Equal(*restored))
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

		// These would be stored as UUID columns in PostgreSQL
		assert.IsType(t, "", userValue)
		assert.IsType(t, "", orderValue)

		// 3. Load from database
		var loadedUserID, loadedOrderID xuid.XUID
		loadedUserID.Scan(userValue)
		loadedOrderID.Scan(orderValue)

		// 4. Restore prefixes based on context (table/column knowledge)
		restoredUserID := loadedUserID.SetPrefix("user")
		restoredOrderID := loadedOrderID.SetPrefix("order")

		// 5. Verify restoration
		assert.True(t, userID.Equal(*restoredUserID))
		assert.True(t, orderID.Equal(*restoredOrderID))
	})

	t.Run("demonstrates prefix restoration strategies", func(t *testing.T) {
		// Strategy 1: Restore based on table/column context
		userID := xuid.MustNewSortable("user")
		value, _ := userID.Value()

		var loaded xuid.XUID
		loaded.Scan(value)

		// In your repository/DAO layer:
		restoredFromUsersTable := loaded.SetPrefix("user")
		assert.Equal(t, "user", restoredFromUsersTable.GetPrefix())

		// Strategy 2: Store prefix separately if needed
		prefix := userID.GetPrefix()
		assert.Equal(t, "user", prefix)
		// You could store this in a separate column if prefix variety is needed

		// Strategy 3: Use WithoutPrefix for prefix-agnostic operations
		withoutPrefix := userID.SetPrefix("")
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
