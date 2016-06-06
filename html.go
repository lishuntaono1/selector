package selector

import (
	"bytes"

	"golang.org/x/net/html"
)

func SelectHtmlNodes(n *html.Node, xpath string) []*html.Node {
	var nav = &NodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	nodes := make([]*html.Node, 0)

	for iter.MoveNext() {
		cur := iter.Current().(*NodeNavigator)
		nodes = append(nodes, cur.currnode)
	}
	return nodes
}

func SelectSingleHtmlNode(n *html.Node, xpath string) *html.Node {
	var nav = &NodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	for iter.MoveNext() {
		cur := iter.Current().(*NodeNavigator)
		return cur.currnode
	}
	return nil
}

func InnerText(n *html.Node) string {
	var b bytes.Buffer
	var output func(*html.Node)
	output = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			output(child)
		}
	}
	output(n)
	return b.String()
}
