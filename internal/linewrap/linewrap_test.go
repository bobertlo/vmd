package linewrap

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapperLorum(t *testing.T) {
	buf := &bytes.Buffer{}
	w := New(buf, 80)

	var testLorum string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	w.WriteToken("Hello, world!")
	w.Newline()
	w.Newline()
	w.WriteTokens(strings.Split(testLorum, " "))

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		assert.Less(t, len(line), 80)
	}

	expected := "Hello, world!\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor\nincididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis\nnostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\nDuis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu\nfugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in\nculpa qui officia deserunt mollit anim id est laborum."
	assert.Equal(t, expected, buf.String())
}

func TestWrapperMarkdown(t *testing.T) {
	buf := &bytes.Buffer{}
	w := New(buf, 80)

	var testLorum string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."

	// mock heading writer
	w.WriteNBytes(1, '#')
	w.WriteByte(' ')
	w.Write([]byte("Hello, world"))
	w.Write([]byte("\n\n"))

	w.WriteTokens(strings.Split(testLorum, " "))
	w.Write([]byte("\n\n"))

	// block quote
	blockquotewrapper := w.NewEmbedded("> ", "> ")
	blockquotewrapper.WriteTokens(strings.Split(testLorum, " "))
	blockquotewrapper.TerminateLine()
	blockquotewrapper.Newline()

	// code block inside block quote
	blockquotewrapper.WriteNBytes(3, '`')
	blockquotewrapper.TerminateLine()
	blockquotewrapper.Write([]byte("hello   world"))
	blockquotewrapper.Newline()
	blockquotewrapper.WriteNBytes(3, '`')
	blockquotewrapper.TerminateLine()

	w.TerminateLine()
	w.WriteTokens([]string{"goodbye", "world"})
	w.TerminateLine()

	expected := "# Hello, world\n\n Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do\neiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim\nveniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo\nconsequat.\n\n> Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor\n> incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis\n> nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n> \n> ```\n> hello   world\n> ```\n\ngoodbye world\n"
	assert.Equal(t, expected, buf.String())
}
