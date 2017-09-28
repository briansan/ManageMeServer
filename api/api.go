package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const envWWWHost = "MANAGEME_WWW_HOST"

var wwwHost string

func init() {
	wwwHost = os.Getenv(envWWWHost)
	if len(wwwHost) == 0 {
		logger.Info("env.MANAGEME_WWW_HOST not specified, defaulting to http://localhost:8889")
		wwwHost = "http://localhost:8889"
	}
}

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{wwwHost},
		AllowCredentials: true,
	}))

	// setup /api
	api := e.Group("/api")

	// setup /api/service
	svc := api.Group("/service")

	// ping pong
	svc.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// setup users
	initAuth(api)
	initUsers(api)
	initTasks(api)

	// setup the rest
	return e
}
