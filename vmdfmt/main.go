package main

import (
	"os"
	"./markdown"
)

func main() {
	buf, err := markdown.RenderFile("test.md")
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(buf)
}
