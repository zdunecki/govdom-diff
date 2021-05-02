package vdom

import "golang.org/x/net/html"

type HTML struct {
	NodeType    html.NodeType     `json:"nT"`
	ID          uint              `json:"i"`
	Name        string            `json:"n"`
	TagName     string            `json:"tN"`
	Attrs       map[string]string `json:"a"`
	TextContent string            `json:"tC"`
	ChildNodes  []*HTML           `json:"cN"`
}
