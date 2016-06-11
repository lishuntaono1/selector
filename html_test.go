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
	//shoudBeEqual(t, "//a[last()]", 3)
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
            <a href="/football">football</a>
            <a href="/basketball">basketball</a>
            <a href="/ping_pong">ping-pong</a>
        </li>
        <li>
            <a href="/swiming">swiming</a>
            <a href="/running">running</a>
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
