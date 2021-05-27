package routes

import (
	"net/http"
	"os"

	"github.com/ambientis-org/hefesto/internal/db"
	"github.com/ambientis-org/hefesto/internal/db/vault"
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

	customCORS := middleware.CORSConfig{
		ExposeHeaders: []string{echo.HeaderSetCookie},
		AllowHeaders: []string{echo.HeaderSetCookie, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}
	server.Use(middleware.CORSWithConfig(customCORS))

	return &Router{Server: server, groups: make(map[string]*echo.Group)}
}

func (router *Router) addGroup(prefix string) {
	router.groups[prefix] = router.Server.Group("/api/v2" + prefix)
}

// API Package element
var API = newRouter()

// DataBase Package element
var DataBase, _ = db.New(os.Getenv("POSTGRES_DSN"))

// MongoDB Package element
var MongoRepo = vault.New(os.Getenv("MONGODB_COLLECTION"))

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
	API.setupLogin()
	API.setupMoods()

	return API
}
