package main

import (
	"flag"
	"fmt"
	"github.com/bobertlo/vmd/internal/markdown"
	"os"
)

func main() {
	pretty := flag.Bool("pretty", false, "enable pretty formatting")
	cols := flag.Int("cols", 80, "number of columns to wrap output")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: vmdfmt [-pretty] <file.md>")
		os.Exit(1)
	}
	path := flag.Arg(0)

	r := markdown.NewRenderer(*cols, *pretty)
	buf, err := r.RenderFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(buf)
}
