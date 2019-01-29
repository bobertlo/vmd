package mdformatter

import (
	"github.com/bobertlo/vmd/internal/renderer"
)

// MDFormatter is a struct with methods to format markdown
type MDFormatter struct {
	render *renderer.Renderer
}

// New returns a new MDFormatter which wraps lines at a specified number
// of columns
func New(cols int) *MDFormatter {
	f := &MDFormatter{}
	f.render = renderer.New(cols)
	return f
}

// RenderBytes renders a markdown []byte slice and returns the formatted
// outout ([]byte, nil) or an error (nil, error)
func (f *MDFormatter) RenderBytes(input []byte) ([]byte, error) {
	return f.render.RenderBytes(input)
}
