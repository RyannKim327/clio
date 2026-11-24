package constellation

import (
	"net/http"

	"constellation/src/middleware"
)

type (
	HandlerWeb     func(r *Request) string
	HandlerApi     func(r *Request, w Response)
	MethodType     string
	PermissionType = middleware.Permission
	Request        http.Request
	Response       http.ResponseWriter
	APIResponse    map[string]any
)

type HandlerSetup struct {
	Web HandlerWeb
	Api HandlerApi
}

const (
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	PUT    MethodType = "PUT"
	PATCH  MethodType = "PATCH"
	DELETE MethodType = "DELETE"
)

const (
	Unauthenticated = middleware.Unauthenticated
	Authenticated   = middleware.Authenticated
	Authorized      = middleware.Authorized
)

type WebComponent struct {
	Class      string
	Attributes map[string]string
	Children   []any
	Text       string
}

type Endpoint struct {
	Path       string
	Method     MethodType
	Handler    HandlerSetup
	Permission PermissionType
	Title      string
}

var Endpoints []Endpoint

func AddEndpoint(endpoint Endpoint) {
	if endpoint.Path == "" {
		endpoint.Path = "/"
	}

	if endpoint.Method == "" {
		endpoint.Method = GET
	}

	if endpoint.Permission == "" {
		endpoint.Permission = Unauthenticated
	}

	Endpoints = append(Endpoints, endpoint)
}
