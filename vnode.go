package vdom

import (
	"golang.org/x/net/html"
)

// VNode represents Virtual HTML Node
type VNode struct {
	// parent is reference to 'parent' VNode
	parent *VNode
	// index identifies current VNode using index array indicator
	index []uint64
	// children list of all VNode childrens
	children []*VNode

	// html.NodeType
	Type html.NodeType
	// Data is the same as "Data" in html.Node
	Data string
	// html.Attribute
	Attr []html.Attribute
}

// getNodeAtIndex finds VNode based on index
func (n *VNode) getNodeAtIndex(index []uint64) (foundNode *VNode) {
	if len(index) <= 1 {
		foundNode = n
		return
	}

	foundNode = n

	// TODO: bug in parser
	if foundNode.Type == html.DocumentNode {
		foundNode = foundNode.children[0]
	}

	tailIndex := index[1:]

	for i := range tailIndex {
		foundNode = foundNode.children[tailIndex[i]]
	}

	return foundNode
}

// attrToMap
func (n *VNode) attrToMap() map[string]string {
	out := make(map[string]string)

	for _, attr := range n.Attr {
		out[attr.Key] = attr.Val
	}

	return out
}
