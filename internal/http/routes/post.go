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

	update := bson.D{bson.E{
		Key: "$set",
		Value: bson.D{
			bson.E{
				Key:   "posts",
				Value: tmpPosts,
			},
		},
	}}

	err = MongoMoodRepo.FindOneAndUpdate(ctx, filter, update, opts).Decode(&bson.M{})
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, newPost)
}

func getUserPosts(c echo.Context) error {
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

		var posts []models.Post

		for _, v := range j.Posts {
			postTimestamp := v.CreatedAt.UTC().Unix()
			if from.UTC().Unix() < postTimestamp && postTimestamp <= to.UTC().Unix() {
				posts = append(posts, v)
			}
		}
		j.Posts = posts
	}

	return c.JSON(http.StatusOK, j)
}

// setupPosts Add User handlers to API
func (router *Router) setupPosts() {

	// User endpoint, protecting it
	router.addGroup("/posts")
	group := API.groups["/posts"]
	config := middleware.JWTConfig{
		Claims:     &auth.CustomClaims{},
		SigningKey: []byte(os.Getenv("API_KEY")),
	}

	group.Use(middleware.JWTWithConfig(config))

	group.POST("/:username", createPost)
	group.GET("/:username", getUserPosts)
}
