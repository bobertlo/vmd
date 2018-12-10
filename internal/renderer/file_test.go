package renderer

import (
	"bufio"
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

func TestColumns(t *testing.T) {
	var files [3]string
	files[0] = "testfiles/lorem.md"
	files[1] = "testfiles/lorem-list.md"
	files[2] = "testfiles/lorem-blocks.md"

	for i := range files {
		src, err := ioutil.ReadFile(files[i])
		if err != nil {
			t.Error("could not load test file")
		}

		for c := 30; c < 101; c += 10 {
			r := New(c)

			out, err := r.RenderBytes(src)
			if err != nil {
				t.Error(err)
			}

			scanner := bufio.NewScanner(bytes.NewReader(out))
			for scanner.Scan() {
				if len(scanner.Text()) > c {
					t.Error("line too long")
				}
			}
			err = scanner.Err()
			if err != nil {
				t.Error(err)
			}
		}
	}
}
