package html

import (
	"bytes"

	"golang.org/x/net/html"
)

func InnerText(n *html.Node) string {
	var b bytes.Buffer
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			b.WriteString(child.Data)
		}
	}
	return b.String()
}
