package main

import (
	"Go-exercise/URL-Short/handler"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	ymlfile := flag.String("ymlfile", "sample.yaml", "Yaml file name")
	jsonfile := flag.String("jsonfile", "sample.json", "json file name")
	flag.Parse()
	mux := defaultMux()
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)
	yaml, err := ioutil.ReadFile(*ymlfile)
	jsondata, err := ioutil.ReadFile(*jsonfile)

	if err != nil {
		panic(err)
	}

	mapHandler = handler.YAMLHandler([]byte(yaml), mapHandler)
	mapHandler = handler.JSONHandler([]byte(jsondata), mapHandler)
	http.ListenAndServe(":8087", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to URL-Short")
	})
	return mux
}
