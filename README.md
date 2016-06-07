Selector
======
A golang package lets you extract data from HTML/XML documents use XPath selectors,inspire by .NET XPath library.

`Notes: some XPath syntax features not supported yet, so please reading the features instruction those already supported features before use.` 

*[XPath Tutorial](http://www.w3schools.com/xsl/xpath_syntax.asp)*

*[XPath Syntax MSDN](https://msdn.microsoft.com/en-us/library/ms256471(v=vs.110).aspx)

Features
======
### Basic Path Expression

| Expression     | examples                                 |
| :-------------:| :------------------------                |
| nodename       | author bookstore                         |
| //             | //author //bookstore bookstore//book  	|
| /              | /author bookstore/book   	            |
| .              | . //a/. .//title                         |
| ..             | .., //a/..                               |
| *              | * //author/* */*                         |
| @              | //@lang title[@*]  //title[@lang]        |

### Predicates

| Expression     | examples                                 |
| :-------------:| :------------------------                |
| @              | //title[@lang='en']  //*[@src]           |
| []             | //book[1]                                |
| last()         | // 

### Functions & Operations


TODO
======
>>>

Examples
======
```golang
package main
import (
	"fmt"
	"strings"
	"github.com/zhengchun/selector"
	"golang.org/x/net/html"
)
func main() {
	s := `<ul><li><a href="foo">Foo</a><a href="ball">Ball</a></li><li><a href="bar">BarBaz</a></li></ul>`
	doc, _ := html.Parse(strings.NewReader(s))
	var nodes = selector.SelectHtmlNodes(doc, "//a")	
	for _, node := range nodes {
		fmt.Println(InnerText(node))
	}
}
```

// functions

//last
//name
//count

`|`

parent::*
ancestor::*
descendant-or-self::*

// mod - + div