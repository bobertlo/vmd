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
}

func flattenSpaces(str []byte) []byte {
	re := regexp.MustCompile("  +")
	replaced := re.ReplaceAll(bytes.TrimSpace(str), []byte(" "))
	return replaced
}

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

func FileRenderer(path string) (*Renderer, error) {
	n, err := loadFile(path)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	r := &Renderer{
		out: buf,
		root: n,
	}
	return r, nil
}

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

func (r *Renderer) writeNBytes (n int, c byte) {
	for i := 0; i < n; i++ {
		r.out.WriteByte(c)
	}
}

func (r *Renderer) Render() error {
	for c := r.root.FirstChild; c != nil; c = c.Next {
		switch (c.Type) {
		case blackfriday.Heading:
			err := r.heading(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Renderer) headingText(n *blackfriday.Node) error {
	for p := n; p != nil; p = n.Next {
		if p.Type != blackfriday.Text {
			return errors.New("Headings may only contain text elements")
		}
		r.out.Write(flattenSpaces(p.Literal))
	}
	return nil
}

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



