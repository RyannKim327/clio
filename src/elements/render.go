package elements

import (
	"fmt"
	"html"
	"sort"
	"strings"

	constellation "constellation/src"
)

var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

func render(tag string, element constellation.WebComponent) string {
	var builder strings.Builder

	builder.WriteByte('<')
	builder.WriteString(tag)
	builder.WriteString(attributes(element))
	builder.WriteByte('>')

	if voidElements[tag] {
		return builder.String()
	}

	builder.WriteString(content(element))
	builder.WriteString("</")
	builder.WriteString(tag)
	builder.WriteByte('>')

	return builder.String()
}

func attributes(element constellation.WebComponent) string {
	keys := make([]string, 0, len(element.Attributes)+1)

	if element.Class != "" {
		keys = append(keys, "class")
	}
	for key := range element.Attributes {
		if key != "class" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		value := element.Attributes[key]
		if key == "class" {
			value = element.Class
		}

		builder.WriteByte(' ')
		builder.WriteString(html.EscapeString(key))
		builder.WriteString(`="`)
		builder.WriteString(html.EscapeString(value))
		builder.WriteByte('"')
	}

	return builder.String()
}

func content(element constellation.WebComponent) string {
	var builder strings.Builder

	if element.Text != "" {
		builder.WriteString(html.EscapeString(element.Text))
	}
	for _, child := range element.Children {
		builder.WriteString(renderChild(child))
	}

	return builder.String()
}

func renderChild(child any) string {
	switch value := child.(type) {
	case nil:
		return ""
	case string:
		return value
	case []any:
		var builder strings.Builder
		for _, nested := range value {
			builder.WriteString(renderChild(nested))
		}
		return builder.String()
	case []string:
		var builder strings.Builder
		for _, nested := range value {
			builder.WriteString(nested)
		}
		return builder.String()
	case fmt.Stringer:
		return html.EscapeString(value.String())
	default:
		return html.EscapeString(fmt.Sprint(value))
	}
}
