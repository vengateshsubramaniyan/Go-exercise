package main

import (
	"Go-exercise/Link-Parser/link"
	"encoding/xml"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type emptyStruct struct{}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}
type urlset struct {
	Urls  []loc  `xml:"url"`
	Xlmns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "Provide the url to create sitemap for...")
	maxDepth := flag.Int("depth", 3, "Maximum depth to traverse")
	flag.Parse()
	urls := bfs(*urlFlag, *maxDepth)
	var locs []loc
	for _, val := range urls {
		locs = append(locs, loc{val})
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	enc.Encode(urlset{locs, xmlns})
}

func bfs(urlstr string, maxDepth int) []string {

	q := make(map[string]emptyStruct)
	nq := map[string]emptyStruct{
		urlstr: emptyStruct{},
	}
	visit := make(map[string]emptyStruct)

	for i := 0; i <= maxDepth; i++ {
		q = nq
		nq = make(map[string]emptyStruct)
		for k := range q {
			if _, ok := visit[k]; ok {
				continue
			}
			visit[k] = emptyStruct{}
			reqURL, domain, paths := get(k)
			if reqURL == k {
				filterPaths := retriveDomainSpecificURLs(paths, domain)
				for _, link := range filterPaths {
					nq[link] = emptyStruct{}
				}
			}
		}
	}
	var ret []string
	for k := range visit {
		ret = append(ret, k)
	}
	return ret
}

func retriveDomainSpecificURLs(urlPaths []string, baseURL string) []string {
	var filterPaths []string
	for _, val := range urlPaths {
		if strings.HasPrefix(val, baseURL) {
			filterPaths = append(filterPaths, val)
		}
	}
	return filterPaths
}

func get(urlStr string) (string, string, []string) {

	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	paths := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}
	baseURL := paths.String()
	urlPaths := hrefs(resp.Body, baseURL)
	return resp.Request.URL.String(), baseURL, urlPaths
}

func hrefs(htmlData io.Reader, base string) []string {
	var ret []string
	links, _ := link.ParseHTML(htmlData)
	for _, li := range links {
		switch {
		case strings.HasPrefix(li.Href, "/") && len(li.Href) > 1:
			ret = append(ret, base+li.Href)
		case strings.HasPrefix(li.Href, "http"):
			ret = append(ret, li.Href)
		}
	}
	return ret
}
