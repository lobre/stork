# stork

**Work in progress**

[Language independent content extraction from web pages](https://github.com/lobre/stork/blob/master/Language_Independent_Content_Extraction.pdf) is a paper that presents a simple, robust, accurate and language-independent solution for extracting the main content of an HTML-formatted Web page.

This package only provides an Golang implementation of the method given in the paper but **this is not affiliated with the research**.

It relies on `golang.org/x/net` to traverse HTML documents.

The core package of the extractor is available as an generic package providing an API that can be called by external packages.

This project also brings a cli to debug and try the extractor. It supports the generation of a graph to plot the document density.
For that purpose, it is using the simple package `github.com/guptarohit/asciigraph`.

## Algorithm adaptation

By implementing the exact same algorithm demonstrated in the paper, the results were not always precise.
I decided to update a part of the algorithm in order to try to improve it.

The algorithm is initially using two constants.

    c1 = 0.333
    c2 = 4

If trying to semantically explain these constants, `c1` would be used to calculate a cutoff as `cutoff = c1 * maxLength` where `maxLenght` is the longest parsed string in the article.
Where iterating on the document, all the strings whose length is lower to the cutoff would not be marked as content.

`c2` is another metric that I would call: "leash". It defines that if a text is long enough to pass the cutoff, it will have to be situated up to 4 slots from another text that is part of the content. Otherwise, this text would not be considered as content.

### Improvements

I first started to play with these variables to try to have better results. But either I was too strict or too gentle and so depending on the document, I had either too much or too few content.
That's why I decided to change a little bit the behavior. I deleted the principle of cutoff and created a function to dynamically calculate the leash from the text length.

I defined limits as constants as followed:

    minLength = 0
    maxLength = 400
    minLeash = 0
    maxLeash = 40

And then I used the following formula to get a leash value from a text length (re-scale operation).

    x = ( ( (maxLeash - minLeash) * (y - minLength) ) / maxLength - minLength ) + minLeash

Where `y` is the length of the text of which we want to get the leash value. If the calculated leash value is below `minLeash`, we set it to `minLeash` and if it is above `maxLeash`, we set it to `maxLeash`.

With this new behavior in the algorithm, we have better results. And if I try again to explain with words what this dynamic leash calculation provides, I would say: the longer a text is, the more likely it belong to the main content. So we put some weight to the leash to say that it can be farther from another text marked as content.

## How to build and run for dev

    go build -o stork cmd/stork/main.go && ./stork -url "https://blog.golang.org/using-go-modules" -o text

## Related projects

 - https://github.com/advancedlogic/GoOse
 - https://github.com/thatguystone/swan

In order to stay in the same lexical field as these projects, I decided to call this implementation `stork`.
