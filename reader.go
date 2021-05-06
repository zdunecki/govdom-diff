package vdom

import (
	"io"

	"golang.org/x/net/html"
)

const (
	HTMLTextarea = "TEXTAREA"
	HTMLInput    = "INPUT"
	HTMLSelect   = "SELECT"
	HTMLScript   = "SCRIPT"
	HTMLNoScript = "NOSCRIPT"
)

const (
	HTMLAttributeValue = "value"
)

type htmlReader struct {
	*HTML
	reader io.Reader
}

func NewHTMLReader(reader io.Reader) *htmlReader {
	return &htmlReader{
		reader: reader,
	}
}

func (h *htmlReader) ToJSON() (*HTML, error) {
	htmlNode, err := html.Parse(h.reader)
	if err != nil {
		return nil, err
	}

	var walk func(htmlNode *html.Node, vNode *HTML, id uint)

	attrToMap := func(htmlNode *html.Node) map[string]string {
		out := make(map[string]string)

		if htmlNode != nil {
			for _, attr := range htmlNode.Attr {
				out[attr.Key] = attr.Val
			}
		}
		return out
	}

	walk = func(htmlNode *html.Node, node *HTML, id uint) {
		if htmlNode.FirstChild != nil {
			node.ChildNodes = make([]*HTML, 0)
		}

		for c, i := htmlNode.FirstChild, 0; c != nil; c, i = c.NextSibling, i+1 {
			attrs := attrToMap(c)
			newNode := &HTML{
				NodeType: c.Type,
				ID:       id,
				Name:     attrs["name"],
				Attrs:    attrs,
			}

			switch c.Type {
			case html.TextNode:
				newNode.TextContent = c.Data
			case html.DoctypeNode, html.DocumentNode, html.ElementNode:
				newNode.TagName = c.Data
			case html.CommentNode:
				newNode.TextContent = c.Data
			}

			node.ChildNodes = append(node.ChildNodes, newNode)

			id++

			walk(c, newNode, id)
		}
	}

	treeAttrs := attrToMap(htmlNode)
	node := &HTML{
		NodeType:    htmlNode.Type,
		ID:          0,
		Name:        treeAttrs["name"],
		TagName:     htmlNode.Data,
		Attrs:       treeAttrs,
		TextContent: "",
		ChildNodes:  nil,
	}

	walk(htmlNode, node, node.ID)

	return node, nil
}
