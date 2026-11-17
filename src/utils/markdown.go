package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// Markdown converts trusted Markdown source into HTML. The source is treated as
// documentation authored in this repository, so the result can be injected with
// HTParser. Never pass untrusted user input here.
func Markdown(source string) string {
	var out bytes.Buffer
	if err := markdownRenderer.Convert([]byte(source), &out); err != nil {
		return "<p>Unable to render this document.</p>"
	}
	return out.String()
}
