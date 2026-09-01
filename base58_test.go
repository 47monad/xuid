package xuid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBase58Alphabet(t *testing.T) {
	t.Run("uses the Bitcoin base58 alphabet", func(t *testing.T) {
		expected := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

		assert.Equal(t, expected, alphabet)
	})

	t.Run("alphabet has 58 unique characters", func(t *testing.T) {
		assert.Len(t, alphabet, 58)

		seen := make(map[byte]bool)
		for i := 0; i < len(alphabet); i++ {
			assert.False(t, seen[alphabet[i]], "duplicate character %q in alphabet", string(alphabet[i]))
			seen[alphabet[i]] = true
		}
	})

	t.Run("alphabet excludes visually ambiguous characters", func(t *testing.T) {
		for _, c := range "0OIl" {
			assert.NotContains(t, alphabet, string(c), "alphabet must not contain ambiguous character %q", string(c))
		}
	})
}

func TestEncodeBase58(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{"empty", []byte{}, ""},
		{"single zero byte", []byte{0x00}, "1"},
		{"single byte 0x61", []byte{0x61}, "2g"},
		{"three bytes", []byte{0x62, 0x62, 0x62}, "a3gV"},
		{"hello world", []byte("hello world"), "StV1DL6CwTryKyV"},
		{"leading zeros map to '1' padding", []byte{0x00, 0x00, 0x00, 0x01}, "1112"},
		{"all zero bytes", make([]byte, 16), strings.Repeat("1", 16)},
		{
			"all 0xff bytes encode to max length",
			[]byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
			"YcVfxkQb6JRzqk5kF2tNLv", // (2^128 - 1) in base58, 22 digits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, encodeBase58(tt.in))
		})
	}

	t.Run("16-byte output never exceeds maxEncodedLen", func(t *testing.T) {
		// The worst case is all-0xff bytes; sweep a range of extreme inputs.
		inputs := [][]byte{
			make([]byte, 16),
			[]byte(strings.Repeat("\xff", 16)),
			{0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00},
		}

		for _, in := range inputs {
			assert.LessOrEqual(t, len(encodeBase58(in)), maxEncodedLen)
		}
	})
}

func TestDecodeBase58(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []byte
		wantErr bool
	}{
		{"empty", "", []byte{}, false},
		{"single '1'", "1", []byte{0x00}, false},
		{"single byte", "2g", []byte{0x61}, false},
		{"three bytes", "a3gV", []byte{0x62, 0x62, 0x62}, false},
		{"hello world", "StV1DL6CwTryKyV", []byte("hello world"), false},
		{"leading '1's decode to zero bytes", "1112", []byte{0x00, 0x00, 0x00, 0x01}, false},
		{"embedded '1' is a zero digit", "12", []byte{0x00, 0x01}, false},
		{"zero digit", "0", nil, true},
		{"capital O", "O", nil, true},
		{"capital I", "I", nil, true},
		{"lowercase l", "l", nil, true},
		{"malformed with zero digit", "StV1DL6CwTryKy0", nil, true},
		{"non-ascii rune", "StV1DL6CwTryKä", nil, true},
		{"exceeds max length", strings.Repeat("2", maxEncodedLen+1), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeBase58(tt.in)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBase58RoundTrip(t *testing.T) {
	tests := [][]byte{
		{},
		{0x00},
		{0x01},
		{0x00, 0x00, 0x00, 0x01},
		[]byte("hello world"),
		make([]byte, 16),
		{0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00},
		[]byte(strings.Repeat("\xff", 16)),
		{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
	}

	for _, in := range tests {
		encoded := encodeBase58(in)

		assert.LessOrEqual(t, len(encoded), maxEncodedLen, "encoding % x exceeds maxEncodedLen", in)

		got, err := decodeBase58(encoded)
		require.NoError(t, err, "decodeBase58(%q)", encoded)
		assert.Equal(t, in, got, "round trip of % x", in)
	}
}

func TestParseRejectsOverlongIdentifier(t *testing.T) {
	t.Run("identifier longer than maxEncodedLen", func(t *testing.T) {
		long := strings.Repeat("2", maxEncodedLen+1)

		_, err := Parse("test_" + long)

		assert.ErrorIs(t, err, ErrParse)
	})

	t.Run("empty identifier", func(t *testing.T) {
		_, err := Parse("test_")

		assert.ErrorIs(t, err, ErrParse)

		_, err = Parse("")

		assert.ErrorIs(t, err, ErrParse)
	})

	t.Run("identifier at max length is not rejected by the guard", func(t *testing.T) {
		// A 22-char identifier is the canonical maximum and must reach
		// the decoder rather than fail the length check.
		_, err := decodeBase58(strings.Repeat("z", maxEncodedLen))

		assert.NotErrorIs(t, err, errBase58TooLong)
	})

	t.Run("22 leading '1's decode to 22 zero bytes and fail Parse", func(t *testing.T) {
		// "111...1" (22 chars) passes the length guard but decodes to 22
		// zero bytes, which can never be a 16-byte UUID.
		decoded, err := decodeBase58(strings.Repeat("1", maxEncodedLen))

		require.NoError(t, err)
		assert.Equal(t, make([]byte, maxEncodedLen), decoded)

		_, err = Parse("test_" + strings.Repeat("1", maxEncodedLen))

		assert.ErrorIs(t, err, ErrParse)
	})
}

// FuzzBase58 asserts the core invariant: any string that decodes
// successfully must re-encode to a string that decodes to the same bytes.
func FuzzBase58(f *testing.F) {
	seeds := []string{
		"",
		"1",
		"2g",
		"StV1DL6CwTryKyV",
		"1112",
		"0",
		"l",
		strings.Repeat("z", maxEncodedLen),
		strings.Repeat("2", maxEncodedLen+1),
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		decoded, err := decodeBase58(s)
		if err != nil {
			return
		}

		round, err := decodeBase58(encodeBase58(decoded))
		if err != nil {
			t.Fatalf("decodeBase58(encodeBase58(decodeBase58(%q))) error = %v", s, err)
		}
		if len(round) != len(decoded) {
			t.Fatalf("round trip mismatch for %q: % x vs % x", s, round, decoded)
		}
		for i := range round {
			if round[i] != decoded[i] {
				t.Fatalf("round trip mismatch for %q: % x vs % x", s, round, decoded)
			}
		}
	})
}
