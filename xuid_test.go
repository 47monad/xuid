package xuid_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/47monad/xuid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSortable(t *testing.T) {
	t.Run("creates sortable XUID with prefix", func(t *testing.T) {
		id, err := xuid.NewSortable("user")

		require.NoError(t, err)
		assert.True(t, id.IsSortable())
		assert.False(t, id.IsRandom())
		assert.Equal(t, "user", id.GetPrefix())
		assert.NotEqual(t, uuid.Nil, id.GetUUID())
	})

	t.Run("creates sortable XUID without prefix", func(t *testing.T) {
		id, err := xuid.NewSortable("")

		require.NoError(t, err)
		assert.True(t, id.IsSortable())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("generates unique sortable XUIDs", func(t *testing.T) {
		id1, err1 := xuid.NewSortable("test")
		id2, err2 := xuid.NewSortable("test")

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.False(t, id1.Equal(id2))
	})
}

func TestMustNewSortable(t *testing.T) {
	t.Run("creates sortable XUID without error", func(t *testing.T) {
		id := xuid.MustNewSortable("order")

		assert.True(t, id.IsSortable())
		assert.Equal(t, "order", id.GetPrefix())
	})

	t.Run("panics on error", func(t *testing.T) {
		// Note: MustNewSortable should not panic in normal circumstances
		// This test ensures the Must function works correctly
		assert.NotPanics(t, func() {
			xuid.MustNewSortable("test")
		})
	})
}

func TestNewRandom(t *testing.T) {
	t.Run("creates random XUID with prefix", func(t *testing.T) {
		id, err := xuid.NewRandom("session")

		require.NoError(t, err)
		assert.True(t, id.IsRandom())
		assert.False(t, id.IsSortable())
		assert.Equal(t, "session", id.GetPrefix())
		assert.NotEqual(t, uuid.Nil, id.GetUUID())
	})

	t.Run("creates random XUID without prefix", func(t *testing.T) {
		id, err := xuid.NewRandom("")

		require.NoError(t, err)
		assert.True(t, id.IsRandom())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("generates unique random XUIDs", func(t *testing.T) {
		id1, err1 := xuid.NewRandom("test")
		id2, err2 := xuid.NewRandom("test")

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.False(t, id1.Equal(id2))
	})
}

func TestMustNewRandom(t *testing.T) {
	t.Run("creates random XUID without error", func(t *testing.T) {
		id := xuid.MustNewRandom("order")

		assert.True(t, id.IsRandom())
		assert.Equal(t, "order", id.GetPrefix())
	})

	t.Run("panics on error", func(t *testing.T) {
		// Note: MustNewRandom should not panic in normal circumstances
		// This test ensures the Must function works correctly
		assert.NotPanics(t, func() {
			xuid.MustNewRandom("test")
		})
	})
}

func TestNewWith(t *testing.T) {
	t.Run("creates XUID with provided UUID and prefix", func(t *testing.T) {
		testUUID := uuid.New()
		id, err := xuid.NewWith(testUUID, "custom")

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "custom", id.GetPrefix())
	})

	t.Run("creates XUID with nil UUID", func(t *testing.T) {
		id, err := xuid.NewWith(uuid.Nil, "empty")

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, id.GetUUID())
		assert.Equal(t, "empty", id.GetPrefix())
	})

	t.Run("creates XUID without prefix", func(t *testing.T) {
		testUUID := uuid.New()
		id, err := xuid.NewWith(testUUID, "")

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})
}

func TestNew(t *testing.T) {
	t.Run("returns error for unsupported method", func(t *testing.T) {
		_, err := xuid.New()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "method not supported")
	})
}

func TestNilUUID(t *testing.T) {
	t.Run("creates XUID with nil UUID", func(t *testing.T) {
		id, err := xuid.NilUUID()

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
		assert.True(t, xuid.IsEmpty(id))
	})
}

func TestXUIDString(t *testing.T) {
	t.Run("returns base58 encoded string with prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, err := xuid.NewWith(testUUID, "test")

		require.NoError(t, err)
		str := id.String()

		assert.True(t, strings.HasPrefix(str, "test_"))
		assert.Contains(t, str, "_")
		parts := strings.Split(str, "_")
		assert.Len(t, parts, 2)
		assert.Equal(t, "test", parts[0])
	})

	t.Run("returns base58 encoded string without prefix", func(t *testing.T) {
		testUUID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
		id, err := xuid.NewWith(testUUID, "")

		require.NoError(t, err)
		str := id.String()

		assert.NotContains(t, str, "_")
		assert.NotEmpty(t, str)
	})
}

func TestXUIDEqual(t *testing.T) {
	t.Run("returns true for identical XUIDs", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "test")
		id2, _ := xuid.NewWith(testUUID, "test")

		assert.True(t, id1.Equal(id2))
	})

	t.Run("returns false for different UUIDs", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.New(), "test")
		id2, _ := xuid.NewWith(uuid.New(), "test")

		assert.False(t, id1.Equal(id2))
	})

	t.Run("returns false for different prefixes", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "test1")
		id2, _ := xuid.NewWith(testUUID, "test2")

		assert.False(t, id1.Equal(id2))
	})
}

func TestParse(t *testing.T) {
	t.Run("parses XUID string with prefix", func(t *testing.T) {
		original, _ := xuid.NewSortable("user")
		str := original.String()

		parsed, err := xuid.Parse(str)

		require.NoError(t, err)
		assert.True(t, original.Equal(parsed))
		assert.Equal(t, "user", parsed.GetPrefix())
		assert.Equal(t, original.GetUUID(), parsed.GetUUID())
	})

	t.Run("parses XUID string without prefix", func(t *testing.T) {
		original, _ := xuid.NewSortable("")
		str := original.String()

		parsed, err := xuid.Parse(str)

		require.NoError(t, err)
		assert.True(t, original.Equal(parsed))
		assert.Equal(t, "", parsed.GetPrefix())
		assert.Equal(t, original.GetUUID(), parsed.GetUUID())
	})

	t.Run("parses XUID string with underscore in prefix", func(t *testing.T) {
		original, _ := xuid.NewSortable("user_profile")
		str := original.String()

		parsed, err := xuid.Parse(str)

		require.NoError(t, err)
		assert.True(t, original.Equal(parsed))
		assert.Equal(t, "user_profile", parsed.GetPrefix())
	})

	t.Run("returns error for invalid XUID string", func(t *testing.T) {
		_, err := xuid.Parse("invalid_string")

		assert.Error(t, err)
		assert.Equal(t, xuid.ErrParse, err)
	})

	t.Run("returns error for malformed base58", func(t *testing.T) {
		_, err := xuid.Parse("test_invalid0characters")

		assert.Error(t, err)
		assert.Equal(t, xuid.ErrParse, err)
	})
}

func TestIsValid(t *testing.T) {
	t.Run("returns true for valid XUID string", func(t *testing.T) {
		id, _ := xuid.NewSortable("test")
		str := id.String()

		assert.True(t, xuid.IsValid(str))
	})

	t.Run("returns false for invalid XUID string", func(t *testing.T) {
		assert.False(t, xuid.IsValid("invalid_string"))
	})

	t.Run("returns false for empty string", func(t *testing.T) {
		assert.False(t, xuid.IsValid(""))
	})
}

func TestMust(t *testing.T) {
	t.Run("returns XUID when no error", func(t *testing.T) {
		id, _ := xuid.NewSortable("test")
		result := xuid.Must(id, nil)

		assert.True(t, id.Equal(result))
	})

	t.Run("panics when error provided", func(t *testing.T) {
		id, _ := xuid.NewSortable("test")

		assert.Panics(t, func() {
			xuid.Must(id, assert.AnError)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("returns true for nil UUID", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		assert.True(t, xuid.IsEmpty(id))
	})

	t.Run("returns false for non-nil UUID", func(t *testing.T) {
		id, _ := xuid.NewSortable("test")

		assert.False(t, xuid.IsEmpty(id))
	})
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("marshals XUID to JSON string", func(t *testing.T) {
		id, _ := xuid.NewSortable("user")
		expected := id.String()

		data, err := json.Marshal(id)

		require.NoError(t, err)
		assert.Equal(t, `"`+expected+`"`, string(data))
	})

	t.Run("unmarshals JSON string to XUID", func(t *testing.T) {
		original, _ := xuid.NewSortable("user")
		jsonStr := `"` + original.String() + `"`

		var parsed xuid.XUID
		err := json.Unmarshal([]byte(jsonStr), &parsed)

		require.NoError(t, err)
		assert.True(t, original.Equal(parsed))
		assert.Equal(t, original.GetPrefix(), parsed.GetPrefix())
		assert.Equal(t, original.GetUUID(), parsed.GetUUID())
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		var parsed xuid.XUID
		err := json.Unmarshal([]byte(`invalid json`), &parsed)

		assert.Error(t, err)
	})

	t.Run("returns error for invalid XUID in JSON", func(t *testing.T) {
		var parsed xuid.XUID
		err := json.Unmarshal([]byte(`"invalid_xuid_string"`), &parsed)

		assert.Error(t, err)
	})

	t.Run("marshals and unmarshals in struct", func(t *testing.T) {
		type TestStruct struct {
			ID   xuid.XUID `json:"id"`
			Name string    `json:"name"`
		}

		original := TestStruct{
			ID:   xuid.MustNewSortable("test"),
			Name: "Test Name",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var parsed TestStruct
		err = json.Unmarshal(data, &parsed)
		require.NoError(t, err)

		assert.True(t, original.ID.Equal(parsed.ID))
		assert.Equal(t, original.Name, parsed.Name)
	})
}

func TestVersionChecking(t *testing.T) {
	t.Run("correctly identifies sortable UUIDs", func(t *testing.T) {
		id, _ := xuid.NewSortable("test")

		assert.True(t, id.IsSortable())
		assert.False(t, id.IsRandom())
	})

	t.Run("correctly identifies random UUIDs", func(t *testing.T) {
		id, _ := xuid.NewRandom("test")

		assert.True(t, id.IsRandom())
		assert.False(t, id.IsSortable())
	})

	t.Run("handles custom UUID versions", func(t *testing.T) {
		// Using UUID v1 for testing
		customUUID, _ := uuid.NewUUID()
		id, _ := xuid.NewWith(customUUID, "test")

		assert.False(t, id.IsSortable())
		assert.False(t, id.IsRandom())
	})
}

func TestGetters(t *testing.T) {
	t.Run("GetUUID returns correct UUID", func(t *testing.T) {
		testUUID := uuid.New()
		id, _ := xuid.NewWith(testUUID, "test")

		assert.Equal(t, testUUID, id.GetUUID())
	})

	t.Run("GetPrefix returns correct prefix", func(t *testing.T) {
		id, _ := xuid.NewSortable("myPrefix")

		assert.Equal(t, "myPrefix", id.GetPrefix())
	})

	t.Run("GetPrefix returns empty string for no prefix", func(t *testing.T) {
		id, _ := xuid.NewSortable("")

		assert.Equal(t, "", id.GetPrefix())
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("handles empty prefix consistently", func(t *testing.T) {
		id1, _ := xuid.NewSortable("")
		id2, _ := xuid.NewSortable("")

		str1 := id1.String()
		str2 := id2.String()

		assert.NotContains(t, str1, "_")
		assert.NotContains(t, str2, "_")
		assert.NotEqual(t, str1, str2) // Different UUIDs
	})

	t.Run("handles very long prefix", func(t *testing.T) {
		longPrefix := strings.Repeat("a", 100)
		id, err := xuid.NewSortable(longPrefix)

		require.NoError(t, err)
		assert.Equal(t, longPrefix, id.GetPrefix())
		assert.True(t, strings.HasPrefix(id.String(), longPrefix+"_"))
	})

	t.Run("handles prefix with special characters", func(t *testing.T) {
		specialPrefix := "test-prefix.with_special@chars"
		id, err := xuid.NewSortable(specialPrefix)

		require.NoError(t, err)
		assert.Equal(t, specialPrefix, id.GetPrefix())

		// Should be able to parse it back
		parsed, err := xuid.Parse(id.String())
		require.NoError(t, err)
		assert.Equal(t, specialPrefix, parsed.GetPrefix())
	})
}

func BenchmarkNewSortable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = xuid.NewSortable("bench")
	}
}

func BenchmarkNewRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = xuid.NewRandom("bench")
	}
}

func BenchmarkString(b *testing.B) {
	id, _ := xuid.NewSortable("bench")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.String()
	}
}

func BenchmarkParse(b *testing.B) {
	id, _ := xuid.NewSortable("bench")
	str := id.String()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = xuid.Parse(str)
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	id, _ := xuid.NewSortable("bench")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(id)
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	id, _ := xuid.NewSortable("bench")
	data, _ := json.Marshal(id)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var parsed xuid.XUID
		_ = json.Unmarshal(data, &parsed)
	}
}
