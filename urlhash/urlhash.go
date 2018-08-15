package urlhash

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/dyong0/go-short-url/base62"
)

const (
	hashValueLength   = 7
	totalOuterBuckets = 62 * 62 * 62      // 3 characters of base62 < 2^18
	totalInnerBuckets = 62 * 62 * 62 * 62 // 4 characters of base62 < 2^24
	maxHashValue      = 62 * 62 * 62 * 62 * 62 * 62 * 62
)

// URLHash provides url-to-hash conversion and vice versa.
// Note that it has the length of 7 fixed so far.
type URLHash struct {
	outerBuckets [totalOuterBuckets][]string
}

// NewURLHash creates a new URLHash.
func NewURLHash() *URLHash {
	return &URLHash{}
}

// StoreURL stores a URL and returns its corresponding hash value with the length of 7.
// It also returns an error with an empty string in case there's no available space to store the URL.
func (h *URLHash) StoreURL(url string) (string, error) {
	hashValue := base62EncodedHashValueFromURL(url)

	outerBucketIndex := hashValue / totalInnerBuckets
	for {
		innerBuckets := &h.outerBuckets[outerBucketIndex]
		if *innerBuckets == nil {
			*innerBuckets = make([]string, totalInnerBuckets)
		}

		storedHashValue, err := storeURLInInnerBuckets(url, hashValue, innerBuckets)
		if err == nil {
			resultBuf := bytes.NewBuffer(make([]byte, 8))
			binary.Write(resultBuf, binary.LittleEndian, storedHashValue)
			return string(base62.EncodeInt64(storedHashValue)), nil
		}

		outerBucketIndex++

		if outerBucketIndex == hashValue/totalOuterBuckets {
			return "", fmt.Errorf("No empty bucket left")
		}
	}
}

// FindURLByHash finds URL by a hash
func (h *URLHash) FindURLByHash(hash string) (string, error) {
	hashValue, err := base62.DecodeToInt64([]byte(hash))
	if err != nil {
		return "", err
	}

	outerBucketIndex := hashValue / totalInnerBuckets
	for {
		innerBuckets := &h.outerBuckets[outerBucketIndex]
		if *innerBuckets == nil {
			return "", fmt.Errorf("No matching URL found")
		}

		url, err := findURLInInnerBuckets(hashValue, innerBuckets)
		if err != nil {
			return "", err
		}

		if url != "" {
			return url, nil
		}

		outerBucketIndex++
	}
}

func base62EncodedHashValueFromURL(url string) int64 {
	// URL to int64
	var urlInInt64 int64
	buf := bytes.NewBuffer([]byte(url))
	binary.Read(buf, binary.LittleEndian, &urlInInt64)

	// int64 to base62
	base62Encoded := base62.EncodeInt64(urlInInt64)

	// base62 to int64
	var hashValue int64
	buf = bytes.NewBuffer(base62Encoded)
	binary.Read(buf, binary.LittleEndian, &hashValue)

	return hashValue % maxHashValue
}

func storeURLInInnerBuckets(url string, hashValue int64, ptrInnerBuckets *[]string) (int64, error) {
	innerBuckets := *ptrInnerBuckets
	innerBucketIndex := hashValue % totalInnerBuckets
	firstTryAt := innerBucketIndex

	var i int64
	for ; ; i++ {
		// Empty bucket found
		if innerBuckets[innerBucketIndex] == "" {
			innerBuckets[innerBucketIndex] = url

			return hashValue + i, nil
		}

		// When the URL already exists
		if innerBuckets[innerBucketIndex] == url {
			return hashValue + i, nil
		}

		// Check the next bucket
		innerBucketIndex = (innerBucketIndex + 1) % totalInnerBuckets
		if innerBucketIndex == firstTryAt {
			return 0, fmt.Errorf("No available bucket")
		}
	}
}

func findURLInInnerBuckets(hashValue int64, ptrInnerBuckets *[]string) (string, error) {
	innerBuckets := *ptrInnerBuckets
	innerBucketIndex := hashValue % totalInnerBuckets
	firstTryAt := innerBucketIndex

	var i int64
	for ; ; i++ {
		// Empty bucket found
		if innerBuckets[innerBucketIndex] == "" {
			return "", nil
		}

		// A matching URL found
		if innerBuckets[innerBucketIndex] != "" {
			return innerBuckets[innerBucketIndex], nil
		}

		// See all buckets until seeing the first bucket again
		innerBucketIndex = (innerBucketIndex + 1) % totalInnerBuckets
		if innerBucketIndex == firstTryAt {
			return "", fmt.Errorf("All buckets are empty")
		}
	}
}
