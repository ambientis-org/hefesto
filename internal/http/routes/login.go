package routes

import (
	"net/http"
	"os"
	"time"

	mongomodels "github.com/ambientis-org/hefesto/internal/db/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ambientis-org/hefesto/internal/db/postgres/models"
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const THREEDAYS = 72

// login receives user and password and returns token
func login(c echo.Context) error {
	requested := &models.User{}
	err := c.Bind(requested)
	if err != nil {
		return err
	}

	user := &models.User{}
	DataBase.Where("Email = ?", requested.Email).First(user)

	if requested.Password != user.Password {
		return echo.ErrUnauthorized
	}

	claims := &auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * THREEDAYS).Unix(),
		},
		Username: requested.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("API_KEY")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"mentiaAuthToken": signedToken,
		"username":        user.Username,
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

// createJournalFor Makes a new Journal on DB for user
func createJournalFor(c echo.Context) error {
	u := getUser(c.Param("username"))

	j := &mongomodels.Journal{}
	err := MongoRepo.FindOne(ctx, bson.D{primitive.E{Key: "user_id", Value: u.ID}}).Decode(&j)
	if err == mongo.ErrNoDocuments {
		j = mongomodels.NewJournal(u.ID, c.Param("username"))
		_, err := MongoRepo.InsertOne(ctx, j)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, j)
	}
	return c.String(http.StatusAlreadyReported, "El usuario ya tiene un Journal")
}

// setupLogin Add Login handlers to API
func (router *Router) setupLogin() {
	router.addGroup("/login")
	router.addGroup("/register")
	router.addGroup("/newJournal")

	loginGroup := API.groups["/login"]
	loginGroup.POST("", login)

	registerGroup := API.groups["/register"]
	registerGroup.POST("", register)

	newJournalGroup := API.groups["/newJournal"]
	newJournalGroup.POST("/:username", createJournalFor)
}
