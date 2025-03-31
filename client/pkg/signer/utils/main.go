package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// LowerCase implement AWS definition:
// Convert the string to lowercase.
func LowerCase(v string) string {
	return strings.ToLower(v)
}

// Hex implementation following AWS specification:
// Lowercase base 16 encoding.
func Hex(v []byte) string {
	return hex.EncodeToString(v)
}

// SHA256Hash implementation following AWS specification:
// Secure Hash Algorithm (SHA) cryptographic hash function.
func SHA256Hash(v []byte) []byte {
	sum := sha256.Sum256(v)
	return sum[:]
}

// HMAC_SHA256 implementation followinf AWS specification:
// Computes HMAC by using the SHA256 algorithm with the signing key provided. This is the final signature.
func HMAC_SHA256(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

// Tim implementation followinf AWS specification:
// Remove any leading or trailing whitespace.
func Trim(v string) string {
	return strings.TrimSpace(v)
}

// URIEncode URI encode every byte. UriEncode() must enforce the following rules:
// URI encode every byte except the unreserved characters: 'A'-'Z', 'a'-'z', '0'-'9', '-', '.', '_', and '~'.
// The space character is a reserved character and must be encoded as "%20" (and not as "+").
// Each URI encoded byte is formed by a '%' and the two-digit hexadecimal value of the byte.
// Letters in the hexadecimal value must be uppercase, for example "%1A".
// Encode the forward slash character, '/', everywhere except in the object key name. For example, if the object key name is photos/Jan/sample.jpg, the forward slash in the key name is not encoded.
func URIEncode(v string, isObjectKey bool) string {
	var ret strings.Builder

	for i := range len(v) {
		c := v[i]
		if !shouldEscape(c, isObjectKey) {
			ret.WriteByte(c)
			continue
		}

		ret.WriteString(fmt.Sprintf("%%%02X", c))
	}

	return ret.String()
}

func shouldEscape(v byte, isObjectKey bool) bool {
	if v >= 'A' && v <= 'Z' ||
		v >= 'a' && v <= 'z' ||
		v >= '0' && v <= '9' {
		return false
	}

	switch v {
	case '-', '.', '_', '~':
		return false

	case '/':
		return !isObjectKey

	default:
		return true
	}
}
