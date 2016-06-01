package main

import (
	"fmt"
	"strings"

	"github.com/zhengchun/selector/html"

	htmlx "golang.org/x/net/html"
)

func main() {
	s := `<p class="a">Links:<br/></p><ul><li><a href="foo">Foo</a></li><li><a href="bar">BarBaz</a></li></ul>`
	doc, _ := htmlx.Parse(strings.NewReader(s))

	var nodes = html.SelectNodes(doc, "//a[@href='foo']")
	for _, node := range nodes {
		fmt.Println(html.InnerText(node))
	}

}
