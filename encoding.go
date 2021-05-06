package vdom

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type jsonencoder struct {
	node *HTML
}

// NewJSONEncoder encoding *HTML to different formats
func NewJSONEncoder(node *HTML) *jsonencoder {
	return &jsonencoder{node: node}
}

// EncodeToHTML encode *HTML to "HTML"
func (encoder *jsonencoder) EncodeToHTML() (io.Reader, error) {
	var walk func(htmlTree, prevSibling, nextSibling *html.Node, node *HTML)

	walk = func(htmlTree, prevSibling, nextSibling *html.Node, node *HTML) {
		htmlTree.Type = node.NodeType
		htmlTree.Data = treeDataByNodeType(node)
		htmlTree.Attr = htmlAttrs(node)

		htmlTree.PrevSibling = prevSibling
		htmlTree.NextSibling = nextSibling

		if node.ChildNodes == nil { // TODO:
			return
		}

		treeChildNodes := make([]*html.Node, len(node.ChildNodes))

		for i, _ := range treeChildNodes {
			treeChildNodes[i] = &html.Node{}
		}

		for i, child := range node.ChildNodes {
			var prev, next *html.Node

			childNodesC := len(node.ChildNodes)
			firstIter := i == 0
			lastIter := i == childNodesC-1

			currentTreeChild := treeChildNodes[i]

			if firstIter {
				htmlTree.FirstChild = currentTreeChild

				if childNodesC > 1 {
					next = treeChildNodes[i+1]
				}
			} else if lastIter {
				htmlTree.LastChild = currentTreeChild

				if childNodesC > 1 {
					prev = treeChildNodes[i-1]
				}
			} else {
				prev, next = treeChildNodes[i-1], treeChildNodes[i+1]
			}

			currentTreeChild.Parent = htmlTree
			if prev != nil {
				prev.Parent = htmlTree
			}
			if next != nil {
				next.Parent = htmlTree
			}

			walk(currentTreeChild, prev, next, child)
		}
	}

	htmlTree := &html.Node{}

	walk(htmlTree, nil, nil, encoder.node)

	b := new(bytes.Buffer)
	if err := html.Render(b, htmlTree); err != nil {
		return nil, err
	}

	return b, nil
}

type walkHTMLNode func(*html.Node, *VNode, uint64)

type htmlencoder struct {
	node *html.Node
}

// NewHTMLEncoder encoding *html.Node to different formats
func NewHTMLEncoder(node *html.Node) *htmlencoder {
	return &htmlencoder{node: node}
}

// TODO: bug - first index is html not document
// TODO: write own tokenizer to avoid double traverse tree
// TODO: bug in html.Node - sometimes same pages has different parsed structure

// EncodeToVNode encode *html.Node to *VNode
func (encoder *htmlencoder) EncodeToVNode() (*VNode, error) {
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

	walk(encoder.node, vNode, 0)

	return vNode, nil
}

type vnodeencoder struct {
	node *VNode
}

// NewVNodeEncoder encoding *VNode to different formats
func NewVNodeEncoder(node *VNode) *vnodeencoder {
	return &vnodeencoder{node: node}
}

// EncodeToJSON encode *VNode to *HTML
func (encoder *vnodeencoder) EncodeToJSON() (*HTML, error) {
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

	treeAttrs := encoder.node.attrToMap()
	htmlTree := &HTML{
		NodeType:    encoder.node.Type,
		ID:          id,
		Name:        treeAttrs["name"],
		TagName:     encoder.node.Data,
		Attrs:       treeAttrs,
		TextContent: "",
		ChildNodes:  nil,
	}

	walk(htmlTree, encoder.node)

	return htmlTree, nil
}

func htmlAttrs(node *HTML) []html.Attribute {
	if node.Attrs == nil {
		return nil
	}

	attrs := make([]html.Attribute, 0)
	for key, val := range node.Attrs {
		attrs = append(attrs, html.Attribute{
			Namespace: "",
			Key:       key,
			Val:       val,
		})
	}

	return attrs
}

func treeDataByNodeType(node *HTML) string {
	switch node.NodeType {
	case html.ErrorNode:
		return ""
	case html.TextNode:
		return node.TextContent
	case html.ElementNode:
		return node.TagName
	case html.CommentNode:
		return node.TextContent
	case html.DoctypeNode:
		return node.TagName
	case html.RawNode:
		return node.TextContent
	default:
		return ""
	}
}
