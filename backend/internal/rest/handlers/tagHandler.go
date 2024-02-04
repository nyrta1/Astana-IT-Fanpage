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

type TagHandlers struct {
	TagRepo  nosql.TagRepo
	NewsRepo nosql.NewsRepo
}

func NewTagHandlers(tagRepo nosql.TagRepo, newsRepo nosql.NewsRepo) *TagHandlers {
	return &TagHandlers{
		TagRepo:  tagRepo,
		NewsRepo: newsRepo,
	}
}

func (h TagHandlers) AddTagToNewsByNewsID(context *gin.Context) {
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

	var tagForm forms.TagForm
	if err := context.BindJSON(&tagForm); err != nil {
		logger.GetLogger().Error("Invalid tag data request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var tagData models.TagData
	tagData.TagName = tagForm.Tag
	tagData.Color = tagForm.Color
	tagData.CreatedAt = time.Now()

	err = h.TagRepo.AddTagToNews(objectID, &tagData)
	if err != nil {
		logger.GetLogger().Error("Error to save the tag data")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error to save the tag data"})
		return
	}

	logger.GetLogger().Info("Tag added successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Tag added successfully", "data": tagData})
}

func (h TagHandlers) GetTagsByNewsID(context *gin.Context) {
	objectIDParam := context.Query("object_id")

	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	tags, err := h.TagRepo.GetTagsByNewsID(objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to get news tags:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News Tags successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": tags})
}

func (h TagHandlers) UpdateTagByNewsID(context *gin.Context) {
	objectIDParam := context.Query("object_id")
	updateTagName := context.Query("tagName")
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

	var tagForm forms.TagForm

	if err := context.BindJSON(&tagForm); err != nil {
		logger.GetLogger().Error("Invalid news request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var tagItem models.TagData
	tagItem.TagName = tagForm.Tag
	tagItem.Color = tagForm.Color

	err = h.TagRepo.UpdateTagByNewsID(objectID, updateTagName, &tagItem)
	if err != nil {
		logger.GetLogger().Error("Unable to update tag data")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update tag data"})
		return
	}

	logger.GetLogger().Info("News tag item updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "News tag item updated successfully", "data": tagItem})
}

func (h TagHandlers) DeleteTagByNewsID(context *gin.Context) {
	objectIDParam := context.Query("object_id")
	deleteTagName := context.Query("tagName")

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

	err = h.TagRepo.DeleteTagByNewsID(objectID, deleteTagName)
	if err != nil {
		logger.GetLogger().Error("Unable to delete the tag item")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to delete the tag item"})
		return
	}

	logger.GetLogger().Info("Delete News successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}
