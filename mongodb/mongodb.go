package mongodb

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once     sync.Once
	mongoCli *mongo.Client
)

// Connect initializes and returns a MongoDB client.
func Connect(mongoUri string) *mongo.Client {
	var err error
	once.Do(func() {
		mongoCli, err = connectToMongoDB(mongoUri)
		if err != nil {
			panic(err)
		}
	})
	return mongoCli
}

func connectToMongoDB(mongoUri string) (*mongo.Client, error) {
	mongoURI := mongoUri
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
