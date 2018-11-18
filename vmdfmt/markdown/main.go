package markdown

import (
	"bytes"
	"io/ioutil"
	"gopkg.in/russross/blackfriday.v2"
)

type Renderer struct {
	out *bytes.Buffer
	root *blackfriday.Node
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
	r.Render()
	return r.out.Bytes(), nil
}

func (r *Renderer) writeNBytes (n int, c byte) {
	for i := 0; i < n; i++ {
		r.out.WriteByte(c)
	}
}

func (r *Renderer) Render() {
	for c := r.root.FirstChild; c != nil; c = c.Next {
		switch (c.Type) {
		case blackfriday.Heading:
			r.heading(c)
		}
	}
}

func (r *Renderer) heading(n *blackfriday.Node) {
	r.writeNBytes(n.HeadingData.Level, '#')
	r.out.WriteByte(' ')
	if n.FirstChild.Type != blackfriday.Text {
		panic("invalid")
	}
}



