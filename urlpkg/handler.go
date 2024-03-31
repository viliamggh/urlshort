package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// PathURLMapping represents a single path to URL mapping
type PathURLMapping struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
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

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil

}

func parseYAML(yml []byte) ([]PathURLMapping, error) {
	var mappings []PathURLMapping
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}

	return mappings, nil
}

func buildMap(mappings []PathURLMapping) map[string]string {
	pathMap := make(map[string]string)
	for _, mapping := range mappings {
		pathMap[mapping.Path] = mapping.URL
	}
	return pathMap
}
