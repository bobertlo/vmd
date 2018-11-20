package linewrap

import (
	"io"
)

type Wrapper struct {
	out           io.Writer
	cols          int // column limit for wrapping words
	count         int
	initialPrefix string // prefix for first line
	prefix        string // prefix for subsequent lines
	firstLine     bool
	newLine       bool
}

func New(w io.Writer, cols int) *Wrapper {
	return NewPrefix(w, cols, "", "")
}

func NewPrefix(writer io.Writer, cols int, initialPrefix, prefix string) *Wrapper {
	return &Wrapper{
		out:           writer,
		cols:          cols,
		initialPrefix: initialPrefix,
		prefix:        prefix,
		firstLine:     true,
		newLine:       true,
	}
	return nil
}

func (w *Wrapper) NewEmbedded(initialPrefix, prefix string) *Wrapper {
	return NewPrefix(w.out, w.cols, w.prefix+initialPrefix, w.prefix+prefix)
}

func (w *Wrapper) WriteToken(token string) {
	if w.firstLine {
		w.out.Write([]byte(w.initialPrefix))
		w.out.Write([]byte(token))
		w.firstLine = false
		w.newLine = false
		w.count += len(w.initialPrefix)
		w.count += len(token)
		if w.count > w.cols {
			w.out.Write([]byte("\n"))
			w.count = 0
			w.newLine = true
		}
	} else if w.newLine {
		w.out.Write([]byte(w.prefix))
		w.out.Write([]byte(token))
		w.newLine = false
		w.count += len(w.initialPrefix)
		w.count += len(token)
		if w.count > w.cols {
			w.out.Write([]byte("\n"))
			w.count = 0
			w.newLine = true
		}
	} else {
		// if the token is too long for this token, create a newline
		// and recurse (to handle prefixes)
		if w.count+len(token)+1 > w.cols {
			w.out.Write([]byte("\n"))
			w.count = 0
			w.newLine = true
			w.WriteToken(token)
		} else {
			w.out.Write([]byte(" "))
			w.out.Write([]byte(token))
			w.count += len(token) + 1
			if w.count > w.cols {
				w.out.Write([]byte("\n"))
				w.count = 0
				w.newLine = true
			}
		}
	}
}

func (w *Wrapper) WriteTokens(tokens []string) {
	for i := range tokens {
		if tokens[i] != "" {
			w.WriteToken(tokens[i])
		}
	}
}

func (w *Wrapper) TerminateLine() {
	if !w.newLine {
		w.Newline()
	}
}

func (w *Wrapper) Newline() {
	w.out.Write([]byte("\n"))
	w.count = 0
	w.newLine = true
}
