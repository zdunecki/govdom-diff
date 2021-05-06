package vdom

import (
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

// Diff returns PatchList - list of steps to produce newNode from oldNode
func Diff(oldNode, newNode *VNode) PatchList {
	patches := make(PatchList, 0)
	diff(oldNode, newNode, &patches)

	return patches
}

// TODO: bugs in patch remove and replace - indalid nodes are marked

// diff is internal function for Diff
// hides technical details
func diff(oldNode, newNode *VNode, patches *PatchList) {
	if oldNode == nil && newNode == nil {
		return
	}

	if oldNode == nil {
		if newNode.children == nil {
			return
		}

		for _, node := range newNode.children {
			*patches = append(*patches, &Patch{
				Type: PatchInsert,
				Node: node,
			})

			diff(oldNode, node, patches)
		}

		return
	}

	oldNodeChildrenLen := len(oldNode.children)
	newNodeChildrenLen := len(newNode.children)
	minNumNodes := newNodeChildrenLen

	//add child nodes
	if newNodeChildrenLen > oldNodeChildrenLen {
		for _, node := range newNode.children[oldNodeChildrenLen:] {
			*patches = append(*patches, &Patch{
				Type: PatchInsert,
				Node: node,
			})
		}

		minNumNodes = oldNodeChildrenLen
	}

	//delete child nodes
	if oldNodeChildrenLen > newNodeChildrenLen {
		for _, node := range oldNode.children[newNodeChildrenLen:] {
			*patches = append(*patches, &Patch{
				Type: PatchRemove,
				Node: node,
			})
		}

		minNumNodes = newNodeChildrenLen
	}

	// compare text
	if oldNode.Type == html.TextNode && newNode.Type == html.TextNode {
		if oldNode.Data != newNode.Data {
			*patches = append(*patches, &Patch{
				Type: PatchText,
				Node: newNode,
			})
		}

		return
	}

	// same tagName
	if oldNode.Data == newNode.Data {
		propsDiff := cmp.Diff(oldNode.Attr, newNode.Attr)
		if propsDiff != "" {
			*patches = append(*patches, &Patch{
				Type: PatchProps,
				Node: newNode,
			})
		}

		// let's compare childs
		for i := 0; i < minNumNodes; i++ {
			diff(oldNode.children[i], newNode.children[i], patches)
		}

		return
	}

	// tagName changed - replace node
	*patches = append(*patches, &Patch{
		Type: PatchReplace,
		Node: newNode,
	})

	return
}
