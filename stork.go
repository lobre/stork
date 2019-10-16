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
	"io"

	"github.com/guptarohit/asciigraph"
	"golang.org/x/net/html"
)

// Article contains all the extracted values of an html document.
// It should be created using the From() method.
type Article struct {
	// A header image to use for the article
	Thumbnail *Image

	// Images contained in the extracted article
	Images []*Image

	// All metadata taken from the html document
	Meta struct {
		Authors     []string
		Canonical   string
		Description string
		Domain      string
		Favicon     string
		Keywords    string
		Links       []string
		Lang        string
		OpenGraph   map[string]string
		PublishDate string
		Tags        []string
		Title       string
	}

	density []struct {
		block *html.Node
		text  string
	}

	output *html.Node
}

// Image contains information taken from a <img> html tag.
type Image struct {
	Src        string
	Width      uint
	Height     uint
	Bytes      int64
	Confidence uint
	Node       *html.Node
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

	// TODO
	if err := a.clean(body); err != nil {
		return nil, err
	}

	// TODO
	// this should create
	if err := a.calculateDensity(body); err != nil {
		return nil, err
	}

	// TODO
	if err := a.generateArticle(body); err != nil {
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
	// fill the density slice
	return nil
}

// body parameter is temporary while the density is not implemented
func (a *Article) generateArticle(body *html.Node) error {
	// initiale node with an html skeleton
	// generate metadata nodes
	// generate thumbnail node
	// calculate article boundaries with density map
	// append relevant tags to article

	// temp
	a.output = body

	return nil
}

// Plot will draw the density graph calculated
// from the extracted article.
//
// It will generate a graph similar to the one on figure 2 at page 3 of the paper.
// https://github.com/lobre/stork/raw/master/Language_Independent_Content_Extraction.pdf
func (a *Article) Plot() string {
	data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	return asciigraph.Plot(data, asciigraph.Height(30))
}
