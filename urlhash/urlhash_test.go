package urlhash

import "testing"

func TestStoreURL_sameURLs(t *testing.T) {
	urlHash := NewURLHash()
	first, _ := urlHash.StoreURL("https://google.com/foo/bar")
	second, _ := urlHash.StoreURL("https://google.com/foo/bar")
	if first != second {
		t.Errorf("The same URLs must have the same hash values")
	}
}

func TestStoreURL_differentURLs(t *testing.T) {
	urlHash := NewURLHash()
	first, _ := urlHash.StoreURL("https://google.com/foo/bar")
	second, _ := urlHash.StoreURL("https://google.com/basmgioasmoi/mlksmrisam")
	if first == second {
		t.Errorf("The same URLs must have the same hash values")
	}
}

func TestFindURLByHash_existingHash(t *testing.T) {
	urlHash := NewURLHash()
	hash, _ := urlHash.StoreURL("https://google.com/foo/bar")
	urlFound, err := urlHash.FindURLByHash(hash)

	if err != nil {
		t.Errorf("Expect no error but got %v", err)
	}

	if urlFound != "https://google.com/foo/bar" {
		t.Errorf("Expect the hash of %s, but got %s", "https://google.com/foo/bar", urlFound)
	}
}

func TestStoreURL_sameHash(t *testing.T) {
	urlHash := NewURLHash()
	first, _ := urlHash.StoreURL("https://google.com/asdfasefafsaf/foo/bar")
	second, _ := urlHash.StoreURL("https://google.com/basmgioasmoi/foo/bar")
	if first == second {
		t.Errorf("The same URLs must have the same hash values")
	}
}
