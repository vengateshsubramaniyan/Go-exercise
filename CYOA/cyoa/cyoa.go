package cyoa

import (
	"html/template"
	"net/http"
)

//Stories is a map with key string and value type chapter
type Stories map[string]Chapter

//Chapter to store title,story, options details
type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

//Option to store text and arc details
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	s   Stories
	tp1 *template.Template
}

//ServeHttp to handle the request
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	if chapter, ok := h.s[path]; ok {
		err := h.tp1.Execute(w, chapter)
		if err != nil {
			http.Error(w, "Internal server error...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Story not found..", http.StatusBadRequest)
}

//StoryHandler to hadler the request
func StoryHandler(s Stories, tp1 *template.Template) handler {
	return handler{s, tp1}
}
