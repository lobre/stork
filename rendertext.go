package stork

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Text renders a text version of the article.
func (a *Article) Text() string {
	b := strings.Builder{}

	var body *html.Node
	iterate(a.output, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
		}
	})

	if body == nil {
		return ""
	}

	iterate(body, func(n *html.Node) {
		switch n.Type {

		case html.ElementNode:

			switch n.Data {
			case "ul", "br", "p", "pre",
				"h1", "h2", "h3", "h4", "h5", "h6":
				b.WriteString("\n")
			case "div":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					b.WriteString("\n")
				}
			}

			if blockTags[n.Data] {
				b.WriteString("\n")
			}

			if n.Data == "li" {
				b.WriteString(" - ")
			}

		case html.TextNode:
			b.WriteString(n.Data)
		}

	})

	text := strings.TrimSpace(b.String())

	regex := regexp.MustCompile("\n{3,}")
	return regex.ReplaceAllString(text, "\n\n")
}
