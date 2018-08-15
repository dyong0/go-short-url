package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/dyong0/go-short-url/urlhash"
)

func main() {
	urlHash := urlhash.NewURLHash()

	router := &http.ServeMux{}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetURLByHash(urlHash, w, r)
			break
		case http.MethodPost:
			handleCreateHashOfURL(urlHash, w, r)
			break
		}
	})

	log.Fatal(http.ListenAndServe(":3000", router))
}

func handleCreateHashOfURL(urlHash *urlhash.URLHash, w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	escapedURL := url.PathEscape(originalURL)

	log.Printf("Storing URL %s\n", escapedURL)
	hash, err := urlHash.StoreURL(escapedURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hash))
}

func handleGetURLByHash(urlHash *urlhash.URLHash, w http.ResponseWriter, r *http.Request) {
	hash := r.URL.EscapedPath()[1:]
	log.Printf("Retrieving the url of hash %s\n", hash)
	urlFound, err := urlHash.FindURLByHash(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	urlUnescaped, err := url.PathUnescape(urlFound)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", urlUnescaped)
	w.WriteHeader(http.StatusFound)
}
