package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// BotProtection rejects requests whose User-Agent matches a known scraper agent
// (scrapy, python-requests, httpx) with 403 Access denied. This is a simple
// heuristic and should be hardened before being relied on in production.
func BotProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := strings.ToLower(c.GetHeader("User-Agent"))

		blocked := []string{
			"scrapy",
			"python-requests",
			"httpx",
		}

		if !strings.ContainsAny(strings.Join(blocked, ","), ua) {
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
			c.Data(403, "text/html; charset=utf-8", []byte(html))

			c.Abort()
			return
		}
		c.Next()
	}
}
