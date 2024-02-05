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

type CommentRepo interface {
	CreateDocumentByNewsID(newsObjectID primitive.ObjectID) error
	AddCommentToNews(newsObjectID primitive.ObjectID, commentItem *models.CommentData) error
	GetAllCommentsByNewsID(newsObjectID primitive.ObjectID) (*models.Comments, error)
	GetCommentByCommentID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID) (*models.Comments, error)
	UpdateCommentsByNewsID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID, updatedCommentItem *models.CommentData) error
	DeleteCommentByNewsIDAndCommentID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID) error
}

type CommentRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		db:   db,
		coll: db.Collection(mongodb.COMMENT_COLLECTION_NAME),
	}
}

func (cr *CommentRepository) CreateDocumentByNewsID(newsObjectID primitive.ObjectID) error {
	var comment models.Comments
	comment.NewsID = newsObjectID
	comment.Comments = []models.CommentData{}

	_, err := cr.coll.InsertOne(context.TODO(), comment)
	if err != nil {
		logger.GetLogger().Fatal("Unable to create document")
		return err
	}

	return nil
}

func (cr *CommentRepository) AddCommentToNews(newsObjectID primitive.ObjectID, commentItem *models.CommentData) error {
	filter := bson.D{{"news_id", newsObjectID}}

	update := bson.D{
		{"$push", bson.D{
			{"comments", commentItem},
		}},
	}

	_, err := cr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to add the mongoDB request! Error: ", err.Error())
		return err
	}

	return nil
}

func (cr *CommentRepository) GetCommentByCommentID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID) (*models.Comments, error) {
	filter := bson.D{
		{"news_id", newsObjectID},
		{"comments", bson.M{"$elemMatch": bson.M{"_id": commentObjectID}}},
	}

	var result models.Comments
	err := cr.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}

	return &result, nil
}

func (cr *CommentRepository) GetAllCommentsByNewsID(newsObjectID primitive.ObjectID) (*models.Comments, error) {
	filter := bson.D{{"news_id", newsObjectID}}

	var result models.Comments
	err := cr.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.GetLogger().Fatal("Unable to find the mongoDB request! Error: ", err.Error())
		return nil, err
	}

	return &result, nil
}

func (cr *CommentRepository) UpdateCommentsByNewsID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID, updatedCommentItem *models.CommentData) error {
	filter := bson.D{
		{"news_id", newsObjectID},
		{"comments", bson.M{"$elemMatch": bson.M{"_id": commentObjectID}}},
	}

	update := bson.D{
		{"$set", bson.D{
			{"comments.$.content", updatedCommentItem.Content},
		}},
	}

	_, err := cr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to update the request! Error: ", err.Error())
		return err
	}

	return nil
}

func (cr *CommentRepository) DeleteCommentByNewsIDAndCommentID(newsObjectID primitive.ObjectID, commentObjectID primitive.ObjectID) error {
	filter := bson.D{
		{"news_id", newsObjectID},
	}

	update := bson.D{
		{"$pull", bson.D{
			{"comments", bson.M{"_id": commentObjectID}},
		}},
	}

	_, err := cr.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.GetLogger().Fatal("Unable to delete the request! Error: ", err.Error())
		return err
	}

	return nil
}
