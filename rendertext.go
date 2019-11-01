package stork

import (
	"strings"

	"golang.org/x/net/html"
)

// Text renders a text version of the article.
func (a *Article) Text() string {
	b := strings.Builder{}

	iterate(a.output, func(n *html.Node) {
		switch n.Type {

		case html.ElementNode:
			if n.Data == "ul" || n.Data == "br" {
				b.WriteString("\n")
				break
			}

			if n.Data == "li" {
				b.WriteString("\n - ")
				break
			}

			if n.Data == "div" {
				if n.FirstChild == nil {
					break
				}
				if n.FirstChild.Type != html.TextNode {
					break
				}
			}

			if blockTags[n.Data] {
				b.WriteString("\n\n")
			}

		case html.TextNode:
			b.WriteString(n.Data)
		}

	})

	return b.String()
}
