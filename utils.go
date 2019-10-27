package stork

import "golang.org/x/net/html"

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

		return
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
