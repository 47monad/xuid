package xuid_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/47monad/xuid"
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
		assert.NotEqual(t, uuid.Nil(), id.GetUUID())
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
}

func TestNewRandom(t *testing.T) {
	t.Run("creates random XUID with prefix", func(t *testing.T) {
		id, err := xuid.NewRandom("session")

		require.NoError(t, err)
		assert.True(t, id.IsRandom())
		assert.False(t, id.IsSortable())
		assert.Equal(t, "session", id.GetPrefix())
		assert.NotEqual(t, uuid.Nil(), id.GetUUID())
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
}

func TestNewWith(t *testing.T) {
	t.Run("creates XUID with provided UUID and prefix", func(t *testing.T) {
		testUUID := uuid.New()
		id, err := xuid.NewWith(testUUID, "custom")

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "custom", id.GetPrefix())
	})

	t.Run("creates XUID with nil UUID and empty prefix", func(t *testing.T) {
		id, err := xuid.NewWith(uuid.Nil(), "")

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})

	t.Run("rejects nil UUID with a prefix", func(t *testing.T) {
		id, err := xuid.NewWith(uuid.Nil(), "empty")

		assert.ErrorIs(t, err, xuid.ErrNilUUIDWithPrefix)
		assert.Equal(t, xuid.XUID{}, id)
	})

	t.Run("creates XUID without prefix", func(t *testing.T) {
		testUUID := uuid.New()
		id, err := xuid.NewWith(testUUID, "")

		require.NoError(t, err)
		assert.Equal(t, testUUID, id.GetUUID())
		assert.Equal(t, "", id.GetPrefix())
	})
}

func TestNilUUID(t *testing.T) {
	t.Run("creates XUID with nil UUID", func(t *testing.T) {
		id, err := xuid.NilUUID()

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil(), id.GetUUID())
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

	t.Run("returns the empty string for the nil UUID", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		assert.Equal(t, "", id.String())
	})

	t.Run("returns the empty string for the zero value", func(t *testing.T) {
		var id xuid.XUID

		assert.Equal(t, "", id.String())
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

	t.Run("does not depend on string encoding", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "")
		id2, _ := xuid.NewWith(testUUID, "")
		require.Equal(t, id1.String(), id2.String())
		assert.True(t, id1.Equal(id2))
	})
}

func TestXUIDEqualUUID(t *testing.T) {
	t.Run("returns true for same UUID with different prefixes", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "user")
		id2, _ := xuid.NewWith(testUUID, "order")

		assert.True(t, id1.EqualUUID(id2))
		assert.False(t, id1.Equal(id2))
	})

	t.Run("returns true for same UUID ignoring empty prefix", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "")
		id2, _ := xuid.NewWith(testUUID, "user")

		assert.True(t, id1.EqualUUID(id2))
	})

	t.Run("returns false for different UUIDs", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.New(), "test")
		id2, _ := xuid.NewWith(uuid.New(), "test")

		assert.False(t, id1.EqualUUID(id2))
	})
}

func TestCompare(t *testing.T) {
	t.Run("returns 0 for equal XUIDs", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "test")
		id2, _ := xuid.NewWith(testUUID, "test")

		assert.Equal(t, 0, xuid.Compare(id1, id2))
	})

	t.Run("returns 0 only when prefix and UUID match", func(t *testing.T) {
		testUUID := uuid.New()
		id1, _ := xuid.NewWith(testUUID, "user")
		id2, _ := xuid.NewWith(testUUID, "order")

		assert.NotEqual(t, 0, xuid.Compare(id1, id2))
	})

	t.Run("sorts by prefix before UUID", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000001"), "a")
		id2, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000000"), "b")

		assert.Equal(t, -1, xuid.Compare(id1, id2))
		assert.Equal(t, 1, xuid.Compare(id2, id1))
	})

	t.Run("sorts empty prefix before non-empty prefix", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.New(), "")
		id2, _ := xuid.NewWith(uuid.New(), "user")

		assert.Equal(t, -1, xuid.Compare(id1, id2))
		assert.Equal(t, 1, xuid.Compare(id2, id1))
	})

	t.Run("sorts by UUID bytes for equal prefixes", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000001"), "test")
		id2, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000002"), "test")

		assert.Equal(t, -1, xuid.Compare(id1, id2))
		assert.Equal(t, 1, xuid.Compare(id2, id1))
		assert.Equal(t, 0, xuid.Compare(id1, id1))
	})

	t.Run("preserves chronological order for UUIDv7 with same prefix", func(t *testing.T) {
		id1 := xuid.MustNewSortable("user")
		time.Sleep(2 * time.Millisecond)
		id2 := xuid.MustNewSortable("user")

		assert.Equal(t, -1, xuid.Compare(id1, id2))
		assert.Equal(t, 1, xuid.Compare(id2, id1))
	})

	t.Run("sorts a slice of XUIDs", func(t *testing.T) {
		ids := []xuid.XUID{
			xuid.Must(xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000002"), "user")),
			xuid.Must(xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000001"), "user")),
			xuid.Must(xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000001"), "order")),
			xuid.Must(xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000000"), "")),
		}
		slices.SortFunc(ids, xuid.Compare)

		assert.Equal(t, "", ids[0].GetPrefix())
		assert.Equal(t, "order", ids[1].GetPrefix())
		assert.Equal(t, "00000000-0000-7000-8000-000000000001", ids[2].GetUUID().String())
		assert.Equal(t, "00000000-0000-7000-8000-000000000002", ids[3].GetUUID().String())
	})
}

func TestLess(t *testing.T) {
	t.Run("matches Compare", func(t *testing.T) {
		id1, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000001"), "test")
		id2, _ := xuid.NewWith(uuid.MustParse("00000000-0000-7000-8000-000000000002"), "test")

		assert.True(t, xuid.Less(id1, id2))
		assert.False(t, xuid.Less(id2, id1))
		assert.False(t, xuid.Less(id1, id1))
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

	t.Run("parses the unprefixed nil UUID string", func(t *testing.T) {
		parsed, err := xuid.Parse("1111111111111111")

		require.NoError(t, err)
		assert.True(t, xuid.IsEmpty(parsed))
		assert.Equal(t, "", parsed.GetPrefix())
	})

	t.Run("rejects a prefixed nil UUID string", func(t *testing.T) {
		_, err := xuid.Parse("user_1111111111111111")

		assert.ErrorIs(t, err, xuid.ErrParse)
		assert.ErrorIs(t, err, xuid.ErrNilUUIDWithPrefix)
		assert.False(t, xuid.IsValid("user_1111111111111111"))
	})

	t.Run("returns error for invalid XUID string", func(t *testing.T) {
		_, err := xuid.Parse("invalid_string")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xuid.ErrParse)
		// The wrapped cause explains why parsing failed.
		assert.ErrorContains(t, err, "want 16")
	})

	t.Run("returns error for malformed base58", func(t *testing.T) {
		_, err := xuid.Parse("test_invalid0characters")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xuid.ErrParse)
		// The wrapped cause identifies the invalid character.
		assert.ErrorContains(t, err, "invalid character")
	})

	t.Run("returns error for empty identifier", func(t *testing.T) {
		_, err := xuid.Parse("")
		assert.ErrorIs(t, err, xuid.ErrParse)

		_, err = xuid.Parse("user_")
		assert.ErrorIs(t, err, xuid.ErrParse)
	})

	t.Run("returns error for overlong identifier", func(t *testing.T) {
		_, err := xuid.Parse("user_" + strings.Repeat("2", 23))

		assert.ErrorIs(t, err, xuid.ErrParse)
		assert.ErrorContains(t, err, "exceeds 22")
	})

	t.Run("returns error for invalid prefix", func(t *testing.T) {
		id, _ := xuid.NewSortable("user")
		// Replace the valid prefix with one containing a disallowed
		// character, keeping the same (valid) identifier part.
		invalid := "bad-prefix_" + id.String()[len("user_"):]

		_, err := xuid.Parse(invalid)

		assert.ErrorIs(t, err, xuid.ErrParse)
		assert.ErrorIs(t, err, xuid.ErrInvalidPrefix)

		_, err = xuid.Parse(strings.Repeat("a", xuid.MaxPrefixLen+1) + "_" + id.String()[len("user_"):])
		assert.ErrorIs(t, err, xuid.ErrInvalidPrefix)
	})
}

func TestValidatePrefix(t *testing.T) {
	t.Run("accepts empty prefix", func(t *testing.T) {
		assert.NoError(t, xuid.ValidatePrefix(""))
	})

	t.Run("accepts letters, digits and underscores", func(t *testing.T) {
		assert.NoError(t, xuid.ValidatePrefix("user_order_42"))
		assert.NoError(t, xuid.ValidatePrefix("ABCxyz_0123"))
	})

	t.Run("accepts prefix at the maximum length", func(t *testing.T) {
		assert.NoError(t, xuid.ValidatePrefix(strings.Repeat("a", xuid.MaxPrefixLen)))
	})

	t.Run("rejects prefix over the maximum length", func(t *testing.T) {
		assert.ErrorIs(t, xuid.ValidatePrefix(strings.Repeat("a", xuid.MaxPrefixLen+1)), xuid.ErrInvalidPrefix)
	})

	t.Run("rejects hyphens and other characters", func(t *testing.T) {
		for _, prefix := range []string{"user-profile", "user profile", "user.profile", "user\n", "üser"} {
			assert.ErrorIs(t, xuid.ValidatePrefix(prefix), xuid.ErrInvalidPrefix, "prefix %q", prefix)
		}
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

	t.Run("returns false for invalid prefix", func(t *testing.T) {
		id, _ := xuid.NewSortable("user")
		invalid := "bad-prefix_" + id.String()[len("user_"):]

		assert.False(t, xuid.IsValid(invalid))
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
		// Using UUID v1 for testing (stdlib uuid has no v1 constructor,
		// so craft one by setting the version nibble directly)
		customUUID := uuid.MustParse("550e8400-e29b-11d4-a716-446655440000") // version 1
		id, _ := xuid.NewWith(customUUID, "test")

		assert.False(t, id.IsSortable())
		assert.False(t, id.IsRandom())
	})

	t.Run("rejects version 7 with a non-RFC variant", func(t *testing.T) {
		// The variant is the top two bits of octet 8; RFC 9562 is 10xx.
		// NCS (0xxx), Microsoft (110x) and future (111x) are not RFC.
		for _, variant := range []byte{0x00, 0x40, 0xc0, 0xe0} {
			raw := uuid.NewV7()
			raw[8] = variant

			id, err := xuid.NewWith(raw, "")
			require.NoError(t, err)
			assert.False(t, id.IsSortable(), "variant %#02x must not be sortable", variant)
		}
	})

	t.Run("rejects version 4 with a non-RFC variant", func(t *testing.T) {
		for _, variant := range []byte{0x00, 0x40, 0xc0, 0xe0} {
			raw := uuid.NewV4()
			raw[8] = variant

			id, err := xuid.NewWith(raw, "")
			require.NoError(t, err)
			assert.False(t, id.IsRandom(), "variant %#02x must not be random", variant)
		}
	})
}

func TestTime(t *testing.T) {
	t.Run("returns the embedded millisecond timestamp", func(t *testing.T) {
		want := time.Date(2024, 3, 14, 15, 9, 26, 0, time.UTC)
		id := mustXUIDv7At(t, want)

		got, err := id.Time()

		require.NoError(t, err)
		assert.Equal(t, want.UnixMilli(), got.UnixMilli())
	})

	t.Run("tracks creation time for sortable XUIDs", func(t *testing.T) {
		before := time.Now().Add(-time.Second)
		id := xuid.MustNewSortable("user")
		after := time.Now().Add(time.Second)

		got, err := id.Time()

		require.NoError(t, err)
		assert.False(t, got.Before(before))
		assert.False(t, got.After(after))
	})

	t.Run("ignores the prefix", func(t *testing.T) {
		want := time.Date(2024, 3, 14, 15, 9, 26, 0, time.UTC)
		id := mustXUIDv7At(t, want).MustWithPrefix("order")

		got, err := id.Time()

		require.NoError(t, err)
		assert.Equal(t, want.UnixMilli(), got.UnixMilli())
	})

	t.Run("returns ErrNotSortable for random UUID", func(t *testing.T) {
		id, _ := xuid.NewRandom("session")

		_, err := id.Time()

		assert.ErrorIs(t, err, xuid.ErrNotSortable)
	})

	t.Run("returns ErrNotSortable for version 7 with a non-RFC variant", func(t *testing.T) {
		raw := uuid.NewV7()
		raw[8] = 0x00 // NCS variant, not RFC 9562

		id, err := xuid.NewWith(raw, "")
		require.NoError(t, err)
		require.False(t, id.IsSortable())

		_, err = id.Time()

		assert.ErrorIs(t, err, xuid.ErrNotSortable)
	})

	t.Run("returns ErrNotSortable for nil UUID", func(t *testing.T) {
		id, _ := xuid.NilUUID()

		_, err := id.Time()

		assert.ErrorIs(t, err, xuid.ErrNotSortable)
	})
}

// mustXUIDv7At builds a UUIDv7 XUID whose 48-bit timestamp encodes at.
func mustXUIDv7At(t *testing.T, at time.Time) xuid.XUID {
	t.Helper()

	ms := at.UnixMilli()
	var raw [16]byte
	for i := 5; i >= 0; i-- {
		raw[i] = byte(ms)
		ms >>= 8
	}
	raw[6] = 0x70 // version 7
	raw[8] = 0x80 // RFC 9562 variant

	id, err := xuid.NewWith(uuid.UUID(raw), "user")
	require.NoError(t, err)
	return id
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

func TestWithPrefix(t *testing.T) {
	t.Run("WithPrefix sets prefix to XUID without prefix", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("")

		withPrefix, err := testXUID.WithPrefix("user")

		require.NoError(t, err)
		assert.Equal(t, "user", withPrefix.GetPrefix())
		assert.Equal(t, testXUID.GetUUID(), withPrefix.GetUUID())
	})

	t.Run("WithPrefix replaces existing prefix", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("old")

		withNewPrefix, err := testXUID.WithPrefix("new")

		require.NoError(t, err)
		assert.Equal(t, "new", withNewPrefix.GetPrefix())
	})

	t.Run("WithPrefix clears prefix when empty", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("old")

		withoutPrefix, err := testXUID.WithPrefix("")

		require.NoError(t, err)
		assert.Equal(t, "", withoutPrefix.GetPrefix())
		assert.Equal(t, testXUID.GetUUID(), withoutPrefix.GetUUID())
	})

	t.Run("WithPrefix leaves original unchanged", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("old")

		_, err := testXUID.WithPrefix("new")

		require.NoError(t, err)
		assert.Equal(t, "old", testXUID.GetPrefix())
	})

	t.Run("WithPrefix rejects invalid prefix", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("old")

		_, err := testXUID.WithPrefix("bad prefix")

		assert.ErrorIs(t, err, xuid.ErrInvalidPrefix)
		assert.Equal(t, "old", testXUID.GetPrefix())
	})

	t.Run("WithPrefix rejects a prefix on the nil UUID", func(t *testing.T) {
		nilXUID, err := xuid.NilUUID()
		require.NoError(t, err)

		_, err = nilXUID.WithPrefix("user")

		assert.ErrorIs(t, err, xuid.ErrNilUUIDWithPrefix)
	})

	t.Run("WithPrefix clears the prefix on the nil UUID", func(t *testing.T) {
		nilXUID, err := xuid.NilUUID()
		require.NoError(t, err)

		cleared, err := nilXUID.WithPrefix("")

		require.NoError(t, err)
		assert.True(t, xuid.IsEmpty(cleared))
		assert.Equal(t, "", cleared.GetPrefix())
	})

	t.Run("MustWithPrefix chains off non-addressable values", func(t *testing.T) {
		original := xuid.MustNewSortable("user")

		// MustParse returns a value (not a pointer), so the result is
		// non-addressable and chaining must work without storing it
		// in a variable first.
		restored := xuid.MustParse(original.String()).MustWithPrefix("user")

		assert.True(t, original.Equal(restored))
	})

	t.Run("MustWithPrefix panics on invalid prefix", func(t *testing.T) {
		testXUID, _ := xuid.NewSortable("old")

		assert.Panics(t, func() {
			testXUID.MustWithPrefix("bad prefix")
		})
	})

	t.Run("MustWithPrefix panics on a prefixed nil UUID", func(t *testing.T) {
		nilXUID, err := xuid.NilUUID()
		require.NoError(t, err)

		assert.Panics(t, func() {
			nilXUID.MustWithPrefix("user")
		})
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

	t.Run("rejects very long prefix", func(t *testing.T) {
		longPrefix := strings.Repeat("a", 100)
		_, err := xuid.NewSortable(longPrefix)

		assert.ErrorIs(t, err, xuid.ErrInvalidPrefix)
	})

	t.Run("accepts prefix at the maximum length", func(t *testing.T) {
		maxPrefix := strings.Repeat("a", xuid.MaxPrefixLen)
		id, err := xuid.NewSortable(maxPrefix)

		require.NoError(t, err)
		assert.Equal(t, maxPrefix, id.GetPrefix())
		assert.True(t, strings.HasPrefix(id.String(), maxPrefix+"_"))
	})

	t.Run("rejects prefix with special characters", func(t *testing.T) {
		specialPrefix := "test-prefix.with_special@chars"
		_, err := xuid.NewSortable(specialPrefix)

		assert.ErrorIs(t, err, xuid.ErrInvalidPrefix)
	})

	t.Run("accepts prefix with underscores", func(t *testing.T) {
		specialPrefix := "test_prefix_with_underscores"
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

func BenchmarkWithPrefix(b *testing.B) {
	id := xuid.MustNewSortable("bench")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = id.WithPrefix("bench")
	}
}

// FuzzParseRoundTrip asserts that every non-nil XUID survives a String/Parse
// round trip: Parse(x.String()) must equal x. It fuzzes both the prefix
// (including underscores and non-ASCII runes) and the raw UUID bytes, so
// malformed prefixes or encodings cannot silently corrupt an identifier.
// Inputs with a prefix that fails validation are skipped, since they
// cannot be constructed in the first place. The nil UUID is checked
// separately: its canonical form is the empty string, which Parse rejects.
func FuzzParseRoundTrip(f *testing.F) {
	// Unprefixed nil UUID.
	f.Add("", []byte{})
	// Prefixed non-nil UUID (a prefixed nil UUID is rejected on purpose).
	f.Add("user", []byte{0, 0, 0, 0, 0, 0, 0x70, 0, 0x80, 0, 0, 0, 0, 0, 0, 0x01})
	// Unprefixed max UUID (all 0xff).
	f.Add("", slices.Repeat([]byte{0xff}, 16))
	// Prefixed v4 (random) sample.
	v4 := uuid.NewV4()
	f.Add("user", v4[:])
	// Prefixed v7 (sortable) sample.
	v7 := uuid.NewV7()
	f.Add("order", v7[:])
	// Prefix with underscores, non-ASCII, and short raw payloads.
	f.Add("user_profile", []byte{0xff})
	f.Add("_", []byte{0x00})
	f.Add("user", []byte{
		0x01, 0x8f, 0x2a, 0x3b, 0x4c, 0x5d, 0x6e, 0x7f,
		0x80, 0x91, 0xa2, 0xb3, 0xc4, 0xd5, 0xe6, 0xf7,
	})

	f.Fuzz(func(t *testing.T, prefix string, raw []byte) {
		var u uuid.UUID
		copy(u[:], raw)

		original, err := xuid.NewWith(u, prefix)
		if err != nil {
			// Invalid prefixes are rejected by the constructor; there
			// is no valid identifier to round-trip.
			return
		}

		if xuid.IsEmpty(original) {
			// The nil UUID renders as the empty string, which Parse
			// rejects by design, so it has no String/Parse round trip.
			assert.Equal(t, "", original.String())
			return
		}

		parsed, err := xuid.Parse(original.String())
		require.NoError(t, err)
		assert.True(t, original.Equal(parsed),
			"round trip mismatch: %q != %q", original.String(), parsed.String())
	})
}

// FuzzParseNoPanic asserts the failure contract of Parse on arbitrary
// input: it never panics, every rejected string produces an error that
// wraps ErrParse, and anything it accepts is either the nil UUID (whose
// canonical form is the empty string) or survives a String/Parse round
// trip. Inputs are decoded byte-for-byte, so non-ASCII and overlong
// strings are exercised alongside well-formed identifiers.
func FuzzParseNoPanic(f *testing.F) {
	seeds := []string{
		"",
		"_",
		"user_",
		"user_2g",
		"user_8M7Qq2vR3kGbF9wN5pL2xA",
		"1",
		strings.Repeat("1", 22),
		strings.Repeat("z", 22),
		strings.Repeat("z", 23),
		"0OIl",
		"user_2g_extra",
		"üser_2g",
		strings.Repeat("a", xuid.MaxPrefixLen+1) + "_2g",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got, err := xuid.Parse(s)
		if err != nil {
			require.ErrorIs(t, err, xuid.ErrParse)
			assert.False(t, xuid.IsValid(s))
			return
		}

		assert.True(t, xuid.IsValid(s))
		if xuid.IsEmpty(got) {
			// The nil UUID renders as the empty string, which Parse
			// rejects by design, so it has no String/Parse round trip.
			assert.Equal(t, "", got.String())
			return
		}
		// Anything else Parse accepts must round-trip through String.
		assert.True(t, xuid.MustParse(got.String()).Equal(got),
			"round trip mismatch for %q", s)
	})
}
