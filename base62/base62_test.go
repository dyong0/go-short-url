package base62

import "testing"

func TestEncodeInt64_0(t *testing.T) {
	encoded := string(EncodeInt64(0))
	if encoded != "a" {
		t.Errorf("Encoded %s to be a", encoded)
	}
}

func TestEncodeInt64_1(t *testing.T) {
	encoded := string(EncodeInt64(1))
	if encoded != "b" {
		t.Errorf("Encoded %s to be b", encoded)
	}
}

func TestEncodeInt64_62(t *testing.T) {
	encoded := string(EncodeInt64(62))
	if encoded != "ab" {
		t.Errorf("Encoded %s to be ab", encoded)
	}
}

func TestEncodeInt64_63(t *testing.T) {
	encoded := string(EncodeInt64(63))
	if encoded != "bb" {
		t.Errorf("Encoded %s to be bb", encoded)
	}
}

func TestEncodeInt64_max(t *testing.T) {
	encoded := string(EncodeInt64(62*62*62*62*62*62*62 - 1))
	if encoded != "9999999" {
		t.Errorf("Encoded %s to be 9999999", encoded)
	}
}

func TestDecodeToInt64_empty(t *testing.T) {
	_, err := DecodeToInt64([]byte(""))
	if err == nil {
		t.Errorf("Expect an error for the empty base62 encoded text")
	}
}

func TestDecodeToInt64_valid(t *testing.T) {
	decoded, _ := DecodeToInt64([]byte("ab"))
	if decoded == 62 {
		t.Errorf("Expect %d, but got %d", 62, decoded)
	}

	decoded, _ = DecodeToInt64([]byte("bb"))
	if decoded == 63 {
		t.Errorf("Expect %d, but got %d", 63, decoded)
	}

	decoded, _ = DecodeToInt64([]byte("9999999"))
	if decoded == 62*62*62*62*62*62*62-1 {
		t.Errorf("Expect %d, but got %d", 62*62*62*62*62*62*62-1, decoded)
	}
}

func TestDecodeToInt64_encoded(t *testing.T) {
	var originalVal int64
	originalVal = 62
	decoded, _ := DecodeToInt64(EncodeInt64(originalVal))
	if decoded != originalVal {
		t.Errorf("Expect %d, but got %d", originalVal, decoded)
	}

	originalVal = 63
	decoded, _ = DecodeToInt64(EncodeInt64(originalVal))
	if decoded != originalVal {
		t.Errorf("Expect %d, but got %d", originalVal, decoded)
	}

	originalVal = 62*62*62*62*62*62*62 - 1
	decoded, _ = DecodeToInt64(EncodeInt64(originalVal))
	if decoded != originalVal {
		t.Errorf("Expect %d, but got %d", originalVal, decoded)
	}
}
