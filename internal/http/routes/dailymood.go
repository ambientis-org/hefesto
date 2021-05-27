package routes

import (
	"context"
	"github.com/ambientis-org/hefesto/internal/db/models"
	"github.com/ambientis-org/hefesto/internal/http/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"strconv"
	"time"
)

// TODO: Move "Models" to other file
type Mood struct {
	ID 		  primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Value 	  int 		`bson:"value"`
}

func NewMood(value int) Mood {
	return Mood{
		ID: primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Value: value,
	}
}

type Journal struct {
	ID 		  primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	UserID	  uint		`bson:"user_id"`
	Username  string	`bson:"username"`
	Moods	  []Mood	`bson:"moods"`
}

func NewJournal(userID uint, username string) *Journal {
	return &Journal{
		ID: primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userID,
		Username: username,
		Moods: []Mood{},
	}
}

// TODO: Move this function
func getUser(username string) *models.User {
	u := &models.User{}
	DataBase.Where("username = ?", username).First(u)
	return u
}

var ctx = context.TODO()

// On Register create a new Journal for user

// createJournalFor Makes a new Journal on DB for user
func createJournalFor(c echo.Context) error {
	u := getUser(c.Param("username"))

	j := &Journal{}
	err := MongoRepo.FindOne(ctx, bson.D{{"user_id", u.ID}}).Decode(&j)
	if err == mongo.ErrNoDocuments {
		j = NewJournal(u.ID, c.Param("username"))
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
	u := getUser(c.Param("username"))

	filter := bson.D{{"user_id", u.ID}}
	j := &Journal{}
	err := MongoRepo.FindOne(ctx, filter).Decode(j)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	newMoodValue, _ := strconv.Atoi(c.Param("moodValue"))
	m := NewMood(newMoodValue)
	opts := options.FindOneAndUpdate().SetUpsert(true)

	tmpMoods := append(j.Moods, m)
	update := bson.D{{"$set", bson.D{{"moods", tmpMoods}}}}
	err = MongoRepo.FindOneAndUpdate(ctx, filter, update, opts).Decode(&bson.M{})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Mood del día de hoy registrado con éxito")
}

// TODO: Make GET using time periods

// getUserMoods return moods
func getUserMoods(c echo.Context) error {
	u := getUser(c.Param("username"))

	filter := bson.D{{"user_id", u.ID}}
	j := &Journal{}
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
		Claims: &auth.CustomClaims{},
		SigningKey: []byte(os.Getenv("API_KEY")),
	}

	group.Use(middleware.JWTWithConfig(config))

	group.POST("/:username", createJournalFor)
	group.POST("/addToday/:username/:moodValue", addMoodForToday)
	group.GET("/:username", getUserMoods)

}
