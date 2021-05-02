package vdom

import (
	"testing"
)

func TestNewPatcher(t *testing.T) {
	tests := NewTestFixtures()

	// TODO: diff compare
	_, err := tests.patchFromFile("./diff1a.html", "./diff2a.html")
	if err != nil {
		t.Error(err)
	}
}
