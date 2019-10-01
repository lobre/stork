package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/guptarohit/asciigraph"
	"golang.org/x/net/html"
)

var IgnoreTags = []string{"script", "link", "style", "meta"}
var StructuralTags = []string{"p", "table", "br", "div", "li", "h1", "h2", "h3", "h4", "h5", "h6"}

type Article struct {
	Img *Image

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

	Doc *goquery.Document

	clip *goquery.Selection
}

type Image struct {
	Src        string
	Width      uint
	Height     uint
	Bytes      int64
	Confidence uint
	Sel        *goquery.Selection
}

// From must be called before Html, Text or Markdown.
func From(r io.Reader) (*Article, error) {
	var (
		a   Article
		err error
	)

	a.Doc, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
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
	// init clip from body
	a.clip = a.Doc.Find("body")

	// remove unwanted tags
	a.clip.Find(strings.Join(IgnoreTags, ",")).Remove()

	// remove classes and id
	a.clip.Find("*").RemoveClass().RemoveAttr("id")

	// remove blank lines

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

	inline := `<html><body><div class="outter-class">
        <h1 class="inner-class">
	        The string I need

	        <span class="other-class" >Some value I don't need</span>
	        <span class="other-class2" title="sometitle"></span>
        </h1>

        <div class="other-class3">
            <h3>Some heading i don't need</h3>
        </div>
    </div></body></html>`

	art, err := From(strings.NewReader(inline))
	//art, err := From(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	html, err := art.Html()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(html)
	fmt.Println(art.Plot())
}

func (a *Article) Text() string {
	return a.clip.Text()
}

func (a *Article) Html() (string, error) {
	return a.clip.Html()
}

func (a *Article) Markdown() (string, error) {
	return "", nil
}

func (a *Article) Plot() string {
	data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	return asciigraph.Plot(data, asciigraph.Height(30))
}

func iterate(doc *html.Node, do func(n *html.Node) error) error {
	var f func(n *html.Node) error
	f = func(n *html.Node) error {
		if err := do(n); err != nil {
			return err
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}

	if err := f(doc); err != nil {
		return err
	}

	return nil
}
