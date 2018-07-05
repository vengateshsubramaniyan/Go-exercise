package main

import (
	"Go-exercise/Quiet_hn/hn"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type item struct {
	hn.Item
	Host  string
	index int
}

type cachedData struct {
	numberStories int
	cache         []item
	duration      time.Duration
	expiration    time.Time
	mutex         sync.Mutex
}

func (c *cachedData) cacheStories() ([]item, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if time.Now().Sub(c.expiration) < 0 {
		return c.cache, nil
	}
	var err error
	c.cache, err = getTopStories(c.numberStories)
	if err != nil {
		return nil, err
	}
	c.expiration = time.Now().Add(6 * time.Minute)
	return c.cache, nil
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

var tp1 *template.Template

func init() {
	tp1 = template.Must(template.ParseFiles("index.html"))
}

func main() {
	port := flag.Int("port", 8081, "Port to run this application")
	noOfRecords := flag.Int("noOfRecords", 30, "Number of records to display")
	flag.Parse()

	http.HandleFunc("/", handler(*noOfRecords))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handler(noOfRecords int) http.HandlerFunc {
	sc := cachedData{numberStories: noOfRecords, duration: 6 * time.Second}

	go func() {
		tick := time.NewTicker(6 * time.Second)
		for {
			temp := cachedData{numberStories: noOfRecords, duration: 6 * time.Second}
			temp.cacheStories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-tick.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := sc.cacheStories()
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		result := templateData{Stories: stories, Time: time.Now().Sub(start)}
		if tp1.Execute(w, result) != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	})
}

func isStoryLink(hnItem hn.Item) bool {
	return hnItem.Type == "story" && hnItem.URL != ""
}

func getTopStories(noOfRecords int) ([]item, error) {
	var client hn.Client
	ids, err := client.GetTopNews()
	if err != nil {
		return nil, err
	}
	var stories []item
	itemChan := make(chan *item)

	for i := 0; i < noOfRecords*5/4; i++ {
		go getItem(itemChan, ids[i], i)
	}

	for i := 0; i < noOfRecords*5/4; i++ {
		ret := <-itemChan
		if ret != nil {
			stories = append(stories, *ret)
		}
	}
	sort.Slice(stories, func(i, j int) bool { return stories[i].index < stories[j].index })
	return stories[:noOfRecords], nil
}

func getItem(ch chan *item, id int, idx int) {
	var client hn.Client
	it, err := client.GetItem(id)
	if err == nil && isStoryLink(it) {
		ite := parseURL(it, idx)
		ch <- &ite
		return
	}
	ch <- nil
}

func parseURL(hnItem hn.Item, idx int) item {
	var ret item
	ret.Item = hnItem
	ret.index = idx
	ul, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(ul.Hostname(), "www.")
	}
	return ret
}
