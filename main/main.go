package main

import (
	"fmt"
	"strings"

	"github.com/zhengchun/selector"
	"golang.org/x/net/html"
)

func main() {
	s := `<html><head></head><body><a href="/"><img src="logo.jpg" /></a><ul><li class="ball"><a href="/football">football</a><a href="/basketball">basketball</a><a href="/ping_pong">ping-pong</a></li><li><a href="/swiming">swiming</a><a href="/running">running</a></li></ul></body></html>`
	doc, _ := html.Parse(strings.NewReader(s))

	var nodes = selector.SelectHtmlNodes(doc, "//a[last()]")
	for _, node := range nodes {
		fmt.Println(selector.InnerText(node))
	}
}
