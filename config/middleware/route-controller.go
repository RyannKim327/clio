package middleware

import "github.com/RyannKim327/constellation/config/utils"

// TODO: Just create a structure inside and call it outside

// Route describes an HTTP endpoint. Method defaults to GET and Path is
// lowercased during registration. Permission selects the access guard and
// Handler produces the response body and template parameters.
type Route struct {
	Title      string
	Method     string
	Path       string
	Permission string
	Handler    utils.HandlerSetup
}
