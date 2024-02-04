package handlers

import (
	"aitu-funpage/backend/internal/repository/nosql"
	"aitu-funpage/backend/internal/rest/forms"
	"aitu-funpage/backend/internal/rest/models"
	"aitu-funpage/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type NewsHandlers struct {
	MongoDb nosql.NewsRepo
}

func NewNewsHandlers(mongoDb nosql.NewsRepo) *NewsHandlers {
	return &NewsHandlers{
		MongoDb: mongoDb,
	}
}

func (h NewsHandlers) CreateNews(context *gin.Context) {
	username, exists := context.Get("username")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	usernamem, ok := username.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user Username")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error while retrieving user Username"})
		return
	}

	var newsForm forms.NewsForm

	if err := context.BindJSON(&newsForm); err != nil {
		logger.GetLogger().Error("Invalid news request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var news models.News
	news.Content = newsForm.Content
	news.Author = usernamem
	news.CreatedAt = time.Now()

	news.Tags = []models.Tag{}
	news.Comments = []models.Comment{}

	if err := h.MongoDb.CreateNews(&news); err != nil {
		logger.GetLogger().Error("Failed to create news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("News created successfully")
	context.JSON(http.StatusOK, gin.H{"message": "News created successfully", "data": news})
}

func (h NewsHandlers) GetNewsByObjectID(context *gin.Context) {
	objectIDParam := context.Query("object_id")

	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	news, err := h.MongoDb.GetNewsByObjectID(objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to get news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": news})
}

func (h NewsHandlers) GetAllNewsByAuthor(context *gin.Context) {
	authorName := context.Query("author")

	news, err := h.MongoDb.GetAllNewsByAuthor(authorName)
	if err != nil {
		logger.GetLogger().Error("Failed to get news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": news})
}

func (h NewsHandlers) GetAllNewsByTag(context *gin.Context) {
	tagName := context.Query("tag")

	tags, err := h.MongoDb.GetAllNewsByTag(tagName)
	if err != nil {
		logger.GetLogger().Error("Failed to get news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": tags})
}

func (h NewsHandlers) UpdateNewsByObjectID(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": "TODO!"})
}

func (h NewsHandlers) DeleteNewsByObjectID(context *gin.Context) {
	objectIDParam := context.Query(" ")

	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	err = h.MongoDb.DeleteNewsByObjectID(objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to delete news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}
