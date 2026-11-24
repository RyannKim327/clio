package main

import (
	"embed"

	"constellation/site"
	constellation "constellation/src"
)

// StaticFS contains the assets under static, embedded at build time.
//
//go:embed all:static
var StaticFS embed.FS

func main() {
	for _, endpoint := range site.Endpoints(StaticFS) {
		constellation.AddEndpoint(endpoint)
	}

	css, err := StaticFS.ReadFile("static/styles.css")
	if err != nil {
		panic(err)
	}

	constellation.StartServer(css)
}
