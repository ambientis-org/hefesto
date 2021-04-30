package routes

import (
	"net/http"

	"github.com/ambientis-org/hefesto/internal/db/models"
	"github.com/labstack/echo/v4"
)

// TODO add validator

// getAllUsers GET method for all users handler
func getAllUsers(c echo.Context) error {
	var users []models.User

	DataBase.Find(&users)
	return c.JSON(http.StatusOK, users)
}

// getUserByID GET method User handler
func getUserByID(c echo.Context) error {
	u := &models.User{}

	DataBase.First(u, c.Param("id"))
	if u.ID == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, u)
}

// getUserByUsername
func getUserByUsername(c echo.Context) error {
	u := &models.User{}
	DataBase.First(u, c.Param("username"))
	if u.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
	}

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

// createUser POST method User handler
func createUser(c echo.Context) error {
	user := &models.User{}
	err := c.Bind(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	DataBase.Create(user)

	return c.JSON(http.StatusCreated, user)
}

// SetupUsers Add User handlers to API
func (router *Router) setupUsers() {

	// User endpoint, protecting it
	router.addGroup("/users")
	group := API.groups["/users"]

	group.GET("", getAllUsers)
	group.GET("/:id", getUserByID)
	group.GET("/:username", getUserByUsername)
	group.POST("", createUser)
	group.PATCH("/:id", updateUser)
	group.DELETE("/:id", deleteUser)
}
