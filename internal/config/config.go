package config

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"solanaindexer/internal/logger"
	"solanaindexer/internal/utils"
	"time"
)

var Client *mongo.Client
var envVariables = utils.LoadEnvVariables()

func getDBUri() (string, error) {
	if envVariables.DbUsername == "" || envVariables.DbPassword == "" || envVariables.DbName == "" {
		return "", errors.New("missing required environment variables")
	}
	return "mongodb://" + envVariables.DbUsername + ":" + envVariables.DbPassword + "@mongodb:27017/" + envVariables.DbName + "?authSource=admin", nil
}

func ConnectDB() error {
	uri, err := getDBUri()
	if err != nil {
		return err
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	logger.Infof("Connected to MongoDB!")
	Client = client
	return nil
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
	if Client == nil {
		return nil, errors.New("database client is not initialized")
	}
	return Client.Database(envVariables.DbName).Collection(collectionName), nil
}
