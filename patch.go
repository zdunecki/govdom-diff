package vdom

type (
	// Patch holds information needed for produce new dom state
	Patch struct {
		Type PatchType
		Node *VNode
	}

	// PatchList is a list of steps to produce new DOM state
	PatchList []*Patch

	// PatchType identifies patch action
	PatchType uint64
)

const (
	// PatchReplace replaces whole element
	PatchReplace PatchType = iota

	// PatchProps updates element properties
	PatchProps

	// PatchText updates element content
	PatchText

	// PatchInsert inserts new node
	PatchInsert

	// PatchRemove deletes node
	PatchRemove

	PatchReorder
)

type patcher struct {
	node *VNode
	ps   PatchList
}

func NewPatcher(node *VNode, ps PatchList) *patcher {
	return &patcher{node: node, ps: ps}
}

// Patch applies PatchList states to current VNode and produce new one's virtual DOM
func (p *patcher) Patch() *VNode {
	deletedIndexes := make([][]uint64, 0)

	for _, ps := range p.ps {
		switch ps.Type {
		case PatchReplace:
			updateNode := p.node.getNodeAtIndex(ps.Node.index)
			updateNode.Data = ps.Node.Data
			updateNode.Attr = ps.Node.Attr
			updateNode.Type = ps.Node.Type
		case PatchProps:
			updateNode := p.node.getNodeAtIndex(ps.Node.index)
			updateNode.Attr = ps.Node.Attr
		case PatchText:
			updateNode := p.node.getNodeAtIndex(ps.Node.index)
			updateNode.Data = ps.Node.Data
		case PatchInsert:
			patchNodeParentIndex := ps.Node.index[:len(ps.Node.index)-1]
			updateParentNode := p.node.getNodeAtIndex(patchNodeParentIndex)
			updateParentNode.children = append(updateParentNode.children, ps.Node)
		case PatchRemove:
			lastPatchIndex := ps.Node.index[len(ps.Node.index)-1]
			patchNodeParentIndex := ps.Node.index[:len(ps.Node.index)-1]
			updateParentNode := p.node.getNodeAtIndex(patchNodeParentIndex)

			updateParentNode.children[lastPatchIndex] = nil
			deletedIndexes = append(deletedIndexes, ps.Node.index)
		}
	}

	// delete nil nodes after PatchRemove
	// it's needed because we can't remove nodes during PatchList loop
	for _, index := range deletedIndexes {
		deletedNodeParent := p.node.getNodeAtIndex(index[:len(index)-1])

		newNodeParentChilds := make([]*VNode, 0)
		for _, parentChild := range deletedNodeParent.children {
			if parentChild != nil {
				newNodeParentChilds = append(newNodeParentChilds, parentChild)
			}
		}

		deletedNodeParent.children = newNodeParentChilds
	}

	return p.node
}
