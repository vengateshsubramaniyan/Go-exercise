package main

import (
	"Go-exercise/Link-Parser/link"
	"flag"
	"fmt"
	"os"
)

func main() {
	fileName := flag.String("file", "a.html", "Name of the file to parse")
	flag.Parse()
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Error while reading the a.html file.")
	}
	links, _ := link.ParseHTML(file)
	for _, li := range links {
		fmt.Printf("Href:%s\nText:%s\n\n", li.Href, li.Data)
	}
}
