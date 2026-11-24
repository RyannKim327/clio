package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var stylesheet string

// SetStylesheet provides the CSS processed by the Tailwind browser runtime.
func SetStylesheet(css string) {
	stylesheet = css
}

func SendApi(w http.ResponseWriter, data map[string]any) {
	json.NewEncoder(w).Encode(data)
}

func WebParser(w http.ResponseWriter, title string, html string) (int, error) {
	if title == "" {
		title = "Constellation"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Powered-By", "Constellation")
	w.Header().Set("Developed-Under", "Master Piece of Paper")
	w.Header().Set("Developed-By", "Ryann Kim Sesgundo")

	_html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		  <!-- Constellation -->
		  <!-- A framework developed with Go language -->
  		  <!-- Developed by Ryann Kim Sesgundo -->
  		  <head>
    		  <meta charset="UTF-8">
    		  <meta name="viewport" content="width=device-width, initial-scale=1">
    		  <meta name="color-scheme" content="dark">
    		  <title>%s</title>
    		  <script src="https://cdn.tailwindcss.com?plugins=typography"></script>
         <style>%s</style>
  		  </head>
  		  <body class="min-h-screen bg-[#0b0b0f] text-slate-50 antialiased">
    		  <main class="min-h-screen">%s</main>
  		  </body>
		</html>
		`, title, stylesheet, html)
	return w.Write([]byte(_html))
}
