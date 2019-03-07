package renderer

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path"
	"testing"
)

const testPath = "testfiles/"

var verbatimFiles = []string{
	"README.md",
	"blockquote-codeblock.md",
	"blockquote-issue7.md",
 }

var columnFiles = []string{"lorem.md", "lorem-list.md", "lorem-blocks.md"}

func TestFiles(t *testing.T) {
	for _, f := range verbatimFiles {
		t.Run(f, func(t *testing.T) {
			src, err := ioutil.ReadFile(path.Join(testPath, f))
			if err != nil {
				t.Error("could not load")
			}

			r := New(80)
			out, err := r.RenderBytes(src)
			if err != nil {
				t.Error("Failed to render")
			}

			if bytes.Compare(src, out) != 0 {
				t.Error("Inconsistency in render")
			}
		})
	}
}

func TestColumns(t *testing.T) {
	for _, f := range columnFiles {
		t.Run(f, func(t *testing.T) {
			src, err := ioutil.ReadFile(path.Join(testPath, f))
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
		})
	}
}
