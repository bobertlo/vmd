package renderer

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestReadme(t *testing.T) {
	src, err := ioutil.ReadFile("testfiles/README.md")
	if err != nil {
		t.Error("could not load README.md")
	}

	r := New(80)

	out, err := r.RenderBytes(src)
	if err != nil {
		t.Error("Failed to render README.md")
	}

	if bytes.Compare(src, out) != 0 {
		t.Error("Inconsistency in rendered README.md")
	}
}
