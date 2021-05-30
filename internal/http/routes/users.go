package routes

import (
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"

	"github.com/ambientis-org/hefesto/internal/db/postgres/models"
	"github.com/labstack/echo/v4"
)

// TODO add validator

// getAllUsers GET method for all users handler
func getAllUsers(c echo.Context) error {
	var users []models.User

	DataBase.Find(&users)
	return c.JSON(http.StatusOK, users)
}


// getUserByUsername
func getUserByUsername(c echo.Context) error {
	u := &models.User{}
	DataBase.Where("username = ?", c.Param("username")).First(u)

	return c.JSON(http.StatusOK, u)
}

// deleteUser DELETE method User handler
func deleteUser(c echo.Context) error {
	user := &models.User{}
	DataBase.Delete(user, c.Param("id"))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, user)
}

// updateUser PATCH method User handler
func updateUser(c echo.Context) error {
	newModel := &models.User{}
	user := &models.User{}

	err := c.Bind(newModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	DataBase.First(user, c.Param("id"))
	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	DataBase.Model(user).Updates(newModel)
	return c.JSON(http.StatusOK, user)
}

// SetupUsers Add User handlers to API
func (router *Router) setupUsers() {

	// User endpoint, protecting it
	router.addGroup("/users")
	group := API.groups["/users"]
	config := middleware.JWTConfig{
		Claims: &auth.CustomClaims{},
		SigningKey: []byte(os.Getenv("API_KEY")),
	}

	group.Use(middleware.JWTWithConfig(config))

	group.GET("", getAllUsers)
	group.GET("/:username", getUserByUsername)
	group.PATCH("/:id", updateUser)
	group.DELETE("/:id", deleteUser)
}
