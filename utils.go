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

func createNode(tag string, content string, attr []html.Attribute) *html.Node {
	elmtNode := html.Node{
		Type: html.ElementNode,
		Data: tag,
		Attr: attr,
	}
	textNode := html.Node{
		Type: html.TextNode,
		Data: content,
	}
	elmtNode.AppendChild(&textNode)
	return &elmtNode
}

func reverseNodes(nodes []*html.Node) {
	for i := len(nodes)/2 - 1; i >= 0; i-- {
		opp := len(nodes) - 1 - i
		nodes[i], nodes[opp] = nodes[opp], nodes[i]
	}
}

func ancestorsWithSameParent(n, p *html.Node) (*html.Node, *html.Node) {
	var nPath, pPath []*html.Node

	search := n
	for search != nil {
		nPath = append(nPath, search)
		search = search.Parent
	}

	if len(nPath) == 0 {
		return nil, nil
	}

	reverseNodes(nPath)

	search = p
	for search != nil {
		pPath = append(pPath, search)
		search = search.Parent
	}

	if len(pPath) == 0 {
		return nil, nil
	}

	reverseNodes(pPath)

	if nPath[0] != pPath[0] {
		return nil, nil
	}

	for i := 0; i < len(nPath) && i < len(pPath); i++ {
		if nPath[i] != pPath[i] {
			return nPath[i], pPath[i]
		}
	}

	return nil, nil
}
