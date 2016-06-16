package selector

import (
	"strings"
	"testing"
)

var xmldoc *XmlNode

func init() {
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
	doc, err := ParseXml(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	xmldoc = doc
}

func shoudBeEqual2(t *testing.T, xpath string, expected int) {
	var nodes []*XmlNode
	nodes = SelectXmlNodes(xmldoc, xpath)
	if num := len(nodes); num != expected {
		t.Fatalf("expected is %d but got %d. %s", expected, num, xpath)
	}
}

func TestXmlPath(t *testing.T) {
	// `//book`
	shoudBeEqual2(t, "//book", 2)
	// `*[@lang='en']`
	shoudBeEqual2(t, "//*[@lang='en']", 2)
}
