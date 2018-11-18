package main

import (
	"os"
	"flag"
	"fmt"
	"github.com/bobertlo/vmd/vmdfmt/markdown"
)

func main() {
	pretty := flag.Bool("pretty", false, "enable pretty formatting")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: vmdfmt [-pretty] <file.md>")
		os.Exit(1)
	}
	path := flag.Arg(0)
	
	buf, err := markdown.RenderFile(path, *pretty)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(buf)
}
