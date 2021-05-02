package vdom

import (
	"testing"
)

func TestDiff(t *testing.T) {
	tests := NewTestFixtures()

	// TODO: diff compare
	_, err := tests.diffFromFiles("./diff1.html", "./diff2.html")
	if err != nil {
		t.Error(err)
	}
}
