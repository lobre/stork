package stork

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func minify(n *html.Node) {
	spacing := regexp.MustCompile(`[ \r\n\t]+`)

	iterate(n, func(n *html.Node) {
		switch n.Type {

		case html.TextNode:
			if n.Parent != nil && (n.Parent.Data == "code" || n.Parent.Data == "pre") {
				break
			}

			// replace all whitespace characters to a single space
			n.Data = spacing.ReplaceAllString(n.Data, " ")
		}
	})

	// trim according to inline/block rules
	trimLeft(n, true)
	trimRight(n, true)
}

func trimLeft(n *html.Node, trim bool) bool {
	switch n.Type {
	case html.ElementNode:
		if blockTags[n.Data] {
			trim = true
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			trim = trimLeft(c, trim)
		}

		if blockTags[n.Data] || n.Data == "br" {
			return true
		}

	case html.TextNode:
		if trim {
			n.Data = strings.TrimLeft(n.Data, " \r\n\t")

			if strings.TrimSpace(n.Data) == "" {
				remove(n)
				return trim
			}
		}

		return false
	}

	return trim
}

func trimRight(n *html.Node, trim bool) bool {
	switch n.Type {
	case html.ElementNode:
		if blockTags[n.Data] {
			trim = true
		}

		for c := n.LastChild; c != nil; c = c.PrevSibling {
			trim = trimRight(c, trim)
		}

		if blockTags[n.Data] || n.Data == "br" {
			return true
		}

	case html.TextNode:
		if trim {
			n.Data = strings.TrimRight(n.Data, " \r\n\t")

			if strings.TrimSpace(n.Data) == "" {
				remove(n)
				return trim
			}
		}

		return false
	}

	return trim
}
