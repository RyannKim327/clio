package constellation

import (
	"io/fs"

	mw "github.com/RyannKim327/constellation/config/middleware"
)

// Constellation bootstraps the application: it sets up CORS, registers every
// supplied route, and starts the HTTP server. static holds the embedded assets
// served at /static.
func Constellation(static fs.FS, _routes []mw.Route) {
	/*
	 * TODO: There's nothing to do with this file
	 * Everything is set
	 */

	// TODO: Call this first before the registration of each routes
	mw.CorsSetup()

	for _, route := range _routes {
		mw.Register(route)
	}

	mw.StartServer(static)
}
