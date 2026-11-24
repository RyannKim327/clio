package middleware

import "net/http"

type Permission string

const (
	Unauthenticated Permission = "__public__"
	Authenticated   Permission = "__private__"
	Authorized      Permission = "__protected__"
)

type AccessCheck func(r *http.Request) bool

var (
	AuthenticatedCheck AccessCheck = func(r *http.Request) bool { return false }
	AuthorizedCheck    AccessCheck = func(r *http.Request) bool { return false }
)

func AccessController(permission Permission, r *http.Request) bool {
	switch permission {
	case "", Unauthenticated:
		return true
	case Authenticated:
		return AuthenticatedCheck != nil && AuthenticatedCheck(r)
	case Authorized:
		return AuthenticatedCheck != nil && AuthorizedCheck != nil && AuthenticatedCheck(r) && AuthorizedCheck(r)
	default:
		return false
	}
}
