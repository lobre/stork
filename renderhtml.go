package main

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var (
	TabStr = "    "
	NewStr = "\n"
)

// Html renders the article in html format.
func (a *Article) Html() (string, error) {
	var b strings.Builder
	if err := renderHtml(&b, a.output, 0); err != nil {
		return "", err
	}
	return b.String(), nil
}

// renderHtml is a simplest version of the one implemented in
// golang.org/x/net/html.Render.
// It does not escape special characters and adds newline
// characters after block elements.
func renderHtml(b *strings.Builder, n *html.Node, depth int) error {
	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("cannot render an ErrorNode node")
	case html.TextNode:
		switch {
		case n.Parent != nil && (n.Parent.Data == "pre" || n.Parent.Data == "code"),
			n.PrevSibling == nil && n.NextSibling == nil:

			_, err := b.WriteString(n.Data)
			return err
		default:
			_, err := b.WriteString(fmt.Sprint(strings.Repeat(TabStr, depth), n.Data))
			return err
		}
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderHtml(b, c, depth+1); err != nil {
				return err
			}
		}
		return nil
	case html.ElementNode:
		// No-op.
	case html.CommentNode:
		// We don't render comment nodes
		return nil
	case html.DoctypeNode:
		if _, err := b.WriteString("<!DOCTYPE "); err != nil {
			return err
		}
		if _, err := b.WriteString(n.Data); err != nil {
			return err
		}
		// We don't render attributes
		return b.WriteByte('>')
	default:
		return errors.New("unknown node type")
	}

	if _, err := b.WriteString(strings.Repeat(TabStr, depth)); err != nil {
		return err
	}

	// Render the <xxx> opening tag.
	if err := b.WriteByte('<'); err != nil {
		return err
	}
	if _, err := b.WriteString(n.Data); err != nil {
		return err
	}
	for _, a := range n.Attr {
		if err := b.WriteByte(' '); err != nil {
			return err
		}
		if a.Namespace != "" {
			if _, err := b.WriteString(a.Namespace); err != nil {
				return err
			}
			if err := b.WriteByte(':'); err != nil {
				return err
			}
		}
		if _, err := b.WriteString(a.Key); err != nil {
			return err
		}
		if _, err := b.WriteString(`="`); err != nil {
			return err
		}
		if _, err := b.WriteString(a.Val); err != nil {
			return err
		}
		if err := b.WriteByte('"'); err != nil {
			return err
		}
	}
	if VoidTags[n.Data] {
		if n.FirstChild != nil {
			return fmt.Errorf("void element <%s> has child nodes", n.Data)
		}
		_, err := b.WriteString("/>")
		return err
	}
	if err := b.WriteByte('>'); err != nil {
		return err
	}

	// Render any child nodes.
	inline := false
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch {
		case c == n.FirstChild && c.NextSibling == nil && c.Type == html.TextNode,
			n.Data == "pre" || n.Data == "code":

			inline = true
			break

		default:
			if _, err := b.WriteString(NewStr); err != nil {
				return err
			}
		}

		if err := renderHtml(b, c, depth+1); err != nil {
			return err
		}
	}

	if !inline {
		if _, err := b.WriteString(fmt.Sprint(NewStr, strings.Repeat(TabStr, depth))); err != nil {
			return err
		}
	}

	// Render the </xxx> closing tag.
	if _, err := b.WriteString("</"); err != nil {
		return err
	}
	if _, err := b.WriteString(n.Data); err != nil {
		return err
	}
	return b.WriteByte('>')
}
