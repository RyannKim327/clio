package middleware

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RyannKim327/constellation/config/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var app = gin.Default()

// CorsSetup installs CORS middleware using ALLOWED_ORIGINS (comma-separated) or
// local development defaults, plus a response-header middleware. Call it before
// registering routes so the middleware applies to them.
func CorsSetup() {
	origins := []string{"https://constellation-red.vercel.app", "http://localhost:3000", "http://localhost:5173"}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		origins = strings.Split(v, ",")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.Use(responseHeaderMiddleware())
}

func _routesHandler(title string, handler utils.HandlerSetup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, params := handler(ctx)
		if body == "" && params == nil {
			return
		}
		utils.Render(ctx, title, body, params)
	}
}

// Register normalizes a route and attaches it to the engine. It defaults the
// method to GET, lowercases the path, appends "/submit" to guarded POST routes,
// resolves the permission, and wires the matching middleware chain.
func Register(routes Route) {
	/*
	 * INFO: To register each of endpoint from endpoint package/folder
	 * The purpose is to make the files cleaner and better looking
	 */

	// TODO: Make GET Request as default
	if routes.Method == "" {
		routes.Method = utils.METHOD_GET
	}

	if routes.Title == "" {
		routes.Title = "Constellation App"
	}

	// TODO: To make the "Method" case insensitive
	method := strings.ToUpper(routes.Method)

	// TODO: To make sure that the endpoint is in lowercase
	routes.Path = strings.ToLower(routes.Path)

	if !strings.HasPrefix(routes.Path, "/") {
		routes.Path = fmt.Sprintf("/%s", routes.Path)
	}

	// TODO: To automatically add "/submit" in every required post requests
	if method == strings.ToUpper(utils.METHOD_POST) && strings.ToUpper(routes.Permission) != utils.PERMISSION_ALL && !strings.HasSuffix(routes.Path, "/submit") {

		// TODO: To prevent double slash in endpoints
		if !strings.HasSuffix(routes.Path, "/") {
			routes.Path += "/"
		}

		routes.Path += "submit"
	}
	//
	// if (strings.ToUpper(method) == strings.ToUpper(utils.METHOD_PUT) || strings.ToUpper(method) == strings.ToUpper(utils.METHOD_PATCH)) && strings.ToUpper(routes.Permission) != utils.PERMISSION_ALL && !strings.HasSuffix(routes.Path, "/edit") {
	//
	// 	// TODO: To prevent double slash in endpoints
	// 	if !strings.HasSuffix(routes.Path, "/") {
	// 		routes.Path += "/"
	// 	}
	//
	// 	routes.Path += "edit"
	// }

	// TODO: To setup default permission
	if routes.Permission == "" {
		if routes.Method == utils.METHOD_POST {
			routes.Permission = utils.PERMISSION_COOKIE
		} else {
			routes.Permission = utils.PERMISSION_ALL
		}
	}

	switch method {
	case utils.METHOD_GET:
		switch routes.Permission {
		case utils.PERMISSION_COOKIE:
			app.GET(routes.Path, BotProtection(), RequestHandlerCookie(), _routesHandler(routes.Title, routes.Handler))
		default:
			app.GET(routes.Path, BotProtection(), _routesHandler(routes.Title, routes.Handler))
		}
	case utils.METHOD_POST:
		switch routes.Permission {
		case utils.PERMISSION_COOKIE:
			app.POST(routes.Path, BotProtection(), RequestHandlerCookie(), _routesHandler(routes.Title, routes.Handler))
		default:
			app.POST(routes.Path, BotProtection(), _routesHandler(routes.Title, routes.Handler))
		}
	case utils.METHOD_PATCH:
		switch routes.Permission {
		case utils.PERMISSION_COOKIE:
			app.PATCH(routes.Path, BotProtection(), RequestHandlerCookie(), _routesHandler(routes.Title, routes.Handler))
		default:
			app.PATCH(routes.Path, BotProtection(), _routesHandler(routes.Title, routes.Handler))
		}
	case utils.METHOD_PUT:
		switch routes.Permission {
		case utils.PERMISSION_COOKIE:
			app.PUT(routes.Path, BotProtection(), RequestHandlerCookie(), _routesHandler(routes.Title, routes.Handler))
		default:
			app.PUT(routes.Path, BotProtection(), _routesHandler(routes.Title, routes.Handler))
		}
	case utils.METHOD_DELETE:
		switch routes.Permission {
		case utils.PERMISSION_COOKIE:
			app.DELETE(routes.Path, BotProtection(), RequestHandlerCookie(), _routesHandler(routes.Title, routes.Handler))
		default:
			app.DELETE(routes.Path, BotProtection(), _routesHandler(routes.Title, routes.Handler))
		}
	}
}

func getPort(port []int) int {
	// TODO: Port from .env
	println("Trying to fetch .env PORT")
	if port_ := os.Getenv("PORT"); port_ != "" {
		if v, err := strconv.Atoi(port_); err == nil {
			return v
		}
	}

	fmt.Println("Trying to fetch dev input PORT")
	if len(port) > 0 {
		return port[0]
	}

	return 8000
}

// StartServer serves the embedded static assets, installs request logging and
// a custom 404 handler, and starts the HTTP server. The port is taken from PORT,
// then the first variadic argument, then 8000.
func StartServer(static fs.FS, port ...int) {
	p := getPort(port)

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	app.StaticFS("/static", http.FS(staticFS))

	app.Use(func(ctx *gin.Context) {
		// TODO: To add manual LOGS to monitor by data or file

		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		// TODO: Header setter

		ctx.Next()

		status := ctx.Writer.Status()
		// INFO: According to ChatGPT: Go’s time formatting uses a reference date instead of formatting codes (like YYYY or MM)
		fmt.Printf("[%s] [%s] %s -> %d\n", time.Now().Format("2006-01-02 15:04:05"), method, path, status)
	})

	app.NoRoute(func(ctx *gin.Context) {
		const html = `
			<!DOCTYPE html>
			<html>
				<head>
					<title>Page Not Found</title>
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
							font-size: .75rem;
							margin: 0;
						}	
					</style>
				</head>
				<body>
					<p><b style="font-size: 1.25rem;">404|</b> Endpoint Not Found</p>
				</body>
			</html>
		`
		ctx.Data(404, "text/html; charset=utf-8", []byte(html))
	})

	fmt.Printf("Running in PORT %d\n\n", p)
	app.Run(fmt.Sprintf(":%d", p))
}

func responseHeaderMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Developer", "MPOP Reverse II")
		ctx.Header("Powered-By", "Constellation")
		ctx.Next()
	}
}
