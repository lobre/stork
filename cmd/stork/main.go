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
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lobre/stork"
)

func main() {
	url := flag.String("url", "", "url to parse")
	file := flag.String("file", "", "file to parse")
	output := flag.String("o", "html", "output [html|markdown|text]")
	flag.Parse()

	var art *stork.Article

	if *url != "" {
		resp, err := http.Get(*url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		art, err = stork.From(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *file != "" {
		f, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}

		art, err = stork.From(bytes.NewReader(f))
		if err != nil {
			log.Fatal(err)
		}
	}

	if art == nil {
		log.Fatal("nothing to process")
	}

	switch *output {

	case "html":
		html, err := art.Html()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(html)

	case "markdown":
		md, err := art.Markdown()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(md)

	case "text":
		fmt.Println(art.Text())
	}
}
