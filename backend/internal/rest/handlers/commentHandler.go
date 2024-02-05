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

type CommentHandlers struct {
	CommentRepo nosql.CommentRepo
}

func NewCommentHandler(commentRepo nosql.CommentRepo) *CommentHandlers {
	return &CommentHandlers{
		CommentRepo: commentRepo,
	}
}

func (h CommentHandlers) AddCommentToNews(context *gin.Context) {
	objectIDParam := context.Query("news_id")
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

	var commentForm forms.CommentForm
	if err := context.BindJSON(&commentForm); err != nil {
		logger.GetLogger().Error("Invalid comment data request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var commentData models.CommentData
	commentData.CommentDataID = primitive.NewObjectID()
	commentData.Content = commentForm.Content
	commentData.Username = usernamem
	commentData.CreatedAt = time.Now()

	if err := h.CommentRepo.AddCommentToNews(objectID, &commentData); err != nil {
		logger.GetLogger().Error("Error to save the comment data")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error to save the comment data"})
		return
	}

	logger.GetLogger().Info("Comment added successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Tag added successfully", "data": commentData})
}

func (h CommentHandlers) GetCommentsByNewsID(context *gin.Context) {
	objectIDParam := context.Query("news_id")

	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	comments, err := h.CommentRepo.GetAllCommentsByNewsID(objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to get news comments:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.GetLogger().Info("Get News Comments successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": comments})
}

func (h CommentHandlers) UpdateCommentsByNewsID(context *gin.Context) {
	objectIDParam := context.Query("news_id")
	commentIDParam := context.Query("comment_id")
	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	commentID, err := primitive.ObjectIDFromHex(commentIDParam)
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

	comments, err := h.CommentRepo.GetCommentByCommentID(objectID, commentID)
	if err != nil {
		logger.GetLogger().Error("Unable to find the comment")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to find the comment"})
		return
	}

	commentData := comments.Comments[0]

	if commentData.Username != usernamem {
		logger.GetLogger().Error("The user can't delete the news item. The user isn't the owner of the news")
		context.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You can't delete the news. You're not the owner of the news"})
		return
	}

	var commentForm forms.CommentForm
	if err := context.BindJSON(&commentForm); err != nil {
		logger.GetLogger().Error("Invalid comment data request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	commentData.Content = commentForm.Content

	if err := h.CommentRepo.UpdateCommentsByNewsID(objectID, commentID, &commentData); err != nil {
		logger.GetLogger().Error("Unable to update the comment")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update the comment"})
		return
	}

	logger.GetLogger().Info("Update News Comments successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "data": commentData})
}

func (h CommentHandlers) DeleteCommentByNewsID(context *gin.Context) {
	objectIDParam := context.Query("news_id")
	commentIDParam := context.Query("comment_id")
	objectID, err := primitive.ObjectIDFromHex(objectIDParam)
	if err != nil {
		logger.GetLogger().Error("Invalid object_id parameter:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object_id parameter"})
		return
	}

	commentID, err := primitive.ObjectIDFromHex(commentIDParam)
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

	comments, err := h.CommentRepo.GetCommentByCommentID(objectID, commentID)
	if err != nil {
		logger.GetLogger().Error("Unable to find the comment")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to find the comment"})
		return
	}

	commentData := comments.Comments[0]

	if commentData.Username != usernamem {
		logger.GetLogger().Error("The user can't delete the news item. The user isn't the owner of the news")
		context.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You can't delete the news. You're not the owner of the news"})
		return
	}

	if err := h.CommentRepo.DeleteCommentByNewsIDAndCommentID(objectID, commentID); err != nil {
		logger.GetLogger().Error("Unable to delete the comment")
		context.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Unable to delete the comment"})
		return
	}

	logger.GetLogger().Info("Delete News Comments successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}
