package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lobre/stork"
)

func main() {
	url := flag.String("url", "", "url to parse")
	file := flag.String("file", "", "file to parse")
	output := flag.String("o", "html", "output [html|markdown|text]")
	plot := flag.Bool("plot", false, "whether to plot the density")
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

	if *plot {
		fmt.Println(art.Plot())
		os.Exit(0)
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
