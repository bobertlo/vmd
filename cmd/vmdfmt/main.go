package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobertlo/vmd/pkg/mdformatter"
)

var (
	cols  = flag.Int("cols", 80, "number of columns to wrap output")
	write = flag.Bool("w", false, "write changes to (source) file")
	list  = flag.Bool("l", false, "list files with modifications")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: vmdfmt [flags] [path ...]")
	flag.PrintDefaults()
}

func processFile(path string, in io.Reader, out io.Writer) error {
	var perm os.FileMode = 0644
	if in == nil {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		in = f
		perm = fi.Mode().Perm()
	}

	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	md := mdformatter.New(*cols)
	output, err := md.RenderBytes(input)
	if err != nil {
		return err
	}

	if !bytes.Equal(output, input) {
		if *list {
			fmt.Fprintln(out, path)
		}
		if *write {
			err = ioutil.WriteFile(path, output, perm)
			if err != nil {
				return err
			}
		}
	}

	if !*write && !*list {
		out.Write(output)
	}

	return nil
}

func isMarkdownFile(f os.FileInfo) bool {
	if f.IsDir() {
		return false
	}
	return strings.HasSuffix(f.Name(), ".md")
}

func walkFile(path string, f os.FileInfo, err error) error {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if isMarkdownFile(f) {
		err := processFile(path, nil, os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		if *write {
			fmt.Fprintln(os.Stderr, "error: cannot use -w when reading stdin")
			os.Exit(1)
		}
		err := processFile("<stdin>", os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}

	for _, f := range flag.Args() {
		dir, err := os.Stat(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if dir.IsDir() {
			filepath.Walk(f, walkFile)
		} else {
			err := processFile(f, nil, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		}
	}
}
