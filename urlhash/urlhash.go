package urlhash

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dyong0/go-short-url/encode"
)

const (
	hashLength  = 7
	hashmapSize = 62 * 62 * 62 * 62 * 62 * 62 * 62 // 7 characters of base62
)

type hashmap [hashmapSize]string

var buckets hashmap

// StoreURL stores a URL and returns its corresponding hash with the length of 7.
// It also returns an error with an empty string in case there's no available space to store the URL.
func StoreURL(url string) (string, error) {
	var strippedIntVal int64
	buf := bytes.NewBuffer([]byte(url)[len(url)-hashLength:]) // last hashLength bytes
	binary.Read(buf, binary.LittleEndian, &strippedIntVal)

	encoded := encode.Int64ToBase62(strippedIntVal)

	var valEncoded int
	buf = bytes.NewBuffer(encoded)
	binary.Read(buf, binary.LittleEndian, &valEncoded)

	bucketIndex := valEncoded % hashmapSize
	firstTryAt := bucketIndex

	for {
		if buckets[bucketIndex] == "" {
			buckets[bucketIndex] = url

			binary.Write(buf, binary.LittleEndian, &valEncoded)
			return string(buf.Bytes()), nil
		}

		bucketIndex++

		if bucketIndex == hashmapSize {
			bucketIndex = 0
		}

		if bucketIndex == firstTryAt {
			return "", fmt.Errorf("All buckets are full: %d out of %d", hashmapSize, hashmapSize)
		}
	}
}
