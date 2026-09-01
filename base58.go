package xuid

import (
	"errors"
	"strings"
)

// alphabet is the Bitcoin base58 alphabet. It excludes the visually
// ambiguous characters 0 (zero), O (capital o), I (capital i), and
// l (lowercase L) for readability.
const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

const (
	// alphabetZero is '1', which represents a leading zero byte.
	alphabetZero = '1'

	// maxEncodedLen is the maximum number of base58 digits needed to
	// represent a 16-byte UUID: ceil(16 * log(256) / log(58)) = 22.
	maxEncodedLen = 22

	invalid = 0xFF
)

// Internal sentinel errors. They are deliberately unexported: Parse maps
// every decode failure to the public ErrParse.
var (
	errBase58TooLong          = errors.New("base58 string exceeds maximum length")
	errBase58InvalidCharacter = errors.New("base58 string contains invalid character")
)

// decodeMap maps each ASCII byte to its base58 digit value, or invalid
// (0xFF) for bytes outside the alphabet.
var decodeMap = func() (m [256]byte) {
	for i := range m {
		m[i] = invalid
	}
	for i := 0; i < len(alphabet); i++ {
		m[alphabet[i]] = byte(i)
	}
	return m
}()

// encodeBase58 encodes b using the Bitcoin base58 alphabet. Each leading
// zero byte of b is represented by a leading '1'.
func encodeBase58(b []byte) string {
	zeros := 0
	for zeros < len(b) && b[zeros] == 0 {
		zeros++
	}

	// Each remaining byte contributes up to log(256)/log(58) ≈ 1.365
	// digits; 138/100 is a safe upper bound. Digits are little-endian.
	digits := make([]byte, 0, (len(b)-zeros)*138/100)
	for _, c := range b[zeros:] {
		carry := int(c)
		for i := range digits {
			carry += 256 * int(digits[i])
			digits[i] = byte(carry % 58)
			carry /= 58
		}
		for carry > 0 {
			digits = append(digits, byte(carry%58))
			carry /= 58
		}
	}

	var sb strings.Builder
	sb.Grow(zeros + len(digits))
	for i := 0; i < zeros; i++ {
		sb.WriteByte(alphabetZero)
	}
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(alphabet[digits[i]])
	}
	return sb.String()
}

// decodeBase58 decodes a Bitcoin base58 string. Leading '1's decode to
// leading zero bytes. It returns an error if s contains a character
// outside the alphabet or exceeds maxEncodedLen.
func decodeBase58(s string) ([]byte, error) {
	if len(s) > maxEncodedLen {
		return nil, errBase58TooLong
	}

	zeros := 0
	for zeros < len(s) && s[zeros] == alphabetZero {
		zeros++
	}

	// Each base58 digit contributes up to log(58)/log(256) ≈ 0.733
	// bytes. Bytes are little-endian.
	b := make([]byte, 0, (len(s)-zeros)*733/1000+1)
	for i := zeros; i < len(s); i++ {
		v := decodeMap[s[i]]
		if v == invalid {
			return nil, errBase58InvalidCharacter
		}
		carry := int(v)
		for j := range b {
			carry += 58 * int(b[j])
			b[j] = byte(carry % 256)
			carry /= 256
		}
		for carry > 0 {
			b = append(b, byte(carry%256))
			carry /= 256
		}
	}

	out := make([]byte, zeros+len(b))
	for i, v := range b {
		out[len(out)-1-i] = v
	}
	return out, nil
}
