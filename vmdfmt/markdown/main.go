package markdown

import (
	"errors"
	"bytes"
	"io/ioutil"
	"regexp"
	"gopkg.in/russross/blackfriday.v2"
)

type Renderer struct {
	out *bytes.Buffer
	root *blackfriday.Node
	pretty bool
	width int
	indent int
	indentIndex int
}

// flattenSpaces removes all leading, trailing, and reduntant spaces from a
// []byte array, leaving internal single spaces.
func flattenSpaces(str []byte) []byte {
	re := regexp.MustCompile("  +")
	replaced := re.ReplaceAll(bytes.TrimSpace(str), []byte(" "))
	replaced = bytes.TrimSpace(replaced)
	return replaced
}

// loadFile reads a file into a []byte buffer and parses it into a blackfriday
// markdown tree
func loadFile(name string) (*blackfriday.Node, error) {
	dat, err := ioutil.ReadFile("test.md")
	if err != nil {
		return nil, err
	}
	m := blackfriday.New(blackfriday.WithExtensions(
		blackfriday.Tables|blackfriday.FencedCode))
	n := m.Parse(dat)

	return n, nil
}

// FileRenderer parses a markdown tree from a file and creates a new Renderer
func FileRenderer(path string) (*Renderer, error) {
	n, err := loadFile(path)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	r := &Renderer{
		out: buf,
		root: n,
		pretty: false,
		width: 80,
		indent: 0,
		indentIndex: 0,
	}
	return r, nil
}

// RenderFile renders a markdown file to the out buffer, returning a formatted
// ([]byte,nil) or (nil,err) if an error occurs
func RenderFile(path string) ([]byte, error) {
	r, err := FileRenderer(path)
	if err != nil {
		return nil, err
	}
	err = r.Render()
	if err != nil {
		return nil, err
	} 
	return r.out.Bytes(), nil
}

// writes 'c' n times
func (r *Renderer) writeNBytes (n int, c byte) {
	for i := 0; i < n; i++ {
		r.out.WriteByte(c)
	}
}

// Render starts at the document root node and renders every valid child, or
// returns an error if invalid nodes are found anywhere in the tree.
func (r *Renderer) Render() error {
	for c := r.root.FirstChild; c != nil; c = c.Next {
		switch (c.Type) {
		case blackfriday.Heading:
			err := r.heading(c)
			if err != nil {
				return err
			}
		case blackfriday.Paragraph:
			err := r.paragraph(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// headingText checks that n and siblings are text nodes (there shouldn't
// be any siblings) and outputs all the text with whitespace flattened, or
// returns an error if an invalid (non Text) node is found
func (r *Renderer) headingText(n *blackfriday.Node) error {
	for p := n; p != nil; p = n.Next {
		if p.Type != blackfriday.Text {
			return errors.New("Headings may only contain text elements")
		}
		r.out.Write(flattenSpaces(p.Literal))
	}
	return nil
}

// heading outputs a heading node (verified before calling) as an atx-heading
// (e.i '#' for each heading level) followed by the contents of each of it's
// text node children (which should be only one) with whitespace flattened
// or returns an error if an invalid (non Text) node is found. Headings are
// line based and cannot be wrapped, so the output is a raw line.
func (r *Renderer) heading(n *blackfriday.Node) error {
	level := n.HeadingData.Level
	r.writeNBytes(level, '#')
	r.out.WriteByte(' ')
	err := r.headingText(n.FirstChild)
	if err != nil {
		return err
	}
	r.out.WriteString("\n\n")
	return nil
}

func (r *Renderer) paragraph(n *blackfriday.Node) error {
	return nil
}


