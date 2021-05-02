package vdom

import (
	"golang.org/x/net/html"
)

type walkHTMLNode func(*html.Node, *VNode, uint64)

type encoder struct{}

// NewEncoder struct for encoding different input formats to VNode
func NewEncoder() *encoder {
	return &encoder{}
}

// TODO: bug - first index is html not document
// TODO: write own tokenizer to avoid double traverse tree

// EncodeHTML encode html.Node to VNode
func (e *encoder) EncodeHTML(htmlNode *html.Node) (*VNode, error) {
	var walk walkHTMLNode

	walk = func(htmlNode *html.Node, vNode *VNode, id uint64) {
		if vNode.parent != nil {
			vNode.index = append(vNode.parent.index, id)
		}

		//set properties for new node
		vNode.Type = htmlNode.Type
		vNode.Data = htmlNode.Data
		vNode.Attr = htmlNode.Attr

		if htmlNode.FirstChild != nil {
			vNode.children = make([]*VNode, 0)
		}

		for c, i := htmlNode.FirstChild, 0; c != nil; c, i = c.NextSibling, i+1 {
			vNode.children = append(vNode.children, &VNode{
				parent: vNode,
			})

			currentChild := vNode.children[i]
			walk(c, currentChild, uint64(i))
		}
	}

	vNode := &VNode{}

	walk(htmlNode, vNode, 0)

	return vNode, nil
}
