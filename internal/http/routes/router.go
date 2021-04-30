package routes

import (
	"net/http"
	"os"

	"github.com/ambientis-org/hefesto/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	Server *echo.Echo
	groups map[string]*echo.Group
}

func newRouter() *Router {
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	return &Router{Server: server, groups: make(map[string]*echo.Group)}
}

func (router *Router) addGroup(prefix string) {
	router.groups[prefix] = router.Server.Group("/api/v2" + prefix)
}

// API Package element
var API = newRouter()

// DataBase Package element
var DataBase, _ = db.New(os.Getenv("POSTGRES_DSN"))

func healthcheck(c echo.Context) error {
	if DataBase == nil {
		return http.ErrAbortHandler
	}
	return c.String(http.StatusOK, "hefesto funciona correctamente")
}

// GetRouter setup handlers and exports the Router instance
func GetRouter() *Router {
	API.Server.GET("/healthcheck", healthcheck)
	API.setupUsers()

	return API
}
