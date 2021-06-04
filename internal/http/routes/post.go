package routes

import (
	"net/http"

	"github.com/ambientis-org/hefesto/internal/db/mongo/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createPost(c echo.Context) error {
	u := GetUser(c.Param("username"))
	requestBody := &models.Post{}

	err := c.Bind(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	newPost := models.NewPost(requestBody.Content)
	filter := bson.D{primitive.E{Key: "user_id", Value: u.ID}}
	j := &models.Journal{}
	err = MongoMoodRepo.FindOne(ctx, filter).Decode(j)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	tmpPosts := append(j.Posts, newPost)

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
}
