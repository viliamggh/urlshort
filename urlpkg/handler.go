package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// PathURLMapping represents a single path to URL mapping
type PathURLMapping struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type DataParser interface {
	Parse(data []byte) ([]PathURLMapping, error)
}

type YamlParser struct{}

func (yp YamlParser) Parse(data []byte) ([]PathURLMapping, error) {
	var mappings []PathURLMapping
	err := yaml.Unmarshal(data, &mappings)
	return mappings, err
}

type JsonParser struct{}

func (jp JsonParser) Parse(data []byte) ([]PathURLMapping, error) {
	var mappings []PathURLMapping
	err := json.Unmarshal(data, &mappings)
	return mappings, err
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func UniversalHandler(parser DataParser, data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedData, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedData)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(mappings []PathURLMapping) map[string]string {
	pathMap := make(map[string]string)
	for _, mapping := range mappings {
		pathMap[mapping.Path] = mapping.URL
	}
	return pathMap
}
