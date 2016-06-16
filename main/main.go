package main

import (
	"fmt"
	"strings"

	"github.com/zhengchun/selector"
	"golang.org/x/net/html"
)

func main() {
	s := `
	<html>
    <head></head>
    <body>
    <a href="/"><img src="logo.jpg" /></a>
    <ul>
        <li class="ball">
            <a id="1" href="/football">football</a>
            <a id="2" href="/basketball">basketball</a>
            <a id="3" href="/ping_pong">ping-pong</a>
        </li>
        <li>
            <a id="4" href="/swiming">swiming</a>
            <a id="5" href="/running">running</a>
        </li>
    </ul> 
    </body>
    </html>`
	doc, _ := html.Parse(strings.NewReader(s))

	var nodes = selector.SelectHtmlNodes(doc, "//a[5 mod 2]")
	for _, node := range nodes {
		fmt.Println(selector.HtmlNodeInnerText(node))
	}
}
