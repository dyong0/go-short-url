package base62

import (
	"fmt"
)

const (
	base62Codes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base62Base  = int64(len(base62Codes))
)

var base62Values = map[byte]int64{
	'a': 0,
	'b': 1,
	'c': 2,
	'd': 3,
	'e': 4,
	'f': 5,
	'g': 6,
	'h': 7,
	'i': 8,
	'j': 9,
	'k': 10,
	'l': 11,
	'm': 12,
	'n': 13,
	'o': 14,
	'p': 15,
	'q': 16,
	'r': 17,
	's': 18,
	't': 19,
	'u': 20,
	'v': 21,
	'w': 22,
	'x': 23,
	'y': 24,
	'z': 25,
	'A': 26,
	'B': 27,
	'C': 28,
	'D': 29,
	'E': 30,
	'F': 31,
	'G': 32,
	'H': 33,
	'I': 34,
	'J': 35,
	'K': 36,
	'L': 37,
	'M': 38,
	'N': 39,
	'O': 40,
	'P': 41,
	'Q': 42,
	'R': 43,
	'S': 44,
	'T': 45,
	'U': 46,
	'V': 47,
	'W': 48,
	'X': 49,
	'Y': 50,
	'Z': 51,
	'0': 52,
	'1': 53,
	'2': 54,
	'3': 55,
	'4': 56,
	'5': 57,
	'6': 58,
	'7': 59,
	'8': 60,
	'9': 61,
}

// EncodeInt64 encodes an integer into a 62 based byte array
func EncodeInt64(plainVal int64) []byte {
	if plainVal == 0 {
		return []byte("a")
	}

	encoded := []byte{}

	for ; plainVal > 0; plainVal = plainVal / base62Base {
		encoded = append(encoded, base62Codes[plainVal%base62Base])
	}

	return encoded
}

// DecodeToInt64 decodes an base62-encoded bytes to int64
func DecodeToInt64(encoded []byte) (int64, error) {
	if len(encoded) == 0 {
		return 0, fmt.Errorf("Empty base62-encoded text has no value")
	}

	var multiplier int64 = 1
	var sum int64
	for i := 0; i < len(encoded); i++ {
		codeValue, ok := base62Values[encoded[i]]
		if !ok {
			return 0, fmt.Errorf("Not base62-encoded")
		}
		sum = sum + (codeValue * multiplier)
		multiplier = multiplier * 62
	}

	return sum, nil
}
