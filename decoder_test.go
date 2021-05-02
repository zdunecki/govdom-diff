package vdom

import (
	"testing"
)

func TestNewDecoder(t *testing.T) {
	tests := NewTestFixtures()

	vnodeB, err := tests.vNodeFromFileToJSON("./diff2.html")
	expectedVnodeB, err := tests.loadFile("./expected/diff2.json")
	if err != nil {
		t.Error(err)
	}

	tests.diff(t, "compare vnode json", expectedVnodeB, vnodeB)
}
