package utils

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var storyMarkdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

var (
	audioToken = regexp.MustCompile(`\{(https?://[^\s{}]+)\}`)
	stopToken  = regexp.MustCompile(`(?i)\{stop\}`)
	unsafeURL  = regexp.MustCompile(`(?i)(href|src)="(?:javascript|data|vbscript):[^"]*"`)
)

// RenderStory converts user-authored Markdown into HTML and turns
// {https://example.com/music.mp3} tokens into audio markers. A {stop} token
// marks where the music should stop. Raw HTML and unsafe link schemes are
// neutralized, so the result can be injected into a page.
func RenderStory(source string) string {
	if strings.TrimSpace(source) == "" {
		return `<p class="story-empty">Nothing written yet.</p>`
	}

	markers := make([]string, 0)
	guarded := audioToken.ReplaceAllStringFunc(source, func(match string) string {
		url := strings.TrimSuffix(strings.TrimPrefix(match, "{"), "}")
		key := fmt.Sprintf("ISTORYAAUDIO%dENDISTORYA", len(markers))
		markers = append(markers,
			`<span class="audio-marker" data-audio-src="`+html.EscapeString(url)+`" role="button" tabindex="0" title="Music starts here">♪ music</span>`,
		)
		return key
	})

	const stopKey = "ISTORYAAUDIOSTOPENDISTORYA"
	guarded = stopToken.ReplaceAllString(guarded, stopKey)

	var out bytes.Buffer
	if err := storyMarkdown.Convert([]byte(guarded), &out); err != nil {
		return `<p class="story-empty">Unable to render this chapter.</p>`
	}

	result := out.String()
	for i, marker := range markers {
		key := fmt.Sprintf("ISTORYAAUDIO%dENDISTORYA", i)
		result = strings.ReplaceAll(result, key, marker)
	}
	result = strings.ReplaceAll(result, stopKey,
		`<span class="audio-stop-marker" data-audio-stop="true" title="Music stops here">■ stop music</span>`,
	)
	result = unsafeURL.ReplaceAllString(result, `$1="#"`)
	return result
}
