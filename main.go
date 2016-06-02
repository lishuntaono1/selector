package main

import (
	"fmt"
	"strings"

	"github.com/zhengchun/selector/html"

	htmlx "golang.org/x/net/html"
)

func main() {
	s := `<a href="foo">Foo</a><a href="bar">BarBaz</a>`
	doc, _ := htmlx.Parse(strings.NewReader(s))

	var nodes = html.SelectNodes(doc, "//a[@href='bar']")

	fmt.Println("===========")
	for _, node := range nodes {
		fmt.Println(html.InnerText(node))
	}

}
