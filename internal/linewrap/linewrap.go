package linewrap

import (
	"io"
)

// Wrapper accepts string tokens and outputs them
type Wrapper struct {
	out           io.Writer
	cols          int // column limit for wrapping words
	count         int
	initialPrefix string // prefix for first line
	prefix        string // prefix for subsequent lines
	firstLine     bool
	newLine       bool
}

// New creates a new Wrapper, given an io.Writer and a column wrap limit
func New(w io.Writer, cols int) *Wrapper {
	return NewPrefix(w, cols, "", "")
}

// NewPrefix creates a new Wrapper, which takes an initialPrefix for the first
// line of output, and a prefix for following lines. As well as an io.Writer
// for output and a column wrap limit.
func NewPrefix(writer io.Writer, cols int, initialPrefix, prefix string) *Wrapper {
	return &Wrapper{
		out:           writer,
		cols:          cols,
		initialPrefix: initialPrefix,
		prefix:        prefix,
		firstLine:     true,
		newLine:       true,
	}
}

// NewEmbedded creates a new Wrapper, based off of a parent Wrapper, given
// an initialPrefix, and a prefix for subsequent lines, both of which are
// appended to the prefix of the parent Wrapper. (Useful for recursion.)
func (w *Wrapper) NewEmbedded(initialPrefix, prefix string) *Wrapper {
	return NewPrefix(w.out, w.cols, w.prefix+initialPrefix, w.prefix+prefix)
}

// WriteToken writes a single token to out the output. If this is a newline,
// it will first write the appropriate prefix. Spaces are inserted between
// tokens, but not between the prefix and the first token, or after the last
// token on a line.
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
		w.count += len(w.prefix)
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

// WriteTokens writes an array of string tokens, calling WriteToken for each.
func (w *Wrapper) WriteTokens(tokens []string) {
	for i := range tokens {
		if tokens[i] != "" {
			w.WriteToken(tokens[i])
		}
	}
}

// TerminateLine writes a newline character, unless the current line is empty.
func (w *Wrapper) TerminateLine() {
	if !w.newLine {
		w.Newline()
	}
}

// Newline creates a new line. If the current line is empty, it writes the
// appropriate prefix first.
func (w *Wrapper) Newline() {
	if w.firstLine {
		w.out.Write([]byte(w.initialPrefix))
		w.count += len(w.initialPrefix)
	} else if w.newLine {
		w.out.Write([]byte(w.prefix))
		w.count += len(w.prefix)
	}
	w.out.Write([]byte("\n"))
	w.count = 0
	w.newLine = true
}
