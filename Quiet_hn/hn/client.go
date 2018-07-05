//Package hn creates a wrapper to access the hacker news
package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://hacker-news.firebaseio.com/v0"
)

//Client struct will be accessible to the end users.
type Client struct {
	baseURL string
}

func (c *Client) defaultify() {
	if c.baseURL == "" {
		c.baseURL = baseURL
	}
}

//GetTopNews returns top news from the Hacker news website.
func (c *Client) GetTopNews() ([]int, error) {
	c.defaultify()
	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", c.baseURL))
	if err != nil {
		return nil, err
	}
	var ids []int
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

//GetItem fetch particular item from the hacker news and return that to the client
func (c *Client) GetItem(id int) (Item, error) {
	c.defaultify()
	var item Item
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.baseURL, id))
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

//Item represents single item returned by the hacker news
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}
