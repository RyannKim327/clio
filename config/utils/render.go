package utils

import (
	"fmt"
	"html"
	"strings"

	"github.com/gin-gonic/gin"
)

// HTParser replaces every "{key}" token in body with the string form of the
// matching value in params. Values are substituted verbatim and are NOT
// HTML-escaped, so only pass trusted markup. Use Render for user-supplied data.
func HTParser(body string, params map[string]any) string {
	/*
	 * NOTE: Use this function for parsing HTML.
	 * Don't use this function for inputs to prevent XSS
	 */
	if len(params) == 0 {
		return body
	}

	pairs := make([]string, 0, len(params)*2)
	for key, value := range params {
		pairs = append(pairs, "{"+key+"}", fmt.Sprint(value))
	}

	return strings.NewReplacer(pairs...).Replace(body)
}

func Parser(body string, params map[string]any) string {
	/*
	 * NOTE: Use this function for parsing texts and escaping the HTML strings.
	 * This function is dedicated to user input-based parsing.
	 * Use this function for inputs to prevent XSS
	 */

	if len(params) == 0 {
		return body
	}

	pairs := make([]string, 0, len(params)*2)
	for key, value := range params {
		pairs = append(pairs, "{"+key+"}", html.EscapeString(fmt.Sprint(value)))
	}

	return strings.NewReplacer(pairs...).Replace(body)
}

// Render wraps body in the full HTML document shell, substitutes params with
// HTML escaping, and writes the result as text/html. It is called automatically
// for handlers that return a non-empty body or a non-nil parameter map.
func Render(ctx *gin.Context, title string, body string, params map[string]any) {
	if title == "" {
		title = "Istorya"
	}

	_html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		  <!-- Constellation -->
		  <!-- A framework developed with Go language -->
  		<!-- Developed by Ryann Kim Sesgundo -->
  		<head>
    		<meta charset="UTF-8">
    		<meta name="viewport" content="width=device-width, initial-scale=1">
   		<title>%s</title>
   		<script src="https://cdn.tailwindcss.com?plugins=typography"></script>
   		<link rel="preconnect" href="https://fonts.googleapis.com">
   		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
   		<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&amp;display=swap" rel="stylesheet">
         <link rel="stylesheet" href="/static/styles.css" />
   		</head>
   		<body>
   			<main class="site-main">%s</main>
   		</body>
		</html>
		`, title, Parser(body, params))
	ctx.Data(200, "text/html; charset=utf-8", []byte(_html))
}
