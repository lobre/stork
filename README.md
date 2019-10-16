# stork

**Work in progress**

[Language independent content extraction from web pages](https://github.com/lobre/stork/blob/master/Language_Independent_Content_Extraction.pdf) is a paper that presents a simple, robust, accurate and language-independent solution for extracting the main content of an HTML-formatted Web page.

This package only provides an Golang implementation of the method given in the paper but **this is not affiliated with the research**.

It relies on `golang.org/x/net` to traverse HTML documents.

The core package of the extractor is available as an generic package providing an API that can be called by external packages.

This project also brings a cli to debug and try the extractor. It supports the generation of a graph to plot the document density.
For that purpose, it is using the simple package `github.com/guptarohit/asciigraph`.

## How to build and run for dev

    go build -o stork cmd/stork/main.go && ./stork -url "https://blog.golang.org/using-go-modules" -o text

## Related projects

 - https://github.com/advancedlogic/GoOse
 - https://github.com/thatguystone/swan

In order to stay in the same lexical field as these projects, I decided to call this implementation `stork`.
