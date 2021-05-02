package vdom

import (
	"testing"
)

func TestParser(t *testing.T) {
	tests := NewTestFixtures()

	// TODO: diff compare
	_, err := tests.encodeHTMLFromFile("./encoder1.html")
	if err != nil {
		t.Error(err)
	}
}
