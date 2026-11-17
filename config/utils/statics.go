package utils

import "github.com/gin-gonic/gin"

// Supported HTTP methods. Route.Method is matched against these values and
// defaults to METHOD_GET when empty.
const (
	// TODO: Methods
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_PATCH  = "PATCH"
	METHOD_DELETE = "DELETE"

	// TODO: Permissions
	PERMISSION_ADMIN  = "ADMIN"
	PERMISSION_ALL    = "ALL"
	PERMISSION_COOKIE = "COOKIE"
)

type PARAMS (map[string]any)

// HandlerSetup is the signature of a Constellation route handler. Returning an
// empty body and a nil parameter map means the handler wrote the response itself
// (for example with ctx.JSON); otherwise the result is rendered as an HTML page.
type HandlerSetup func(ctx *gin.Context) (string, PARAMS)
