package vdom

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetNodeAtIndex(t *testing.T) {
	tests := NewTestFixtures()

	text1 := "Hello"
	text2 := "World"
	expect := text1 + " " + text2

	testNode := &VNode{
		index: []uint64{0},
		children: []*VNode{
			{
				index: []uint64{0, 0},
				Type:  html.ElementNode,
				Data:  "div",
				Attr:  nil,
				children: []*VNode{
					{
						index: []uint64{0, 0, 0},
						children: []*VNode{
							{
								index:    []uint64{0, 0, 0, 0},
								children: nil,
								Type:     html.TextNode,
								Data:     text1,
								Attr:     nil,
							},
						},
						Type: html.ElementNode,
						Data: "h1",
						Attr: nil,
					},
					{
						index: []uint64{0, 0, 1},
						children: []*VNode{
							{
								index:    []uint64{0, 0, 1, 0},
								children: nil,
								Type:     html.TextNode,
								Data:     text2,
								Attr:     nil,
							},
						},
						Type: html.ElementNode,
						Data: "p",
						Attr: nil,
					},
				},
			},
			{
				index:    []uint64{0, 1},
				Type:     html.ElementNode,
				Data:     "div",
				Attr:     nil,
				children: nil,
			},
		},
		Type: html.ElementNode,
		Data: "body",
		Attr: nil,
	}

	head := testNode.getNodeAtIndex([]uint64{0})

	helloNode := testNode.getNodeAtIndex([]uint64{0, 0, 0, 0})
	worldNode := testNode.getNodeAtIndex([]uint64{0, 0, 1, 0})

	textToTest := helloNode.Data + " " + worldNode.Data

	tests.diff(t, "should find text nodes and match with expected text", expect, textToTest)
	tests.diff(t, "should find first head", []uint64{0}, head.index)
}
