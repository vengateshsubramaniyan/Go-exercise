package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-yaml/yaml"
)

type urls struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

//MapHandler to server for the values stored in the map
func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if _, ok := paths[path]; ok {
			http.Redirect(w, r, paths[path], http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

//YAMLHandler to handle yaml files
func YAMLHandler(paths []byte, fallback http.Handler) http.HandlerFunc {

	var pathUrls []urls
	var err error
	err = yaml.Unmarshal(paths, &pathUrls)
	if err != nil {
		panic(err)
	}

	mapPaths := make(map[string]string)
	for _, val := range pathUrls {
		mapPaths[val.Path] = val.URL
	}
	return MapHandler(mapPaths, fallback)
}

//JSONHandler json file
func JSONHandler(paths []byte, fallback http.Handler) http.HandlerFunc {

	var pathUrls []urls
	var err error
	err = json.Unmarshal(paths, &pathUrls)
	if err != nil {
		panic(err)
	}

	mapPaths := make(map[string]string)
	for _, val := range pathUrls {
		fmt.Println("json")
		mapPaths[val.Path] = val.URL
	}
	return MapHandler(mapPaths, fallback)
}
