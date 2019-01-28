package mdformatter

import (
	"github.com/bobertlo/vmd/internal/renderer"
)

type MDFormatter struct {
	render *renderer.Renderer
}

func New(cols int) *MDFormatter {
	f := &MDFormatter{}
	f.render = renderer.New(cols)
	return f
}

func (f *MDFormatter) RenderBytes(input []byte) ([]byte, error) {
	return f.render.RenderBytes(input)
}
