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
	NewsRepo nosql.NewsRepo
	TagRepo  nosql.TagRepo
}

func NewNewsHandlers(newsRepo nosql.NewsRepo, tagRepo nosql.TagRepo) *NewsHandlers {
	return &NewsHandlers{
		NewsRepo: newsRepo,
		TagRepo:  tagRepo,
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
	news.Title = newsForm.Title
	news.Content = newsForm.Content
	news.Author = usernamem
	news.CreatedAt = time.Now()

	insertedId, err := h.NewsRepo.CreateNews(&news)
	if err != nil {
		logger.GetLogger().Error("Failed to create news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TagRepo.CreateDocumentByNewsID(insertedId); err != nil {
		logger.GetLogger().Error("Failed to create tag document:", err)
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

	news, err := h.NewsRepo.GetNewsByObjectID(objectID)
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

	news, err := h.NewsRepo.GetAllNewsByAuthor(authorName)
	if err != nil {
		logger.GetLogger().Error("Failed to get news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": news})
}

func (h NewsHandlers) UpdateNewsByObjectID(context *gin.Context) {
	objectIDParam := context.Query("object_id")
	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

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

	news, err := h.NewsRepo.GetNewsByObjectID(objectID)
	if news.Author != usernamem {
		logger.GetLogger().Error("The user can't delete the news item. The user isn't the owner of the news")
		context.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You can't delete the news. You're not the owner of the news"})
		return
	}

	var newsForm forms.NewsForm

	if err := context.BindJSON(&newsForm); err != nil {
		logger.GetLogger().Error("Invalid news request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var updateNews models.News
	news.Title = newsForm.Title
	news.Content = newsForm.Content

	if err := h.NewsRepo.UpdateNewsByObjectID(objectID, &updateNews); err != nil {
		logger.GetLogger().Error("Failed to update news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("News updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "News updated successfully", "data": news})
}

func (h NewsHandlers) DeleteNewsByObjectID(context *gin.Context) {
	objectIDParam := context.Query("object_id")
	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

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

	news, err := h.NewsRepo.GetNewsByObjectID(objectID)
	if news.Author != usernamem {
		logger.GetLogger().Error("The user can't delete the news item. The user isn't the owner of the news")
		context.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You can't delete the news. You're not the owner of the news"})
		return
	}

	err = h.NewsRepo.DeleteNewsByObjectID(objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to delete news:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Delete News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}
