package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/guptarohit/asciigraph"
	"golang.org/x/net/html"
)

type extractor struct {
	l []string
}

func main() {
	var url string
	flag.StringVar(&url, "url", "", "url to parse")
	flag.Parse()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// initialize the linear list with an empty string
	e := extractor{l: []string{""}}
	e.traverse(resp.Body)
	e.plot()
}

func (e *extractor) traverse(r io.Reader) {
	z := html.NewTokenizer(r)
	ignore := false

	for {
		tt := z.Next()

		switch tt {

		case html.ErrorToken:
			return

		case html.TextToken:
			if ignore {
				continue
			}

			token := z.Token()
			e.l[len(e.l)-1] = e.l[len(e.l)-1] + token.Data

		case html.StartTagToken, html.SelfClosingTagToken:
			if ignore {
				continue
			}

			token := z.Token()

			if isIgnored(token.Data) {
				ignore = true
				continue
			}

			if isStructural(token.Data) {
				e.l = append(e.l, "")
			}

		case html.EndTagToken:
			token := z.Token()

			if isIgnored(token.Data) {
				ignore = false
			}
		}
	}
}

func (e *extractor) plot() {
	var data []float64

	for _, s := range e.l {
		data = append(data, float64(len(s)))
	}

	graph := asciigraph.Plot(data, asciigraph.Height(30), asciigraph.Width(60))
	graph2 := asciigraph.Plot(data, asciigraph.Height(30))
	fmt.Println(graph)
	fmt.Println(graph2)
}

func isStructural(tag string) bool {
	tags := []string{
		"p",
		"table",
		"br",
		"div",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"li",
	}

	for _, t := range tags {
		if tag == t {
			return true
		}
	}

	return false
}

func isIgnored(tag string) bool {
	tags := []string{
		"link",
		"script",
		"style",
		"meta",
	}

	for _, t := range tags {
		if tag == t {
			return true
		}
	}

	return false
}
