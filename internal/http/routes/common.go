package routes

import (
	"context"

	postgresmodels "github.com/ambientis-org/hefesto/internal/db/postgres/models"
)

func GetUser(username string) *postgresmodels.User {
	u := &postgresmodels.User{}
	DataBase.Where("username = ?", username).First(u)
	return u
}

var ctx = context.TODO()
