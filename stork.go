// package stork implements an algorithm of html content extraction.
//
// It claims to bring a simple, robust, accurate and language-independent solution
// for extracting the main content of an HTML-formatted Web page and for
// removing additional content such as navigation menus, functional
// and design elements, and commercial advertisements.
//
// This method creates a text density graph of a given Web page and then
// selects the region of the Web page with the highest density.
//
// For more information about the original method, please have a look
// at the following paper.
//
// https://github.com/lobre/stork/raw/master/Language_Independent_Content_Extraction.pdf
//
// It provides here an implementation of the method given in the paper
// but is not affiliated with the research.
//
// Before analysing the html document, the process first applies some simple techniques
// to simplify the content.
//  - strip everything that is not in the body tag
//  - strip some unwanted tags
//  - apply a simple whitespace removal strategy
package stork

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/guptarohit/asciigraph"
	"golang.org/x/net/html"
)

// leashParams represents parameters needed
// to calculate a leash from a text size.
type leashParams struct {
	minLength, maxLength float64
	minLeash, maxLeash   float64
}

// default values for leash calculation
var defaultLeashParams = leashParams{0, 400, 0, 40}

var htmlSkeleton string = "<!DOCTYPE html><html><head><meta charset=\"utf-8\" /><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" /></head><body></body></html>"

// blockText stores the textual representation of
// a structural block element on an html page.
// It aims to be used in a slice to
// calculate the density.
type blockText struct {
	block *html.Node
	text  string
}

// Article contains all the extracted values of an html document.
// It should be created using the From() method.
type Article struct {
	// A header image to use for the article
	Thumbnail *Image

	// Images contained in the extracted article
	Images []*Image

	// Links contained in the extracted article
	Links []string

	// Metadata taken from the html document
	Meta struct {
		Lang        string
		Canonical   string
		Title       string
		Favicon     string
		Description string
		Keywords    string
		OpenGraph   map[string]string
	}

	density []blockText

	output *html.Node
}

// Image contains information taken from a <img> html tag.
type Image struct {
	Src    string
	Width  uint
	Height uint
	node   *html.Node
}

// From parses an html document from an io.Reader
// and extracts the content into an Article.
func From(r io.Reader) (*Article, error) {
	var a Article

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	// search body
	var body *html.Node
	iterate(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
		}
	})

	if body == nil {
		return nil, errors.New("body not found")
	}

	// TODO
	if err := a.extractMetadata(doc); err != nil {
		return nil, err
	}

	// TODO
	if err := a.extractThumbnail(doc); err != nil {
		return nil, err
	}

	if err := a.clean(body); err != nil {
		return nil, err
	}

	// TODO
	// this should create
	if err := a.calculateDensity(body); err != nil {
		return nil, err
	}

	// TODO
	if err := a.extractContent(body); err != nil {
		return nil, err
	}

	// TODO parse article images

	// TODO assert if really an article

	return &a, nil
}

func (a *Article) extractThumbnail(doc *html.Node) error {
	return nil
}

func (a *Article) extractMetadata(doc *html.Node) error {
	return nil
}

// clean will apply a first layer of cleaning to the parsed html.
//
// It will:
//  - remove unwanted tags
//  - remove comments
//  - apply a whitespace removal strategy
func (a *Article) clean(body *html.Node) error {
	iterate(body, func(n *html.Node) {
		switch n.Type {

		case html.CommentNode:
			remove(n)

		case html.ElementNode:
			if IgnoreTags[n.Data] {
				remove(n)
			}

			// remove class
			var keep []html.Attribute
			for _, attr := range n.Attr {
				if attr.Key != "class" {
					keep = append(keep, attr)
				}
			}
			n.Attr = keep
		}
	})

	minify(body)

	return nil
}

func (a *Article) calculateDensity(body *html.Node) error {
	a.density = nil
	a.density = append(a.density, blockText{body, ""})
	idx := 0

	iterate(body, func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			if blockTags[n.Data] {
				a.density = append(a.density, blockText{n, ""})
				idx++
			}
		case html.TextNode:
			a.density[idx].text += n.Data
		}
	})

	return nil
}

// body parameter is temporary while the density is not implemented
func (a *Article) extractContent(body *html.Node) error {
	if len(a.density) <= 0 {
		return errors.New("wrong density")
	}

	// find longest text
	smax, maxl := 0, 0
	for i, d := range a.density {
		if len(d.text) > maxl {
			smax = i
			maxl = len(d.text)
		}
	}

	// high density region
	hdr := []int{smax}

	restart := true
	for restart {
		restart = false
		for i, d := range a.density {
			add := false
			for _, j := range hdr {
				// already exists
				if i == j {
					add = false
					break
				}
				leash := calculateLeash(defaultLeashParams, len(d.text))
				if abs(i-j) < leash {
					add = true
				}
			}
			if add {
				hdr = append(hdr, i)
				restart = true
			}
		}
	}

	start, end := smax, smax
	for _, i := range hdr {
		if i < start {
			start = i
		}
		if i > end {
			end = i
		}
	}

	if err := a.assembleOutput(start, end); err != nil {
		return err
	}

	return nil
}

func (a *Article) assembleOutput(start, end int) error {
	reader := strings.NewReader(htmlSkeleton)

	var err error
	a.output, err = html.Parse(reader)
	if err != nil {
		return err
	}

	// search head and body
	var root, head, body *html.Node
	iterate(a.output, func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "html":
				root = n
			case "head":
				head = n
			case "body":
				body = n
			}
		}
	})

	if root == nil || head == nil || body == nil {
		return errors.New("error parsing html skeleton")
	}

	if a.Meta.Lang != "" {
		root.Attr = append(root.Attr, html.Attribute{Key: "lang", Val: a.Meta.Lang})
	}

	if a.Meta.Description != "" {
		head.AppendChild(createNode("meta", "", []html.Attribute{
			html.Attribute{Key: "name", Val: "description"},
			html.Attribute{Key: "content", Val: a.Meta.Description},
		}))
	}

	if a.Meta.Keywords != "" {
		head.AppendChild(createNode("meta", "", []html.Attribute{
			html.Attribute{Key: "name", Val: "keywords"},
			html.Attribute{Key: "content", Val: a.Meta.Keywords},
		}))
	}

	for k, v := range a.Meta.OpenGraph {
		head.AppendChild(createNode("meta", "", []html.Attribute{
			html.Attribute{Key: "property", Val: "og:" + k},
			html.Attribute{Key: "content", Val: v},
		}))
	}

	if a.Meta.Title != "" {
		head.AppendChild(createNode("title", a.Meta.Title, nil))
		body.AppendChild(createNode("h1", a.Meta.Title, nil))
	}

	if a.Meta.Canonical != "" {
		head.AppendChild(createNode("link", "", []html.Attribute{
			html.Attribute{Key: "rel", Val: "canonical"},
			html.Attribute{Key: "href", Val: a.Meta.Canonical},
		}))
	}

	if a.Meta.Favicon != "" {
		head.AppendChild(createNode("link", "", []html.Attribute{
			html.Attribute{Key: "rel", Val: "shortcut icon"},
			html.Attribute{Key: "href", Val: a.Meta.Favicon},
			html.Attribute{Key: "type", Val: "image/x-icon"},
		}))
	}

	if a.Thumbnail != nil {
		// shallow copy
		thumb := *a.Thumbnail.node
		thumb.Parent, thumb.PrevSibling, thumb.NextSibling = nil, nil, nil
		body.AppendChild(&thumb)
	}

	// append article content
	idx := start
	for idx <= end {
		// shallow copy
		block := *a.density[idx].block
		block.Parent, block.PrevSibling, block.NextSibling = nil, nil, nil
		body.AppendChild(&block)
		idx++
	}

	return nil
}

// Density returns the content from the density table
// with index and text length prepended to each text item.
// This function is for debug purposes.
func (a *Article) Density() string {
	var b strings.Builder
	for i, d := range a.density {
		b.WriteString(fmt.Sprintf("%d (%d) - %s\n", i, len(d.text), d.text))
	}
	return b.String()
}

// Plot will draw the density graph calculated
// from the extracted article.
//
// It will generate a graph similar to the one on figure 2 at page 3 of the paper.
// https://github.com/lobre/stork/raw/master/Language_Independent_Content_Extraction.pdf
func (a *Article) Plot() string {
	var data []float64
	for _, t := range a.density {
		data = append(data, float64(len(t.text)))
	}
	return asciigraph.Plot(data, asciigraph.Height(30))
}

func calculateLeash(lp leashParams, length int) int {
	var res float64
	res = (((lp.maxLeash - lp.minLeash) * (float64(length) - lp.minLength)) /
		(lp.maxLength - lp.minLength)) + lp.minLeash

	if res < lp.minLeash {
		res = lp.minLeash
	}

	if res > lp.maxLeash {
		res = lp.maxLeash
	}

	return int(res)
}
