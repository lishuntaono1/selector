Selector
======
A golang package lets you extract data from HTML/XML documents use XPath selectors,inspire by .NET XPath library.

`NOTES: some XPath syntax features not supported yet, so please reading the features instruction those already supported features before use.` 

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
| [ ]            | //book[1]                                |
| last()         | //book[last()]                           | 
| position()     | //book[position()=1]                     |
| [num]          | //book[1]  /book[3]                      |


### Functions & Operations
| Expression     | examples                                 |
| :-------------:| :------------------------                |
| |              | //node | //node                          |
| +              | 6 + 4                                    |
| -              | 6 - 4                                    |
| *              | 6 * 4                                    |
| div            | 8 div 4                                  |
| =              | price=9.80                               |
| !=             | price!=9.80                              |
| <              | price < 9.80                             |
| <=             | price <= 9.80                            |
| >              | price > 9.80                             |
| >=             | price >= 9.80                            |
| or             | price = 9.80 or price = 9.70             |
| and            | price > 9.00 and price < 9.90            |
| mod            | 5 mod 2                                  |

TODO
======
