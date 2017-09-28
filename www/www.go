package www

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	envAssetsDir = "MANAGEME_ASSETS_DIR"
	envAPIHost   = "MANAGEME_API_HOST"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// Allow assets filepath to be configurable
	assetsDir := os.Getenv(envAssetsDir)
	if len(assetsDir) == 0 {
		assetsDir = "www/client"
	}

	// Get the API host env var
	apiHost := os.Getenv(envAPIHost)
	if len(apiHost) == 0 {
		apiHost = "http://localhost:8888"
	}

	// Set up route to report the api url
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, apiHost)
	})

	// Serve that folder and default index.html
	e.Use(middleware.Static(assetsDir))
	return e
}
