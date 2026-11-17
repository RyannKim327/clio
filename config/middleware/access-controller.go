package middleware

import "github.com/gin-gonic/gin"

// RequestHandlerCookie guards a route with PERMISSION_COOKIE. It requires a
// "temporary" cookie whose value is at least 30 characters, otherwise the
// request is aborted with 403 Unauthorized.
func RequestHandlerCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, err := ctx.Cookie("temporary")

		if len(key) < 30 || err != nil {
			const html = `
			<!DOCTYPE html>
			<html>
				<head>
					<title>Restricted Area</title>
					<style>
						body {
							display: flex;
							flex-direction: row;
							height: 100dvh;
							width: 100dvw;
							align-items: center;
							justify-content: center;
							background-color: black;
							color: white;
						}
						p{
							font-size: 0.75rem;
							margin: 0;
						}	
					</style>
				</head>
				<body>
					<p><b style="font-size: 1.25rem;">403|</b> Restricted area</p>
				</body>
			</html>
		`
			ctx.Data(403, "text/html; charset=utf-8", []byte(html))
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
