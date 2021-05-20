package vault

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoDriver() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
}

func getCollection(client *mongo.Client, collectionName string) (*mongo.Collection) {
	return client.Database("hefestoVault").Collection(collectionName)
}

