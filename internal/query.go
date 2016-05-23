package internal

import "golang.org/x/net/html"

type Query interface {
	Select(*html.Node) bool
	Evaluate(*html.Node) interface{}
}
