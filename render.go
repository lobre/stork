package main

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// selfClosingTags is a list of void elements. Void elements
// are those that can't have any contents.
var selfClosingTags = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

// renderHtml is a simplest version of the one implemented in
// golang.org/x/net/html.Render.
// It does not escape special characters and adds newline
// characters after block elements.
func renderHtml(b *strings.Builder, n *html.Node) error {
	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("cannot render an ErrorNode node")
	case html.TextNode:
		if _, err := b.WriteString(n.Data); err != nil {
			return err
		}
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderHtml(b, c); err != nil {
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
	if selfClosingTags[n.Data] {
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
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := renderHtml(b, c); err != nil {
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
