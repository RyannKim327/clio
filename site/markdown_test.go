package site

import (
	"strings"
	"testing"
)

func TestRenderMarkdownEscapesRawHTML(t *testing.T) {
	result := RenderMarkdown("<script>alert(1)</script> and **bold**")
	if strings.Contains(result, "<script>") {
		t.Fatalf("raw HTML was not escaped: %s", result)
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Fatalf("expected escaped script tag, got: %s", result)
	}
	if !strings.Contains(result, "<strong>bold</strong>") {
		t.Fatalf("expected strong text, got: %s", result)
	}
}

func TestRenderMarkdownBlocks(t *testing.T) {
	source := "# Chapter\n\nA *quiet* line.\n\n> A note.\n\n- one\n- two\n\n| Left | Right |\n| :--- | ---: |\n| a | b |\n\n```go\nfmt.Println(\"hi\")\n```"
	result := RenderMarkdown(source)

	for _, expected := range []string{
		"<h1>Chapter</h1>",
		"<em>quiet</em>",
		"<blockquote class=\"story-quote\">",
		"<ul class=\"story-list\">",
		"<table class=\"story-table\">",
		"<pre class=\"story-pre\">",
		"fmt.Println(&#34;hi&#34;)",
	} {
		if !strings.Contains(result, expected) {
			t.Fatalf("expected %q in output: %s", expected, result)
		}
	}
}

func TestRenderMarkdownAudioMarkers(t *testing.T) {
	result := RenderMarkdown("Listen {https://cdn.example.com/theme.mp3} closely. {stop}")
	if !strings.Contains(result, `data-audio-src="https://cdn.example.com/theme.mp3"`) {
		t.Fatalf("missing audio marker: %s", result)
	}
	if !strings.Contains(result, `data-audio-stop="true"`) {
		t.Fatalf("missing stop marker: %s", result)
	}
}

func TestRenderMarkdownRejectsUnsafeLinks(t *testing.T) {
	result := RenderMarkdown("[click](javascript:alert(1))")
	if strings.Contains(result, `href="javascript`) {
		t.Fatalf("unsafe link survived: %s", result)
	}
	if !strings.Contains(result, "[click]") {
		t.Fatalf("expected the link text to remain, got: %s", result)
	}
}

func TestRenderMarkdownKeepsMarkersOutOfCode(t *testing.T) {
	result := RenderMarkdown("Use `{stop}` literally.\n\n```\n{https://cdn.example.com/theme.mp3}\n```")
	if strings.Contains(result, "data-audio-stop") || strings.Contains(result, "data-audio-src") {
		t.Fatalf("markers leaked into code: %s", result)
	}
	if !strings.Contains(result, "{stop}") || !strings.Contains(result, "{https://cdn.example.com/theme.mp3}") {
		t.Fatalf("original markers were not preserved: %s", result)
	}
	if strings.ContainsRune(result, '\x00') {
		t.Fatalf("sentinel leaked into output: %s", result)
	}
}

func TestRenderMarkdownEmpty(t *testing.T) {
	if result := RenderMarkdown("   "); !strings.Contains(result, "Nothing written yet.") {
		t.Fatalf("unexpected empty output: %s", result)
	}
}
