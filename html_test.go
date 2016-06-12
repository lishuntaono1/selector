package selector

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var doc *html.Node

func TestXPathPathExpressions(t *testing.T) {
	// `//a`
	shoudBeEqual(t, "//a", 6)
	// `//li/a`
	shoudBeEqual(t, "//li/a", 5)
	// `//a/@href`
	shoudBeEqual(t, "//a/@href", 6)
	// `//a/..`
	shoudBeEqual(t, "//a/..", 3)
	// `//a/.`
	shoudBeEqual(t, "//a/.", 6)
	// `/html/body/a`
	shoudBeEqual(t, "/html/body/a", 1)
	// `//@href`
	shoudBeEqual(t, "//@href", 6)
}

func TestXPathPredicates(t *testing.T) {
	// `[@src]`
	shoudBeEqual(t, "//*[@src]", 1)
	// `//a[1]`
	shoudBeEqual(t, "//a[1]", 3)
	// `//a[last()]`
	shoudBeEqual(t, "//a[last()]", 3)
	// `//a[position()]`
	shoudBeEqual(t, "//a[position()=3]", 1)
	// `//class='ball`
	shoudBeEqual(t, "//*[@class='ball']", 1)

}

func TestXPathFunctions(t *testing.T) {
	// `|`
	shoudBeEqual(t, "//head | //body", 2)
	// `+`
	shoudBeEqual(t, "//a[1+2]", 1)
	// `div`
	shoudBeEqual(t, "//a[6 div 3]", 2)
	// `!=`
	shoudBeEqual(t, "//a[@href!='/running']", 5)
	// `<=`
	shoudBeEqual(t, "//a[@id<=3]", 3)
	// `>`
	shoudBeEqual(t, "//a[@id>3]", 2)
	// `or`
	shoudBeEqual(t, "//a[@id=1 or @id=3]", 2)
	// `and`
	shoudBeEqual(t, "//a[@id=1 and 1=1]", 1)
}

func shoudBeEqual(t *testing.T, xpath string, expected int) {
	var nodes []*html.Node
	nodes = SelectHtmlNodes(doc, xpath)
	if num := len(nodes); num != expected {
		t.Fatalf("expected is %d but got %d. %s", expected, num, xpath)
	}
}

func init() {
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
    </html>   
    `
	root, err := html.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	doc = root
}
