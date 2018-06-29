package main

import (
	"Go-exercise/CYOA/cyoa"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var tp1 *template.Template

func init() {
	tp1 = template.Must(template.ParseFiles("index.html"))
}

func main() {
	fileName := flag.String("file", "story.json", "Filename of the json story")
	flag.Parse()
	fmt.Printf("story was parsed from the %s file", *fileName)
	file, err := os.Open(*fileName)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing %s file", *fileName))
	}

	decoder := json.NewDecoder(file)
	var stories cyoa.Stories
	err = decoder.Decode(&stories)
	if err != nil {
		panic("Error while decoding the json.")
	}

	handler := cyoa.StoryHandler(stories, tp1)
	http.ListenAndServe(":8087", handler)
	fmt.Println("Server is listening at the port 8087")
}
