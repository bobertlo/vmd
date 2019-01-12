package renderer

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/bobertlo/vmd/internal/linewrap"
	blackfriday "gopkg.in/bobertlo/blackfriday.v2"
)

// Renderer renders blackfriday markdown trees into []byte output
type Renderer struct {
	out  *bytes.Buffer
	cols int
}

// flattenSpaces removes all reduntant spaces from a []byte array, leaving
// single spaces
func flattenSpaces(str []byte) []byte {
	re := regexp.MustCompile("  +")
	return re.ReplaceAll(str, []byte(" "))
}

// trimFlattenSpaces removes all redundant spaces from a string, as well as
// any exterior spaces
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

// New creates a new markdown Renderer. cols specifies how many columns to
// wrap lines at.
func New(cols int) *Renderer {
	buf := new(bytes.Buffer)
	r := &Renderer{
		out:  buf,
		cols: cols,
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
		case blackfriday.Table:
			err := r.table(c)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported node type %s ignored", c.Type)
		}
	}

	// remove empty newline at end of file
	out := r.out.Bytes()
	if len(out) > 2 && out[len(out)-1] == '\n' && out[len(out)-2] == '\n' {
		return out[:len(out)-1], nil
	}
	if len(out) == 1 && (out[0] == '\n' || out[0] == ' ') {
		return []byte(""), nil
	}
	return out, nil
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

// link renders a Link node into a string
func link(n *blackfriday.Node) (string, error) {
	dst := string(n.LinkData.Destination)
	if n.FirstChild == nil {
		return "", errors.New("Invalid Link Node")
	}

	text, err := compileInline(n.FirstChild)
	if err != nil {
		return "", err
	}

	if strings.Compare(dst, text) == 0 {
		return ("<" + dst + ">"), nil
	}
	return ("[" + text + "](" + dst + ")"), nil
}

// image renders an Image node into a string
func image(n *blackfriday.Node) (string, error) {
	dst := string(n.LinkData.Destination)
	if n.FirstChild == nil {
		return "", errors.New("Invalid Image node")
	}

	text, err := compileInlineText(n.FirstChild)
	if err != nil {
		return "", err
	}

	return ("![" + text + "](" + dst + ")"), nil
}

// compileText joins text nodes. Takes a Node, and processes it and any
// siblings after it, returning a formatted (string, nil) or ("", err)
func compileInlineText(n *blackfriday.Node) (string, error) {
	if n == nil {
		return "", errors.New("invalid italic or bold formatting")
	}
	str := ""
	for c := n; c != nil; c = c.Next {
		if c.Type != blackfriday.Text {
			return "", errors.New("invalid italic or bold formatting")
		}
		str += string(c.Literal)
	}
	re := regexp.MustCompile("  +")
	str = strings.Replace(str, "\n", " ", -1)
	return re.ReplaceAllString(str, " "), nil
}

// inlineNode compiles the text inside a nodes children, and returns
// a string deliminated by the string "delim"
func inlineNode(n *blackfriday.Node, delim string) (string, error) {
	str, err := compileInlineText(n)
	if err != nil {
		return "", err
	}
	return (delim + str + delim), nil
}

// compileInline returns a string consisting of all Node n, and all of it's
// siblings, rendered (string, nil) or ("", err)
func compileInline(n *blackfriday.Node) (string, error) {
	var b strings.Builder

	for c := n; c != nil; c = c.Next {
		switch c.Type {
		case blackfriday.Link:
			if c.Parent.Type == blackfriday.Link {
				return "", errors.New("link text may not contain links")
			}
			str, err := link(c)
			if err != nil {
				return "", err
			}
			b.WriteString(str)
		case blackfriday.Image:
			str, err := image(c)
			if err != nil {
				return "", err
			}
			b.WriteString(str)
		case blackfriday.Text:
			str := strings.Replace(string(c.Literal), "\n", " ", -1)
			b.WriteString(str)
		case blackfriday.Emph:
			str, err := inlineNode(c.FirstChild, "*")
			if err != nil {
				return "", err
			}
			b.WriteString(str)
		case blackfriday.Strong:
			str, err := inlineNode(c.FirstChild, "**")
			if err != nil {
				return "", err
			}
			b.WriteString(str)
		case blackfriday.Code:
			b.WriteByte('`')
			str := strings.Replace(string(c.Literal), "\n", " ", -1)
			b.WriteString(str)
			b.WriteByte('`')
		default:
			return "", fmt.Errorf("Unsupported node type %s ignored", c.Type)
		}
	}

	return b.String(), nil
}

// wrapInline renders a node and all following children into a string, then
// tokenizes it based on whitespace and emits those tokens using a linewrapper,
// which is supplied to support recursion in lists and other types of blocks
func (r *Renderer) wrapInline(w *linewrap.Wrapper, n *blackfriday.Node) error {
	line, err := compileInline(n)
	if err != nil {
		return err
	}
	tokens := strings.Split(line, " ")
	w.WriteTokens(tokens)
	w.TerminateLine()
	return nil
}

// paragaph takes a Wrapper (because it is used to process code blocks and list
// bodies recursively) and a Paragraph Node, and renders all text and inline
// formatting nodes contained in the paragraph.
func (r *Renderer) paragraph(w *linewrap.Wrapper, n *blackfriday.Node) error {
	return r.wrapInline(w, n.FirstChild)
}

// blockQuote takes a BlockQuote node, and emits it. Only paragraphs and
// embedded BlockQuote nodes are admitted
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

// codeBlock emits a CodeBlock node.
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

// list emits a list, including any sublists recursively to a linewrap writer
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

func tableWidth(n *blackfriday.Node) (int, error) {
	head := n.FirstChild
	if head == nil || head.Type != blackfriday.TableHead {
		return 0, errors.New("invalid table structure")
	}

	row := head.FirstChild
	if row == nil || row.Type != blackfriday.TableRow {
		return 0, errors.New("invalid table structure")
	}

	cols := 0
	for c := row.FirstChild; c != nil; c = c.Next {
		if c.Type != blackfriday.TableCell {
			return 0, errors.New("invalid table structure")
		}
		cols++
	}

	return cols, nil
}

func (r *Renderer) table(n *blackfriday.Node) error {
	width, err := tableWidth(n)
	if err != nil {
		return err
	}

	max := make([]int, width)
	headData := make([]string, width)

	// process head
	hrow := n.FirstChild.FirstChild
	i := 0
	for c := hrow.FirstChild; c != nil; c = c.Next {
		if c.FirstChild != nil {
			str, err := compileInline(c.FirstChild)
			if err != nil {
				return err
			}
			headData[i] = str
			if len(str) > max[i] {
				max[i] = len(str)
			}
		}
		i++
	}
	if i < width {
		return errors.New("table row too short")
	}

	values := [][]string{}

	body := n.FirstChild.Next
	if body == nil {
		return errors.New("invalid table structure")
	}

	// process rows
	rows := 0
	for row := body.FirstChild; row != nil; row = row.Next {
		rowData := make([]string, width)
		i = 0
		for c := row.FirstChild; c != nil; c = c.Next {
			if i > width {
				return errors.New("table row too long")
			}

			if c.FirstChild != nil {
				str, err := compileInline(c.FirstChild)
				if err != nil {
					return err
				}
				rowData[i] = str
				if len(str) > max[i] {
					max[i] = len(str)
				}
			}
			i++
		}
		if i < width {
			return errors.New("table row too short")
		}
		values = append(values, rowData)
		rows++
	}

	if rows == 0 {
		return errors.New("invalid table structure")
	}

	// output table head
	for i := 0; i < width; i++ {
		fmt.Fprintf(r.out, "| %s", headData[i])
		r.writeNBytes(max[i]-len(headData[i])+1, ' ')
	}
	fmt.Fprintln(r.out, "|")

	for i := 0; i < width; i++ {
		r.out.WriteByte('|')
		r.writeNBytes(max[i]+2, '-')
	}
	fmt.Fprintln(r.out, "|")

	for i := range values {
		for j := 0; j < width; j++ {
			fmt.Fprintf(r.out, "| %s", values[i][j])
			r.writeNBytes(max[j]-len(values[i][j])+1, ' ')
		}
		fmt.Fprintln(r.out, "|")
	}

	r.out.WriteByte('\n')

	return nil
}
