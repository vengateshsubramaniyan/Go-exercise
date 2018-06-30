package link

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link is used to store the struct data.
type Link struct {
	Href string
	Data string
}

//ParseHTML is used to parse the html file.
func ParseHTML(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, errors.New("Error while parsing html")
	}
	nodes := linkNodes(node)
	var linkes []Link
	for _, n := range nodes {
		var tlink Link
		for _, att := range n.Attr {
			if att.Key == "href" {
				tlink.Href = att.Val
				break
			}
		}
		tlink.Data = strings.TrimSpace(text(n))
		linkes = append(linkes, tlink)
	}
	return linkes, nil
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var ret string

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
