package internal

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	var node = ParseXPathExpression("//a[@href]")
	processNode(node)
}

func processAxis(root *Axis) {
	if root.input != nil {
		processNode(root.input)
	}

	switch root.axistype {
	case AxisChild:
		fmt.Println("AxisChild")
	case AxisAncestor:
		fmt.Println("AxisAncestor")
	case AxisAncestorOrSelf:
		fmt.Println("AxisAncestorOrSelf")
	case AxisAttribute:
		fmt.Println("AxisAttribute")
	case AxisDescendant:
		fmt.Println("AxisDescendant")
	case AxisDescendantOrSelf:
		fmt.Println("AxisDescendantOrSelf")
	case AxisFollowing:
		fmt.Println("AxisFollowing")
	case AxisFollowingSibling:
		fmt.Println("AxisFollowingSibling")
	case AxisNamespace:
		fmt.Println("AxisNamespace")
	case AxisParent:
		fmt.Println("AxisParent")
	case AxisPreceding:
		fmt.Println("AxisPreceding")
	case AxisPrecedingSibling:
		fmt.Println("AxisPrecedingSibling")
	case AxisSelf:
		fmt.Println("AxisSelf")
	}
}

func processFilter(root *Filter) {
	fmt.Println("====Filter(Begin)========")
	processNode(root.condition)

	if root.input != nil {
		processNode(root.input)
	}
	fmt.Println("======Filter(END)========")
}
func processNode(root AstNode) {
	switch root.Type() {
	case AxisAst:
		processAxis(root.(*Axis))
	case FilterAst:
		processFilter(root.(*Filter))
	}
}
