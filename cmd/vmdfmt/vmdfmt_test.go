package main

import (
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
