package vdom

import (
	"golang.org/x/net/html"
)

type decoder struct {
	node *VNode
}

func NewDecoder(node *VNode) *decoder {
	return &decoder{node: node}
}

func (d *decoder) ToJSON() (*HTML, error) {
	var id uint = 0
	var walk func(*HTML, *VNode)

	walk = func(tree *HTML, node *VNode) {
		if tree.ChildNodes == nil {
			tree.ChildNodes = make([]*HTML, 0)
		}

		for _, c := range node.children {
			attrs := c.attrToMap()
			htmlNode := &HTML{
				NodeType: c.Type,
				ID:       id,
				Name:     attrs["name"],
				Attrs:    attrs,
			}

			switch c.Type {
			case html.TextNode:
				htmlNode.TextContent = c.Data
			case html.DoctypeNode, html.DocumentNode, html.ElementNode:
				htmlNode.TagName = c.Data
			case html.CommentNode:
				htmlNode.TextContent = c.Data
			}

			tree.ChildNodes = append(tree.ChildNodes, htmlNode)

			walk(htmlNode, c)
			id++
		}
	}

	treeAttrs := d.node.attrToMap()
	htmlTree := &HTML{
		NodeType:    d.node.Type,
		ID:          id,
		Name:        treeAttrs["name"],
		TagName:     d.node.Data,
		Attrs:       treeAttrs,
		TextContent: "",
		ChildNodes:  nil,
	}

	walk(htmlTree, d.node)

	return htmlTree, nil
}
