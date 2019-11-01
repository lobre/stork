package stork

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
	if err := renderHtml(&b, a.output, 0, true); err != nil {
		return "", err
	}
	return b.String(), nil
}

// renderHtml is a simplest version of the one implemented in
// golang.org/x/net/html.Render.
// It does not escape special characters and adds newline
// characters after block elements.
func renderHtml(b *strings.Builder, n *html.Node, depth int, indent bool) error {
	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("cannot render an ErrorNode node")
	case html.TextNode:
		_, err := b.WriteString(n.Data)
		return err
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderHtml(b, c, depth, indent); err != nil {
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
		if err := b.WriteByte('>'); err != nil {
			return err
		}
		// We don't render attributes
		_, err := b.WriteString(NewStr)
		return err
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
	if voidTags[n.Data] {
		if n.FirstChild != nil {
			return fmt.Errorf("void element <%s> has child nodes", n.Data)
		}
		_, err := b.WriteString("/>")
		return err
	}
	if err := b.WriteByte('>'); err != nil {
		return err
	}

	if inlineTags[n.Data] || n.Data == "pre" {
		indent = false
	}

	if n.FirstChild == nil {
		indent = false
	}

	// Render any child nodes.
	collapse := false
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if singleTxt := (c == n.FirstChild && c.NextSibling == nil && c.Type == html.TextNode); singleTxt {
			indent = false
		}

		if c.PrevSibling != nil {
			isPrevInline := c.PrevSibling.Type == html.ElementNode && inlineTags[c.PrevSibling.Data]
			isPrevTxt := c.PrevSibling.Type == html.TextNode
			isCurInline := c.Type == html.ElementNode && inlineTags[c.Data]
			isCurTxt := c.Type == html.TextNode

			if isCurTxt && isPrevInline || isCurInline && isPrevTxt {
				collapse = true
			}
		}

		if indent && !collapse {
			if _, err := b.WriteString(fmt.Sprint(NewStr, strings.Repeat(TabStr, depth+1))); err != nil {
				return err
			}

		}

		if err := renderHtml(b, c, depth+1, indent); err != nil {
			return err
		}
	}

	if indent {
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
