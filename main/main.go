package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zhengchun/selector"
)

func html() {
	//HTML
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
	doc, _ := selector.ParseHtml(strings.NewReader(s))
	var nodes = selector.SelectHtmlNodes(doc, "//a")
	for _, node := range nodes {
		fmt.Println(selector.HtmlNodeInnerText(node))
	}
}

func xml() {
	//XML
	s := `<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
<book>
  <title lang="en">Harry Potter</title>
  <price>29.99</price>
</book>
<book>
  <title lang="en">Learning XML</title>
  <price>39.95</price>
</book>
</bookstore>`
	doc, _ := selector.ParseXml(strings.NewReader(s))
	var nodes = selector.SelectXmlNodes(doc, "//book")
	for _, node := range nodes {
		fmt.Println(selector.XmlNodeInnerText(node))
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		html()
		wg.Done()
	}()
	go func() {
		xml()
		wg.Done()
	}()
	wg.Wait()
}
