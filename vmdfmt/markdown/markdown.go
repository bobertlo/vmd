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
	pretty bool
	width int
	prefix string
	indent int
}

// flattenSpaces removes all leading, trailing, and reduntant spaces from a
// []byte array, leaving internal single spaces.
func flattenSpaces(str []byte) []byte {
	re := regexp.MustCompile("  +")
	replaced := re.ReplaceAll(bytes.TrimSpace(str), []byte(" "))
	replaced = bytes.TrimSpace(replaced)
	return replaced
}

// LoadMarkdown reads a file into a []byte buffer and parses it into a 
// blackfriday markdown tree.
func LoadMarkdown(path string) (*blackfriday.Node, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := blackfriday.New(blackfriday.WithExtensions(
		blackfriday.Tables|blackfriday.FencedCode))
	n := m.Parse(dat)

	return n, nil
}

// FileRenderer parses a markdown tree from a file and creates a new Renderer
func NewRenderer(pretty bool) *Renderer {
	buf := new(bytes.Buffer)
	r := &Renderer{
		out: buf,
		pretty: false,
		width: 80,
		indent: 0,
		prefix: "",
	}
	return r
}

// RenderFile renders a markdown file to the out buffer, returning a formatted
// ([]byte,nil) or (nil,err) if an error occurs
func RenderFile(path string, pretty bool) ([]byte, error) {
	r := NewRenderer(pretty)

	n, err := LoadMarkdown(path)
	if err != nil {
		return nil, err
	}

	err = r.Render(n)
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

func (r *Renderer) Render(root *blackfriday.Node) error {
	// if passed a full document, start on the first child node
	if root.Type == blackfriday.Document {
		root = root.FirstChild
	}

	for c := root; c != nil; c = c.Next {
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


