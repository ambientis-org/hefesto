package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ambientis-org/hefesto/internal/db/models"
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const THREEDAYS = 72

// login receives user and password and returns token
func login(c echo.Context) error {
	requestedUser := &models.User{}
	err := c.Bind(requestedUser)
	if err != nil {
		return err
	}

	user := &models.User{}
	DataBase.Where("Username = ?", requestedUser.Username).First(user)

	if requestedUser.Password != user.Password {
		return echo.ErrUnauthorized
	}

	claims := &auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * THREEDAYS).Unix(),
		},
		Username: requestedUser.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("API_KEY")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"token": signedToken,
	})
}

// register POST method User handler
func register(c echo.Context) error {
	user := &models.User{}
	err := c.Bind(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	DataBase.Create(user)

	return c.JSON(http.StatusCreated, user)
}

// setupLogin Add Login handlers to API
func (router *Router) setupLogin() {
	router.addGroup("/login")
	router.addGroup("/register")

	loginGroup := API.groups["/login"]
	loginGroup.POST("", login)

	registerGroup := API.groups["/register"]
	registerGroup.POST("", register)
}
