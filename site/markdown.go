package site

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	audioToken      = regexp.MustCompile(`\{(https?://[^\s{}<>"]+)\}`)
	stopToken       = regexp.MustCompile(`(?i)\{stop\}`)
	headingLine     = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	horizontalRule  = regexp.MustCompile(`^(-{3,}|\*{3,}|_{3,})$`)
	fenceLine       = regexp.MustCompile("^(```|~~~)(.*)$")
	languageName    = regexp.MustCompile(`^[A-Za-z0-9_+-]+$`)
	unorderedItem   = regexp.MustCompile(`^[-*+]\s+(.*)$`)
	orderedItem     = regexp.MustCompile(`^\d+[.)]\s+(.*)$`)
	tableDelimiter  = regexp.MustCompile(`^\|?\s*:?-{3,}:?\s*(\|\s*:?-{3,}:?\s*)+\|?$`)
	autolinkPattern = regexp.MustCompile(`^<(https?://[^>\s]+)>`)
	taskItemLine    = regexp.MustCompile(`^\[([ xX])\]\s+(.*)$`)
	codeRegion      = regexp.MustCompile(`(?is)<pre\b[^>]*>.*?</pre>|<code\b[^>]*>.*?</code>`)
)

const stopMarker = "\x00clio-stop\x00"

func RenderMarkdown(source string) string {
	if strings.TrimSpace(source) == "" {
		return `<p class="story-empty">Nothing written yet.</p>`
	}

	source = strings.ReplaceAll(source, "\x00", "")
	replacements := map[string]string{
		stopMarker: `<span class="audio-stop-marker" data-audio-stop="true" title="Music stops here">&#9632; stop music</span>`,
	}
	originals := map[string]string{stopMarker: "{stop}"}

	guarded := audioToken.ReplaceAllStringFunc(source, func(match string) string {
		address := match[1 : len(match)-1]
		if !safeMediaURL(address) {
			return match
		}
		key := fmt.Sprintf("\x00clio-audio-%d\x00", len(replacements))
		replacements[key] = `<span class="audio-marker" data-audio-src="` + html.EscapeString(address) + `" role="button" tabindex="0" title="Music starts here">&#9834; music</span>`
		originals[key] = match
		return key
	})
	guarded = stopToken.ReplaceAllStringFunc(guarded, func(match string) string {
		originals[stopMarker] = match
		return stopMarker
	})

	return injectMarkers(renderBlocks(guarded), replacements, originals)
}

func injectMarkers(rendered string, replacements, originals map[string]string) string {
	var out strings.Builder
	last := 0

	for _, region := range codeRegion.FindAllStringIndex(rendered, -1) {
		out.WriteString(applyMarkers(rendered[last:region[0]], replacements))
		out.WriteString(applyMarkers(rendered[region[0]:region[1]], originals))
		last = region[1]
	}
	out.WriteString(applyMarkers(rendered[last:], replacements))
	return out.String()
}

func applyMarkers(segment string, replacements map[string]string) string {
	for key, value := range replacements {
		segment = strings.ReplaceAll(segment, key, value)
	}
	return segment
}

func renderBlocks(source string) string {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	var out strings.Builder

	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}

		if match := fenceLine.FindStringSubmatch(line); match != nil {
			fence, language := match[1], strings.TrimSpace(match[2])
			if !languageName.MatchString(language) {
				language = ""
			}
			i++
			var code strings.Builder
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), fence) {
				code.WriteString(lines[i])
				code.WriteByte('\n')
				i++
			}
			if i < len(lines) {
				i++
			}
			class := "story-code"
			if language != "" {
				class = "language-" + language
			}
			out.WriteString(`<pre class="story-pre"><code class="` + class + `">` + html.EscapeString(strings.TrimSuffix(code.String(), "\n")) + `</code></pre>`)
			continue
		}

		if match := headingLine.FindStringSubmatch(line); match != nil {
			level := len(match[1])
			out.WriteString(fmt.Sprintf("<h%d>%s</h%d>", level, renderInline(match[2]), level))
			i++
			continue
		}

		if horizontalRule.MatchString(line) {
			out.WriteString(`<hr class="story-rule" />`)
			i++
			continue
		}

		if strings.HasPrefix(line, ">") {
			var quoted []string
			for i < len(lines) {
				candidate := strings.TrimSpace(lines[i])
				if !strings.HasPrefix(candidate, ">") {
					break
				}
				candidate = strings.TrimPrefix(candidate, ">")
				quoted = append(quoted, strings.TrimPrefix(candidate, " "))
				i++
			}
			out.WriteString(`<blockquote class="story-quote">` + renderBlocks(strings.Join(quoted, "\n")) + `</blockquote>`)
			continue
		}

		if strings.Contains(line, "|") && i+1 < len(lines) && tableDelimiter.MatchString(strings.TrimSpace(lines[i+1])) {
			table, next := renderTable(lines, i)
			out.WriteString(table)
			i = next
			continue
		}

		if unorderedItem.MatchString(line) || orderedItem.MatchString(line) {
			ordered := orderedItem.MatchString(line)
			tag := "ul"
			matcher := unorderedItem
			if ordered {
				tag = "ol"
				matcher = orderedItem
			}
			var items []string
			for i < len(lines) {
				match := matcher.FindStringSubmatch(strings.TrimSpace(lines[i]))
				if match == nil {
					break
				}
				items = append(items, renderListItem(match[1]))
				i++
			}
			out.WriteString("<" + tag + " class=\"story-list\">" + strings.Join(items, "") + "</" + tag + ">")
			continue
		}

		var paragraph []string
		for i < len(lines) && !startsBlock(lines, i) {
			current := strings.TrimSpace(lines[i])
			if current == "" {
				break
			}
			hardBreak := strings.HasSuffix(lines[i], "  ")
			paragraph = append(paragraph, renderInline(strings.TrimSpace(current)))
			if hardBreak {
				paragraph = append(paragraph, "<br />")
			}
			i++
		}
		if len(paragraph) > 0 {
			out.WriteString("<p>" + strings.Join(paragraph, " ") + "</p>")
		}
	}

	return out.String()
}

func startsBlock(lines []string, index int) bool {
	line := strings.TrimSpace(lines[index])
	if line == "" || strings.HasPrefix(line, ">") || headingLine.MatchString(line) || horizontalRule.MatchString(line) || fenceLine.MatchString(line) {
		return true
	}
	if unorderedItem.MatchString(line) || orderedItem.MatchString(line) {
		return true
	}
	return strings.Contains(line, "|") && index+1 < len(lines) && tableDelimiter.MatchString(strings.TrimSpace(lines[index+1]))
}

func renderListItem(text string) string {
	if match := taskItemLine.FindStringSubmatch(text); match != nil {
		checked := ""
		if strings.EqualFold(match[1], "x") {
			checked = " checked"
		}
		return `<li class="task-item"><input type="checkbox" disabled` + checked + ` /> ` + renderInline(match[2]) + `</li>`
	}
	return "<li>" + renderInline(text) + "</li>"
}

func renderTable(lines []string, start int) (string, int) {
	headers := tableCells(lines[start])
	delimiters := tableCells(lines[start+1])
	alignments := make([]string, len(headers))
	for i := range alignments {
		if i < len(delimiters) {
			alignments[i] = tableAlignment(delimiters[i])
		}
	}

	var out strings.Builder
	out.WriteString(`<div class="story-table-wrap"><table class="story-table"><thead><tr>`)
	for i, cell := range headers {
		out.WriteString(`<th` + alignmentAttribute(alignments[i]) + `>` + renderInline(cell) + `</th>`)
	}
	out.WriteString(`</tr></thead><tbody>`)

	i := start + 2
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.Contains(line, "|") {
			break
		}
		cells := tableCells(lines[i])
		out.WriteString(`<tr>`)
		for j, cell := range cells {
			alignment := ""
			if j < len(alignments) {
				alignment = alignments[j]
			}
			out.WriteString(`<td` + alignmentAttribute(alignment) + `>` + renderInline(cell) + `</td>`)
		}
		out.WriteString(`</tr>`)
		i++
	}
	out.WriteString(`</tbody></table></div>`)
	return out.String(), i
}

func tableCells(line string) []string {
	value := strings.TrimSpace(line)
	value = strings.TrimSuffix(strings.TrimPrefix(value, "|"), "|")
	parts := strings.Split(value, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func tableAlignment(delimiter string) string {
	left, right := strings.HasPrefix(delimiter, ":"), strings.HasSuffix(delimiter, ":")
	switch {
	case left && right:
		return "center"
	case right:
		return "right"
	case left:
		return "left"
	default:
		return ""
	}
}

func alignmentAttribute(alignment string) string {
	if alignment == "" {
		return ""
	}
	return ` style="text-align:` + alignment + `"`
}

func renderInline(source string) string {
	var out strings.Builder

	for i := 0; i < len(source); {
		switch {
		case source[i] == '\\' && i+1 < len(source) && strings.ContainsRune(`\*_{}[]()#+-.!>`, rune(source[i+1])):
			out.WriteString(html.EscapeString(string(source[i+1])))
			i += 2

		case source[i] == '`':
			if end := strings.IndexByte(source[i+1:], '`'); end > 0 {
				out.WriteString(`<code class="story-code">` + html.EscapeString(source[i+1:i+1+end]) + `</code>`)
				i += end + 2
			} else {
				out.WriteString("&#96;")
				i++
			}

		case strings.HasPrefix(source[i:], "**") || strings.HasPrefix(source[i:], "__"):
			delimiter := source[i : i+2]
			if end := strings.Index(source[i+2:], delimiter); end > 0 {
				out.WriteString(`<strong>` + renderInline(source[i+2:i+2+end]) + `</strong>`)
				i += end + 4
			} else {
				out.WriteString(html.EscapeString(delimiter))
				i += 2
			}

		case strings.HasPrefix(source[i:], "~~"):
			if end := strings.Index(source[i+2:], "~~"); end > 0 {
				out.WriteString(`<del>` + renderInline(source[i+2:i+2+end]) + `</del>`)
				i += end + 4
			} else {
				out.WriteString("~~")
				i += 2
			}

		case strings.HasPrefix(source[i:], "=="):
			if end := strings.Index(source[i+2:], "=="); end > 0 {
				out.WriteString(`<mark>` + renderInline(source[i+2:i+2+end]) + `</mark>`)
				i += end + 4
			} else {
				out.WriteString("==")
				i += 2
			}

		case source[i] == '*' || source[i] == '_':
			delimiter := string(source[i])
			if end := strings.Index(source[i+1:], delimiter); end > 0 {
				out.WriteString(`<em>` + renderInline(source[i+1:i+1+end]) + `</em>`)
				i += end + 2
			} else {
				out.WriteString(html.EscapeString(delimiter))
				i++
			}

		case strings.HasPrefix(source[i:], "!["):
			if rendered, next, ok := renderImage(source, i); ok {
				out.WriteString(rendered)
				i = next
			} else {
				out.WriteString("!")
				i++
			}

		case source[i] == '[':
			if rendered, next, ok := renderLink(source, i); ok {
				out.WriteString(rendered)
				i = next
			} else {
				out.WriteString("[")
				i++
			}

		case source[i] == '<':
			if match := autolinkPattern.FindStringSubmatch(source[i:]); match != nil {
				out.WriteString(renderAnchor(match[1], html.EscapeString(match[1])))
				i += len(match[0])
			} else {
				out.WriteString("&lt;")
				i++
			}

		default:
			j := i
			for j < len(source) && !isInlineSpecial(source[j]) {
				j++
			}
			if j == i {
				j++
			}
			out.WriteString(html.EscapeString(source[i:j]))
			i = j
		}
	}

	return out.String()
}

func renderLink(source string, start int) (string, int, bool) {
	labelEnd := strings.IndexByte(source[start+1:], ']')
	if labelEnd < 0 {
		return "", start, false
	}
	labelEnd = start + 1 + labelEnd
	if labelEnd+1 >= len(source) || source[labelEnd+1] != '(' {
		return "", start, false
	}
	targetEnd := strings.IndexByte(source[labelEnd+2:], ')')
	if targetEnd < 0 {
		return "", start, false
	}
	targetEnd = labelEnd + 2 + targetEnd

	target := strings.TrimSpace(source[labelEnd+2 : targetEnd])
	if safeHref(target) == "" {
		return "", start, false
	}
	label := renderInline(source[start+1 : labelEnd])
	return renderAnchor(target, label), targetEnd + 1, true
}

func renderImage(source string, start int) (string, int, bool) {
	altEnd := strings.IndexByte(source[start+2:], ']')
	if altEnd < 0 {
		return "", start, false
	}
	altEnd = start + 2 + altEnd
	if altEnd+1 >= len(source) || source[altEnd+1] != '(' {
		return "", start, false
	}
	targetEnd := strings.IndexByte(source[altEnd+2:], ')')
	if targetEnd < 0 {
		return "", start, false
	}
	targetEnd = altEnd + 2 + targetEnd

	target := strings.TrimSpace(source[altEnd+2 : targetEnd])
	href := safeHref(target)
	if href == "" || (!safeMediaURL(target) && !isRelativeURL(target)) {
		return "", start, false
	}
	alt := html.EscapeString(source[start+2 : altEnd])
	return `<img class="story-image" src="` + href + `" alt="` + alt + `" loading="lazy" />`, targetEnd + 1, true
}

func renderAnchor(target, label string) string {
	href := safeHref(target)
	if href == "" {
		return label
	}
	return `<a href="` + href + `">` + label + `</a>`
}

func safeHref(target string) string {
	target = strings.TrimSpace(target)
	if target == "" || strings.ContainsAny(target, "\"'<>") {
		return ""
	}

	parsed, err := url.Parse(target)
	if err != nil {
		return ""
	}
	if parsed.Scheme == "" {
		return html.EscapeString(target)
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https", "mailto":
		return html.EscapeString(target)
	default:
		return ""
	}
}

func safeMediaURL(target string) bool {
	parsed, err := url.Parse(strings.TrimSpace(target))
	if err != nil {
		return false
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
		return parsed.Host != ""
	default:
		return false
	}
}

func isRelativeURL(target string) bool {
	return strings.HasPrefix(target, "/") || strings.HasPrefix(target, "./") || strings.HasPrefix(target, "../") || strings.HasPrefix(target, "#")
}

func isInlineSpecial(char byte) bool {
	switch char {
	case '\\', '`', '*', '_', '~', '!', '[', '<', '=':
		return true
	default:
		return false
	}
}
