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
package main

import (
	"bytes"
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

// IgnoreTags are tags that will be removed from the document before analysing it.
// This list contains tags such as metadata elements that don't make sense in
// the context of the extracted content.
//
// Note:
// Fresh and clean new metadata will be added afterwards if using the Html() method
// in order to re-create a full and complete html document.
//
// https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/Content_categories#Metadata_content
//
// Other tags can be added to this list if they should be removed from the extracted document.
var IgnoreTags = []string{"base", "command", "link", "meta", "noscript", "script", "style", "title"}

// BlockTags are elements that always start on a new line and takes up the full width available
// They are used to determine what is a structural tag in order to extract the main content of the page.
//
// https://www.w3schools.com/html/html_blocks.asp
var BlockTags = []string{
	"address",
	"article",
	"aside",
	"blockquote",
	"canvas",
	"dd",
	"div",
	"dl",
	"dt",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"header",
	"hr",
	"li",
	"main",
	"nav",
	"noscript",
	"ol",
	"p",
	"pre",
	"section",
	"table",
	"tfoot",
	"ul",
	"video",
}

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

	// The parent node representing the overall html document
	Doc *html.Node

	// An html node representing the body
	body *html.Node
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
	var (
		a   Article
		err error
	)

	a.Doc, err = html.Parse(r)
	if err != nil {
		return nil, err
	}

	// search body
	iterate(a.Doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			a.body = n
		}
	})

	if a.body == nil {
		return nil, errors.New("body not found")
	}

	// TODO(lobre) parse metadata
	if err := a.metadata(); err != nil {
		return nil, err
	}

	// TODO(lobre) parse image
	if err := a.image(); err != nil {
		return nil, err
	}

	// TODO(lobre) strip useless bits of the document
	if err := a.strip(); err != nil {
		return nil, err
	}

	// TODO(lobre) extract document in some sort of data structure
	if err := a.extract(); err != nil {
		return nil, err
	}

	// TODO(lobre) parse article images

	// TODO assert if really an article

	return &a, nil
}

func (a *Article) image() error {
	return nil
}

func (a *Article) metadata() error {
	return nil
}

func (a *Article) strip() error {
	spacing := regexp.MustCompile(`[ \r\n\t]+`)

	iterate(a.body, func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			// remove unwanted tags
			for _, ignore := range IgnoreTags {
				if n.Type == html.ElementNode && n.Data == ignore {
					remove(n)
				}
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
			if n.Parent.Data != "code" && n.Parent.Data != "pre" {
				n.Data = strings.TrimSpace(spacing.ReplaceAllString(n.Data, " "))
				if n.Data == "" {
					remove(n)
				}
			}
		}
	})

	return nil
}

func (a *Article) extract() error {
	// analyse structural tags "p", "table", "br", "div", "h1", "h2", "h3", "h4", "h5", "h6", "li"
	// parse into a data structure that will easily allow outputs
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

	html, err := art.Html()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(html)
	//fmt.Println(art.Plot())
}

func (a *Article) Text() string {
	return ""
}

func (a *Article) Html() (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	if err := html.Render(w, a.body); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *Article) Markdown() (string, error) {
	return "", nil
}

// Plot will draw the density graph calculated
// from the extracted article.
//
// It will generate a graph alike what we can find on figure 2 at page 3 of the paper.
// https://github.com/lobre/stork/raw/master/Language_Independent_Content_Extraction.pdf
func (a *Article) Plot() string {
	data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	return asciigraph.Plot(data, asciigraph.Height(30))
}

func iterate(doc *html.Node, do func(*html.Node)) {
	if doc == nil {
		return
	}

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}

		do(n)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func remove(n *html.Node) {
	// save next because it is removed by RemoveChild
	// but we need it to continue iterating
	next := n.NextSibling
	n.Parent.RemoveChild(n)
	n.NextSibling = next
}
