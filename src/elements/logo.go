package elements

import (
	"fmt"
	"html"
)

func Logo(class string) string {
	return fmt.Sprintf(`<svg class="%s" viewBox="0 0 64 64" width="28" height="28" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true" focusable="false">
  <path d="M44 13 30 20 18 30 16 44 28 52 42 46" stroke="#fb923c" stroke-opacity="0.7" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
  <g fill="#f97316">
    <circle cx="44" cy="13" r="4"/>
    <circle cx="30" cy="20" r="3.5"/>
    <circle cx="18" cy="30" r="4"/>
    <circle cx="16" cy="44" r="3.5"/>
    <circle cx="28" cy="52" r="4"/>
    <circle cx="42" cy="46" r="3.5"/>
  </g>
  <g fill="#fff7ed" fill-opacity="0.9">
    <circle cx="44" cy="13" r="1.3"/>
    <circle cx="30" cy="20" r="1.1"/>
    <circle cx="18" cy="30" r="1.3"/>
    <circle cx="16" cy="44" r="1.1"/>
    <circle cx="28" cy="52" r="1.3"/>
    <circle cx="42" cy="46" r="1.1"/>
  </g>
  <path d="M34 26.5 36.5 32 42 34l-5.5 2L34 41.5 31.5 36 26 34l5.5-2z" fill="#fb923c"/>
</svg>`, html.EscapeString(class))
}
