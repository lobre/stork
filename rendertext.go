package stork

import (
	"bytes"

	"golang.org/x/net/html"
)

// Text renders a text version of the article.
func (a *Article) Text() string {
	buf := bytes.Buffer{}

	iterate(a.output, func(n *html.Node) {
		switch n.Type {

		case html.ElementNode:
			if n.Data == "ul" || n.Data == "br" {
				buf.WriteString("\n")
				break
			}

			if n.Data == "li" {
				buf.WriteString("\n - ")
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
				buf.WriteString("\n\n")
			}

		case html.TextNode:
			buf.WriteString(n.Data)
		}

	})

	return buf.String()
}
