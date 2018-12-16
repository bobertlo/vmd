package renderer

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestTrimFlatten(t *testing.T) {
	a := trimFlattenSpaces([]byte(" testing  leading  spaces"))
	b := []byte("testing leading spaces")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("leading spaces failed")
	}

	a = trimFlattenSpaces([]byte("testing trailing   spaces  "))
	b = []byte("testing trailing spaces")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("traling spaces failed")
	}

	a = trimFlattenSpaces([]byte("testing   interior spaces"))
	b = []byte("testing interior spaces")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("leading spaces failed")
	}
}

func TestFlatten(t *testing.T) {
	a := flattenSpaces([]byte(" testing  leading spaces"))
	b := []byte(" testing leading spaces")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("leading spaces failed")
	}

	a = flattenSpaces([]byte("testing trailing   spaces  "))
	b = []byte("testing trailing spaces ")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("traling spaces failed")
	}

	a = flattenSpaces([]byte("testing   interior  spaces"))
	b = []byte("testing interior spaces")
	if bytes.Compare(a, b) != 0 {
		t.Errorf("leading spaces failed")
	}
}

func TestBytes(t *testing.T) {
	r := New(80)
	out, err := r.RenderBytes([]byte("<p>html not supported</p>"))
	if err == nil || out != nil {
		t.Error("bytes failed")
	}
	out, err = r.RenderBytes([]byte(""))
	if err != nil || out != nil {
		t.Error("bytes failed")
	}
	out, err = r.RenderBytes([]byte(" \n \n"))
	if err != nil || out != nil {
		t.Error("bytes failed")
	}
}

func TestRender(t *testing.T) {
	r := New(80)
	src, err := ioutil.ReadFile("testfiles/README.md")
	if err != nil {
		t.Error("read error")
	}

	fout, err := r.RenderFile("testfiles/README.md")
	if err != nil {
		t.Error("RenderFile failed")
	}
	if bytes.Compare(src, fout) != 0 {
		t.Error("RenderFile inconsistent")
	}

	fout, err = r.RenderFile("testfiles/nonexistantfile.go")
	if fout != nil || err == nil {
		t.Error("invalid handling of nonexistant file")
	}
}
