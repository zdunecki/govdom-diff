package vdom

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

const (
	DataPath = "./data/tests"
)

type fixtures struct{}

func NewTestFixtures() *fixtures {
	return &fixtures{}
}

func (f *fixtures) diff(t *testing.T, title string, vExpect, vCurrent interface{}) {
	if diff := cmp.Diff(vExpect, vCurrent); diff != "" {
		t.Errorf("%s mismatch (-want +got):\n%s", title, diff)
	}
}

func (f *fixtures) loadVNodeFromFile(fileName string) (*VNode, error) {
	diffF, err := f.loadFile(fileName)
	if err != nil {
		return nil, err
	}

	diffHTML, err := html.Parse(bytes.NewReader(diffF))
	if err != nil {
		return nil, err
	}

	vnode, err := NewHTMLEncoder(diffHTML).EncodeToVNode()
	if err != nil {
		return nil, err
	}

	return vnode, nil
}

func (f *fixtures) vNodeFromFileToJSON(fileName string) ([]byte, error) {
	vnode, err := f.loadVNodeFromFile(fileName)
	if err != nil {
		return nil, nil
	}

	d := NewVNodeEncoder(vnode)

	htmlJSON, err := d.EncodeToJSON()
	if err != nil {
		return nil, nil
	}

	return json.Marshal(htmlJSON)
}

func (f *fixtures) diffFromFiles(fileName, fileName2 string) (PatchList, error) {
	vnode, err := f.loadVNodeFromFile(fileName)
	if err != nil {
		return nil, err
	}

	vnode2, err := f.loadVNodeFromFile(fileName2)
	if err != nil {
		return nil, err
	}

	return Diff(vnode, vnode2), nil
}

func (f *fixtures) encodeHTMLFromFile(fileName string) (*VNode, error) {
	parserF, err := f.loadFile(fileName)
	if err != nil {
		return nil, err
	}

	htmlNode, err := html.Parse(bytes.NewReader(parserF))
	if err != nil {
		return nil, err
	}

	return NewHTMLEncoder(htmlNode).EncodeToVNode()
}

func (f *fixtures) patchFromFile(fileName, fileName2 string) (*VNode, error) {
	vnode, err := f.loadVNodeFromFile(fileName)
	if err != nil {
		return nil, nil
	}
	ps, err := f.diffFromFiles(fileName, fileName2)
	if err != nil {
		return nil, nil
	}

	p := NewPatcher(vnode, ps)

	return p.Patch(), nil
}

func (f *fixtures) loadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(DataPath, fileName))
}
