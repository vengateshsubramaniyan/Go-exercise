package hn

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() (string, func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "[1,2,3]")
	})
	mux.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"by\":\"test_user\",\"descendants\":10,\"id\":1,\"kids\":[16732999,16729637,16729517,16729595],\"score\":34,\"time\":1522599083,\"title\":\"Test Story Title\",\"type\":\"story\",\"url\":\"https://www.test-story.com\"}")
	})
	server := httptest.NewServer(mux)
	return server.URL, func() {
		server.Close()
	}
}

func TestClient_GetTopItems(t *testing.T) {
	url, teardown := setup()
	defer teardown()
	c := Client{url}
	ids, err := c.GetTopNews()
	if err != nil {
		t.Errorf("GetTopItems() fails with error message %s", err.Error())
	}
	if len(ids) != 3 {
		t.Errorf("got %d, want 3", len(ids))
	}
}

func TestClient_GetItem(t *testing.T) {
	url, teardown := setup()
	defer teardown()
	c := Client{url}
	item, err := c.GetItem(1)
	if err != nil {
		t.Errorf("GetItem() fails with error message %s", err.Error())
	}
	if item.By != "test_user" {
		t.Errorf("got %s, want test_user", item.By)
	}
}

func TestDeafaultify(t *testing.T) {
	c := Client{}
	c.defaultify()
	if c.baseURL != baseURL {
		t.Errorf("got %s, but want %s", c.baseURL, baseURL)
	}
}
