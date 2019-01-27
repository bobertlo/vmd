package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsMarkdown(t *testing.T) {
	f, _ := os.Stat("../../README.md")
	if !isMarkdownFile(f) {
		t.Error("isMarkdownFile failed")
	}
	f, _ = os.Stat(".")
	if isMarkdownFile(f) {
		t.Error("isMarkdownFile failed")
	}
	f, _ = os.Stat("main.go")
	if isMarkdownFile(f) {
		t.Error("isMarkdownFile failed")
	}
}

func TestProcessFile(t *testing.T) {
	src, err := ioutil.ReadFile("../../internal/renderer/testfiles/README.md")
	if err != nil {
		t.Error("could not load README.md")
	}

	//out := new(bytes.Buffer)
	var out = bytes.NewBuffer(nil)

	err = processFile("../../internal/renderer/testfiles/README.md", nil, out)
	if err != nil {
		t.Error("processFile failed")
	}

	if bytes.Compare(src, out.Bytes()) != 0 {
		t.Error("file mismatch")
	}
}
