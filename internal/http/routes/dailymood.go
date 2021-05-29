package routes

import (
	"context"
	"net/http"
	"os"

	postgresmodels "github.com/ambientis-org/hefesto/internal/db/models"
	mongomodels "github.com/ambientis-org/hefesto/internal/db/vault/models"
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Move this function
func getUser(username string) *postgresmodels.User {
	u := &postgresmodels.User{}
	DataBase.Where("username = ?", username).First(u)
	return u
}

var ctx = context.TODO()

// On Register create a new Journal for user

// createJournalFor Makes a new Journal on DB for user
func createJournalFor(c echo.Context) error {
	u := getUser(c.Param("username"))

	j := &mongomodels.Journal{}
	err := MongoRepo.FindOne(ctx, bson.D{{"user_id", u.ID}}).Decode(&j)
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

// TODO: Limit it to only one per day
// addMoodForToday add a new value to Moods array for user
func addMoodForToday(c echo.Context) error {
	// Processing requesr
	u := getUser(c.Param("username"))
	m := &mongomodels.Mood{}
	err := c.Bind(m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	// Querying from MongoDB
	filter := bson.D{{"user_id", u.ID}}
	j := &mongomodels.Journal{}
	err = MongoRepo.FindOne(ctx, filter).Decode(j)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Adding mood to user's journal
	opts := options.FindOneAndUpdate().SetUpsert(true)
	tmpMoods := append(j.Moods, m)
	update := bson.D{bson.E{
		Key: "$set",
		Value: bson.D{
			bson.E{
				Key:   "moods",
				Value: tmpMoods,
			},
		},
	}}

	if m.CreatedAt.Date() == j.Moods[len(j.Moods)-1].Date() {
		return c.String(http.StatusBadRequest, "Ys has regustrado un mood para hoy")
	} else {
		err = MongoRepo.FindOneAndUpdate(ctx, filter, update, opts).Decode(&bson.M{})
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Mood del día de hoy registrado con éxito")
}

// TODO: Make GET using time periods

// getUserMoods return moods
func getUserMoods(c echo.Context) error {
	u := getUser(c.Param("username"))

	filter := bson.D{{"user_id", u.ID}}
	j := &mongomodels.Journal{}
	err := MongoRepo.FindOne(ctx, filter).Decode(j)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, j)
}

// setupMoods Add User handlers to API
func (router *Router) setupMoods() {

	// User endpoint, protecting it
	router.addGroup("/mood")
	group := API.groups["/mood"]
	config := middleware.JWTConfig{
		Claims:     &auth.CustomClaims{},
		SigningKey: []byte(os.Getenv("API_KEY")),
	}

	group.Use(middleware.JWTWithConfig(config))

	group.POST("/:username", createJournalFor)
	group.POST("/:username/add", addMoodForToday)
	group.GET("/:username", getUserMoods)

}
