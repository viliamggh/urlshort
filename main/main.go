package main

import (
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB
var err error

func main() {

	db, err = bolt.Open("paths.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// // Create bucket if doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("URLMappings"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	// Example: Insert a mapping
	err = insertMapping("/carss", "https://www.sauto.cz")
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	handler := DbHandler(db, mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func insertMapping(path, url string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("URLMappings"))
		return bucket.Put([]byte(path), []byte(url))
	})
}

func getMapping(path string) (string, error) {
	var url string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("URLMappings"))
		urlBytes := bucket.Get([]byte(path))
		url = string(urlBytes)
		return nil
	})
	return url, err
}

func DbHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("Request path: %s", path) // Log the request path

		dest, err := getMapping(path)
		if err == nil && dest != "" {
			log.Printf("Redirecting to: %s", dest) // Log the redirection
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		log.Printf("No mapping found for path: %s, falling back to default handler.", path) // Log the fallback

		fallback.ServeHTTP(w, r)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
