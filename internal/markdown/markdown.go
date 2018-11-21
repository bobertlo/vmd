package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/bobertlo/vmd/internal/linewrap"
	"gopkg.in/bobertlo/blackfriday.v2"
)

// Renderer renders blackfriday markdown trees into []byte output
type Renderer struct {
	out    *bytes.Buffer
	pretty bool
	cols   int
}

// flattenSpaces removes all reduntant spaces from a []byte array, leaving
// single spaces
func flattenSpaces(str []byte) []byte {
	re := regexp.MustCompile("  +")
	return re.ReplaceAll(str, []byte(" "))
}

func trimFlattenSpaces(str []byte) []byte {
	return bytes.TrimSpace(flattenSpaces(str))
}

// LoadMarkdown reads a file and parses it into a blackfriday markdown tree,
// returning the document root (*Node, nil) or (nil, err)
func LoadMarkdown(path string) (*blackfriday.Node, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseMarkdown(dat)
}

// ParseMarkdown parses a markdown document from a []byte and returns
// a blackfriday markdown root (*Node, nil) or (nil, err)
func ParseMarkdown(dat []byte) (*blackfriday.Node, error) {
	m := blackfriday.New(blackfriday.WithExtensions(
		blackfriday.Tables | blackfriday.FencedCode |
			blackfriday.NoIntraEmphasis))
	n := m.Parse(dat)

	return n, nil
}

// NewRenderer creates a new markdown renderer. cols specifies how many
// columns to wrap lines at, and pretty specifies whether to format tables
// with whitespace.
func NewRenderer(cols int, pretty bool) *Renderer {
	buf := new(bytes.Buffer)
	r := &Renderer{
		out:    buf,
		pretty: pretty,
		cols:   80,
	}
	return r
}

// RenderFile renders a markdown file to the out buffer, returning a formatted
// ([]byte,nil) or (nil,err) if an error occurs
func (r *Renderer) RenderFile(path string) ([]byte, error) {
	n, err := LoadMarkdown(path)
	if err != nil {
		return nil, err
	}

	return r.Render(n)
}

// RenderBytes parses a markdown document in a []byte and renders it,
// returning a formatted document in a []byte. Returns ([]byte,nil) or
// (nil,err)
func (r *Renderer) RenderBytes(dat []byte) ([]byte, error) {
	n, err := ParseMarkdown(dat)
	if err != nil {
		return nil, err
	}

	return r.Render(n)
}

// writes 'c' n times
func (r *Renderer) writeNBytes(n int, c byte) {
	for i := 0; i < n; i++ {
		r.out.WriteByte(c)
	}
}

// Render a blackfriday markdown tree and return the output as a []byte.
// Returns ([]byte,nil) or (nil,err) if invalid input is encountered.
func (r *Renderer) Render(root *blackfriday.Node) ([]byte, error) {
	// if passed a full document, start on the first child node
	if root.Type == blackfriday.Document {
		root = root.FirstChild
	}

	for c := root; c != nil; c = c.Next {
		switch c.Type {
		case blackfriday.Heading:
			err := r.heading(c)
			if err != nil {
				return nil, err
			}
		case blackfriday.Paragraph:
			w := linewrap.New(r.out, r.cols)
			err := r.paragraph(w, c)
			if err != nil {
				return nil, err
			}
			w.Newline()
		case blackfriday.CodeBlock:
			r.codeBlock(c)
		case blackfriday.BlockQuote:
			w := linewrap.New(r.out, r.cols)
			err := r.blockQuote(w, c)
			if err != nil {
				return nil, err
			}
			r.out.WriteByte('\n')
		case blackfriday.List:
			w := linewrap.New(r.out, r.cols)
			err := r.list(w, c)
			if err != nil {
				return nil, err
			}
			w.Newline()
		}
	}

	return r.out.Bytes(), nil
}

// headingText checks that n and siblings are text nodes (there shouldn't
// be any siblings) and outputs all the text with whitespace flattened, or
// returns an error if an invalid (non Text) node is found
func (r *Renderer) headingText(n *blackfriday.Node) error {
	for p := n; p != nil; p = n.Next {
		if p.Type != blackfriday.Text {
			return errors.New("Headings may only contain text elements")
		}
		r.out.Write(trimFlattenSpaces(p.Literal))
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

func link(n *blackfriday.Node) (string, error) {
	dst := string(n.LinkData.Destination)
	if n.FirstChild == nil || n.FirstChild.Type != blackfriday.Text {
		return "", errors.New("Invalid Link Node")
	}
	text := string(trimFlattenSpaces(n.FirstChild.Literal))

	if strings.Compare(dst, text) == 0 {
		return ("<" + dst + ">"), nil
	}
	return ("[" + text + "](" + dst + ")"), nil
}

func (r *Renderer) paragraph(w *linewrap.Wrapper, n *blackfriday.Node) error {
	for c := n.FirstChild; c != nil; c = c.Next {
		switch c.Type {
		case blackfriday.Link:
			str, err := link(c)
			if err != nil {
				return err
			}
			w.WriteToken(str)
		case blackfriday.Text:
			s := strings.Replace(string(c.Literal), "\n", " ", -1)
			tokens := strings.Split(s, " ")
			w.WriteTokens(tokens)
		case blackfriday.Code:
			w.WriteToken("`" + string(c.Literal) + "`")
		}
	}
	w.TerminateLine()

	return nil
}

func (r *Renderer) blockQuote(w *linewrap.Wrapper, n *blackfriday.Node) error {
	subw := w.NewEmbedded("> ", "> ")
	first := true
	for c := n.FirstChild; c != nil; c = c.Next {
		if first == true {
			first = false
		} else {
			subw.Newline()
		}

		if c.Type == blackfriday.Paragraph {
			r.paragraph(subw, c)
			subw.TerminateLine()
		} else if c.Type == blackfriday.BlockQuote {
			r.blockQuote(subw, c)
			subw.TerminateLine()
		} else {
			return errors.New("BlockQuotes may only contain paragraphs or BlockQuotes")
		}
	}
	return nil
}

func (r *Renderer) codeBlock(n *blackfriday.Node) {
	fenceLength := 3
	if n.CodeBlockData.IsFenced && n.CodeBlockData.FenceLength > 0 {
		fenceLength = n.CodeBlockData.FenceLength
	}
	r.writeNBytes(fenceLength, '`')
	r.out.WriteByte('\n')
	r.out.Write(n.Literal)
	r.writeNBytes(fenceLength, '`')
	r.out.WriteByte('\n')
	r.out.WriteByte('\n')
}

func (r *Renderer) list(w *linewrap.Wrapper, n *blackfriday.Node) error {
	ordered := n.ListData.ListFlags&blackfriday.ListTypeOrdered > 0
	index := 1

	for c := n.FirstChild; c != nil; c = c.Next {
		if c.Type != blackfriday.Item {
			return errors.New("all list children must be 'Item' type")
		}
		if c.FirstChild.Type == blackfriday.Paragraph {
			if ordered {
				prefix := fmt.Sprintf("%d. ", index)
				subw := w.NewEmbedded(prefix, "   ")
				r.paragraph(subw, c.FirstChild)
			} else {
				subw := w.NewEmbedded("- ", "   ")
				r.paragraph(subw, c.FirstChild)
			}
		}
		if c.FirstChild.Next != nil && c.FirstChild.Next.Type == blackfriday.List {
			r.list(w.NewEmbedded("   ", "   "), c.FirstChild.Next)
		}
		index++
	}

	return nil
}
