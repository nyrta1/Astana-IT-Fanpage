package nosql

import (
	"aitu-funpage/backend/internal/db/mongodb"
	"aitu-funpage/backend/internal/rest/models"
	"aitu-funpage/backend/pkg/logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsRepo interface {
	CreateNews(newNews *models.News) error
	GetNewsByObjectID(objectId string) (*models.News, error)
	GetAllNewsByAuthor(author string) ([]*models.News, error)
	GetAllNewsByTag(tagName string) ([]*models.News, error)
	UpdateNewsByObjectID(objectId string, updateNews *models.News) error
	DeleteNewsByObjectID(objectId string) error
}

type NewsRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewNewsRepository(db *mongo.Database) *NewsRepository {
	return &NewsRepository{
		db:   db,
		coll: db.Collection(mongodb.NEWS_COLLECTION_NAME),
	}
}

func (nr *NewsRepository) CreateNews(newNews *models.News) error {
	_, err := nr.coll.InsertOne(context.TODO(), newNews)
	if err != nil {
		logger.GetLogger().Fatal("Error to save the news: ", err.Error())
		return err
	}

	return nil
}

func (nr *NewsRepository) GetNewsByObjectID(objectId string) (*models.News, error) {
	id, _ := primitive.ObjectIDFromHex(objectId)
	filter := bson.D{{"_id", id}}

	var result models.News
	err := nr.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}

	return &result, nil
}

func (nr *NewsRepository) GetAllNewsByAuthor(author string) ([]*models.News, error) {
	filter := bson.D{{"author", author}}

	cursor, err := nr.coll.Find(context.TODO(), filter)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []*models.News
	for cursor.Next(context.TODO()) {
		var result models.News
		if err := cursor.Decode(&result); err != nil {
			logger.GetLogger().Fatal("Error decoding result: ", err.Error())
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

func (nr *NewsRepository) GetAllNewsByTag(tagName string) ([]*models.News, error) {
	filter := bson.D{
		{"tags", bson.M{"$elemMatch": bson.M{"tag": tagName}}}}

	cursor, err := nr.coll.Find(context.TODO(), filter)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []*models.News
	for cursor.Next(context.TODO()) {
		var result models.News
		if err := cursor.Decode(&result); err != nil {
			logger.GetLogger().Fatal("Error decoding result: ", err.Error())
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

func (nr *NewsRepository) UpdateNewsByObjectID(objectId string, updateNews *models.News) error {
	id, _ := primitive.ObjectIDFromHex(objectId)
	filter := bson.D{{"_id", id}}

	update := bson.D{{"$set", updateNews}}

	_, err := nr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Error unable to update the news: ", err.Error())
		return err
	}

	return nil
}

func (nr *NewsRepository) DeleteNewsByObjectID(objectId string) error {
	id, _ := primitive.ObjectIDFromHex(objectId)
	filter := bson.D{{"_id", id}}

	_, err := nr.coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.GetLogger().Fatal("Error unable to update the news: ", err.Error())
		return err
	}

	return nil
}
