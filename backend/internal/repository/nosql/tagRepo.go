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

type TagRepo interface {
	CreateDocumentByNewsID(newsObjectID primitive.ObjectID) error
	AddTagToNews(newsObjectID primitive.ObjectID, tagItem *models.TagData) error
	GetTagsByNewsID(newsObjectID primitive.ObjectID) (*models.Tags, error)
	UpdateTagByNewsID(newsObjectID primitive.ObjectID, tagName string, updatedTagItem *models.TagData) error
	DeleteTagByNewsID(newsObjectID primitive.ObjectID, tagName string) error
}

type TagRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewTagRepository(db *mongo.Database) *TagRepository {
	return &TagRepository{
		db:   db,
		coll: db.Collection(mongodb.TAGS_COLLECTION_NAME),
	}
}

func (tr *TagRepository) CreateDocumentByNewsID(newsObjectID primitive.ObjectID) error {
	var tags models.Tags
	tags.NewsID = newsObjectID
	tags.Tags = []models.TagData{}

	_, err := tr.coll.InsertOne(context.TODO(), tags)
	if err != nil {
		logger.GetLogger().Fatal("Unable to create document")
		return err
	}

	return nil
}

func (tr *TagRepository) AddTagToNews(newsObjectID primitive.ObjectID, tagItem *models.TagData) error {
	filter := bson.D{{"news_id", newsObjectID}}

	update := bson.D{
		{"$push", bson.D{
			{"tags", tagItem},
		}},
	}

	_, err := tr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to add the mongoDB request! Error: ", err.Error())
		return err
	}

	return nil
}

func (tr *TagRepository) GetTagsByNewsID(newsObjectID primitive.ObjectID) (*models.Tags, error) {
	filter := bson.D{{"news_id", newsObjectID}}

	var result models.Tags
	err := tr.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}

	return &result, nil
}

func (tr *TagRepository) UpdateTagByNewsID(newsObjectID primitive.ObjectID, tagName string, updatedTagItem *models.TagData) error {
	filter := bson.D{
		{"news_id", newsObjectID},
		{"tags", bson.M{"$elemMatch": bson.M{"tag": tagName}}},
	}

	update := bson.D{
		{"$set", bson.D{
			{"tags.$.tag", updatedTagItem.TagName},
			{"tags.$.color", updatedTagItem.Color},
			{"tags.$.created_at", updatedTagItem.CreatedAt},
		}},
	}

	_, err := tr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to update the request! Error: ", err.Error())
		return err
	}

	return nil
}

func (tr *TagRepository) DeleteTagByNewsID(newsObjectID primitive.ObjectID, tagName string) error {
	filter := bson.D{
		{"news_id", newsObjectID},
	}

	update := bson.D{
		{"$pull", bson.D{
			{"tags", bson.M{"tag": tagName}},
		}},
	}

	_, err := tr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to delete the request! Error: ", err.Error())
		return err
	}

	return nil
}
