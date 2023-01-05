package urlshort

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	bolt "go.etcd.io/bbolt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if path, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}

	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type T struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func YAMLtoMap(yml []byte) map[string]string {

	t := []T{}

	err := yaml.Unmarshal(yml, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	result := make(map[string]string)
	for _, entry := range t {
		result[entry.Path] = entry.URL
	}

	return result
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathsToUrls := YAMLtoMap(yml)

	return func(w http.ResponseWriter, r *http.Request) {

		if path, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}

	}, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	{
//			"mapping": [
//				{
//					"path": "..."
//					"url": "..."
//				},
//				...
//			]
//	}
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type URLs struct {
	URLs []URL `json:"mapping"`
}

type URL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func JSONtoMap(jsn []byte) map[string]string {

	u := URLs{}

	err := json.Unmarshal(jsn, &u)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	result := make(map[string]string)
	for _, entry := range u.URLs {
		result[entry.Path] = entry.URL
	}

	return result
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathsToUrls := JSONtoMap(jsn)

	return func(w http.ResponseWriter, r *http.Request) {

		if path, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}

	}, nil
}

// BoltDBHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the BoltDB, then the
// fallback http.Handler will be called instead.
//
// BoltDB entries are encoded and must be written to the
// database using golang.
//
// The only errors that can be returned all related to having
// invalid BoltDB entries.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func BoltDBtoMap(db *bolt.DB) map[string]string {

	result := make(map[string]string)

	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("mapping"))

		fmt.Println("BoltDB entries:")
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("%s: %s,\n", string(k), string(v))
			result[string(k)] = string(v)
			return nil
		})

		return errors.New("database failed during reading")
	})

	return result
}

func BoltHandler(blt *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {

	pathsToUrls := BoltDBtoMap(blt)

	return func(w http.ResponseWriter, r *http.Request) {

		if path, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}

	}, nil
}
