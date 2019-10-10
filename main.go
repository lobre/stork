// package main implements an algorithm of html content extraction.
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
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

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
	iterate(doc, func(n *html.Node, last *html.Node) {
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
	if err := a.stripContent(body); err != nil {
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

// stripContent will apply a first layer of cleaning to the parsed html.
//
// It will:
//  - remove unwanted tags
//  - remove comments
//  - apply a whitespace removal strategy
func (a *Article) stripContent(body *html.Node) error {
	spacing := regexp.MustCompile(`[ \r\n\t]+`)

	iterate(body, func(n *html.Node, last *html.Node) {
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

		case html.TextNode:
			if n.Parent != nil && (n.Parent.Data == "code" || n.Parent.Data == "pre") {
				return
			}

			// replace all whitespace characters to a single space
			n.Data = spacing.ReplaceAllString(n.Data, " ")

			// If after a block tag that ends
			if n.PrevSibling != nil && n.PrevSibling.Type == html.ElementNode && (BlockTags[n.PrevSibling.Data] || n.PrevSibling.Data == "br") {
				if strings.TrimSpace(n.Data) == "" {
					remove(n)
				}
				n.Data = strings.TrimLeft(n.Data, " ")
			}

			// remove leading space character if after a block element
			if last != nil && last.Type == html.ElementNode && (BlockTags[last.Data] || last.Data == "br") {
				if strings.TrimSpace(n.Data) == "" {
					remove(n)
				}
				n.Data = strings.TrimLeft(n.Data, " ")
			}

			// remove ending space character if last element of a block parent
			if n.Parent != nil && BlockTags[n.Parent.Data] && n.NextSibling == nil {
				n.Data = strings.TrimRight(n.Data, " ")
			}
		}
	})

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

func main() {
	var url string
	flag.StringVar(&url, "url", "", "url to parse")
	flag.Parse()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//inline := `<html><body><div id="toto" class="outter-class">
	//        <h1 class="inner-class">
	//	        The string I need
	//
	//	        <span class="other-class" >Some value I don't need</span>
	//	        <span class="other-class2" title="sometitle"></span>
	//            <script></script>
	//        </h1>
	//
	//        <pre>function toto()
	//        toto
	//          toto</pre>
	//
	//        <div class="other-class3">
	//            <h3>Some heading i don't need</h3>
	//        </div>
	//    </div></body></html>`

	//art, err := From(strings.NewReader(inline))
	art, err := From(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//md, err := art.Markdown()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(md)

	html, err := art.Html()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html)

	//fmt.Println(art.Text())
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

func iterate(doc *html.Node, do func(*html.Node, *html.Node)) {
	if doc == nil {
		return
	}

	var f func(n *html.Node, last *html.Node) *html.Node
	f = func(n *html.Node, last *html.Node) *html.Node {
		if n == nil {
			return last
		}

		do(n, last)

		last = n
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			last = f(c, last)
		}

		return last
	}
	f(doc, nil)
}

func remove(n *html.Node) {
	// save next because it is removed by RemoveChild
	// but we need it to continue iterating
	next := n.NextSibling
	n.Parent.RemoveChild(n)
	n.NextSibling = next
}
