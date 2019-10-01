package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/guptarohit/asciigraph"
)

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

	// TODO(lobre) parse document in some sort of data structure
	if err := a.parse(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *Article) image() error {
	return nil
}

func (a *Article) metadata() error {
	return nil
}

func (a *Article) parse() error {
	// init clip from body
	a.clip = a.Doc.Find("body")

	// remove unwanted tags "link", "script", "style", "meta"
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

	art, err := From(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(art.Text())
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
