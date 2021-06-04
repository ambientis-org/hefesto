package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ambientis-org/hefesto/internal/db/mongo/models"
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// addMoodForToday add a new value to Moods array for user
func addMoodForToday(c echo.Context) error {
	// Processing request
	u := GetUser(c.Param("username"))
	requestBody := &models.Mood{}
	err := c.Bind(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	m := models.NewMood(requestBody.Value)

	// Querying from MongoDB
	filter := bson.D{primitive.E{Key: "user_id", Value: u.ID}}
	j := &models.Journal{}
	err = MongoMoodRepo.FindOne(ctx, filter).Decode(j)
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

	lenMoods := len(j.Moods)
	newMoodWeekday := m.CreatedAt.UTC().Weekday()

	if len(j.Moods) == 0 || j.Moods[lenMoods-1].CreatedAt.Weekday() != newMoodWeekday {
		err = MongoMoodRepo.FindOneAndUpdate(ctx, filter, update, opts).Decode(&bson.M{})
	} else {
		return c.String(http.StatusAlreadyReported, "Ya has registrado un mood para hoy")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, m)
}

// getUserMoods return moods
func getUserMoods(c echo.Context) error {
	u := GetUser(c.Param("username"))
	j := &models.Journal{}

	filter := bson.D{primitive.E{Key: "user_id", Value: u.ID}}
	err := MongoMoodRepo.FindOne(ctx, filter).Decode(j)

	if c.QueryParam("from") == "" && c.QueryParam("to") == "" {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		layout := "2006-01-02T15:04:05.000Z"
		from, _ := time.Parse(layout, c.QueryParam("from"))
		to, _ := time.Parse(layout, c.QueryParam("to"))

		var moods []models.Mood

		for _, v := range j.Moods {
			moodTimestamp := v.CreatedAt.UTC().Unix()
			if from.UTC().Unix() < moodTimestamp && moodTimestamp <= to.UTC().Unix() {
				moods = append(moods, v)
			}
		}
		j.Moods = moods
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

	group.POST("/:username/add", addMoodForToday)
	group.GET("/:username", getUserMoods)
}
