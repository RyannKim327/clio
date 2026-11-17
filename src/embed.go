package src

import "embed"

// StaticFS contains the assets under src/static, embedded at build time.
//
//go:embed all:static
var StaticFS embed.FS
