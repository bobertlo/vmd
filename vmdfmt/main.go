package main

import (
	"os"
	"fmt"
	"./markdown"
)

func main() {
	buf, err := markdown.RenderFile("test.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(buf)
}
