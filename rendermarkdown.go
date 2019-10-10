package main

import (
	"bytes"
	"io"

	md "github.com/evorts/html-to-markdown"
	"golang.org/x/net/html"
)

// Markdown renders the article in markdown format.
func (a *Article) Markdown() (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	if err := html.Render(w, a.output); err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html.UnescapeString(buf.String()))
	if err != nil {
		return "", err
	}

	return markdown, nil
}
