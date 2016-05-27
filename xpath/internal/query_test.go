package internal

import "testing"

func TestQueryBuilder(t *testing.T) {
	//s := `<p class="a">Links:<br/></p><ul><li><a href="foo">Foo</a></li><li><a href="/bar/baz">BarBaz</a></li></ul>`
	//doc, _ := html.Parse(strings.NewReader(s))

	root := ParseXPathExpression("//a")
	builder := &QueryBuilder{
		query: "//a",
	}
	props := noneProp
	builder.processNode(root, noneFlag, &props)
}
