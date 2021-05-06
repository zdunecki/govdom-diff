package vdom

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var flagURL = flag.String("url", "https://livesession.io", "Name of URL to HTML")

func TestJsonToHTML(t *testing.T) {
	u, err := url.Parse(*flagURL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}

	respBodyCopy, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyCopy)) // reset response body to unread state

	node, err := NewHTMLReader(resp.Body).ToJSON()

	jsonEncoder := NewJSONEncoder(node)

	nodeHTML, err := jsonEncoder.EncodeToHTML()
	if err != nil {
		t.Fatal(err)
	}

	// TODO: compare respBodyCopy and htmlTreeBuffer
	// htmlTreeBuffer has different order than respBodyCopy and some attributes are lower instead of camel case

	if nodeHTML == nil {
		t.Fatal("htmlTreeBuffer should not be nil")
	}
}

func TestParser(t *testing.T) {
	tests := NewTestFixtures()

	// TODO: diff compare
	_, err := tests.encodeHTMLFromFile("./encoder1.html")
	if err != nil {
		t.Error(err)
	}
}

func TestNewDecoder(t *testing.T) {
	tests := NewTestFixtures()

	vnodeB, err := tests.vNodeFromFileToJSON("./diff2.html")
	expectedVnodeB, err := tests.loadFile("./expected/diff2.json")
	if err != nil {
		t.Error(err)
	}

	tests.diff(t, "compare vnode json", expectedVnodeB, vnodeB)
}
