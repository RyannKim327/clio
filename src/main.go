package constellation

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"constellation/src/middleware"
)

func StartServer(css []byte) {
	serve := http.NewServeMux()

	for _, endpoint := range Endpoints {
		serve.HandleFunc(fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path), func(w http.ResponseWriter, r *http.Request) {
			if !middleware.AccessController(endpoint.Permission, r) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			request := (*Request)(r)

			switch {
			case endpoint.Handler.Api != nil:
				endpoint.Handler.Api(request, w)
			case endpoint.Handler.Web != nil:
				middleware.WebParser(w, endpoint.Title, endpoint.Handler.Web(request))
			}
		})
	}
	println("System Powering Up")
	println("Powered by Go")
	println("Developed under Master Piece of Paper")
	println("Welcome to Constellation Server")

	middleware.SetStylesheet(string(css))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, serve))
}
