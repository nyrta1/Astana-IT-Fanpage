package mongodb

import (
	"aitu-funpage/backend/internal/config"
	"aitu-funpage/backend/pkg/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

var (
	NEWS_COLLECTION_NAME    string = "news"
	COMMENT_COLLECTION_NAME string = "comments"
	TAGS_COLLECTION_NAME    string = "tags"
)

func initializeMongoDB(mongoDbConfig config.MongoDbConfig) error {
	uri := mongoDbConfig.Addr
	if uri == "" {
		logger.GetLogger().Fatal("You must set your 'MONGODB_URI' environment variable.")
		return fmt.Errorf("you must set your 'MONGODB_URI' environment variable")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.GetLogger().Fatal("Unable to connect to MongoDB. Error: ", err.Error())
		return err
	}

	db = client.Database(mongoDbConfig.DatabaseName)

	return nil
}

func GetMongoDbInstance(mongoDbConfig config.MongoDbConfig) (*mongo.Database, error) {
	db = nil
	var errGetMongoDB error
	if db == nil {
		if err := initializeMongoDB(mongoDbConfig); err != nil {
			errGetMongoDB = err
		}
	}
	return db, errGetMongoDB
}
