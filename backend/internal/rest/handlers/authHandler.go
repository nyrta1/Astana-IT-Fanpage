package handlers

import (
	"aitu-funpage/backend/internal/config"
	"aitu-funpage/backend/internal/repository"
	"aitu-funpage/backend/internal/rest/forms"
	"aitu-funpage/backend/internal/rest/models"
	"aitu-funpage/backend/pkg/logger"
	"aitu-funpage/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AuthHandlers struct {
	Repo         repository.UserRepo
	UserTypeRepo repository.UserTypeRepo
	RedisConfig  config.RedisConfig
}

func NewAuthHandlers(repo repository.UserRepo, userTypeRepo repository.UserTypeRepo, redisConfig config.RedisConfig) *AuthHandlers {
	return &AuthHandlers{
		Repo:         repo,
		UserTypeRepo: userTypeRepo,
		RedisConfig:  redisConfig,
	}
}

func (h AuthHandlers) Register(context *gin.Context) {
	logger.GetLogger().Info("Received registration request")

	var user models.User

	if err := context.BindJSON(&user); err != nil {
		logger.GetLogger().Error("Invalid registration request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.Repo.GetUserByUsername(user.Username)
	if err == nil {
		logger.GetLogger().Error("Account already registered for username:", user.Username)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}

	userTypeID, err := h.UserTypeRepo.GetIDByTypeName(models.USER_ROLE)
	if err != nil {
		logger.GetLogger().Error("Account cannot get userTypeID:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Account cannot get userTypeID"})
		return
	}
	user.UserType = userTypeID

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.GetLogger().Error("Unable to hash the password")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash the password"})
		return
	}
	user.Password = hashedPassword

	userTypeName, err := h.UserTypeRepo.GetTypeByID(user.UserType)
	if err != nil {
		logger.GetLogger().Error("Unable to find the user type")
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find the user type"})
		return
	}

	if err := h.Repo.CreateUser(&user); err != nil {
		logger.GetLogger().Error("Failed to create user:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signedToken, _ := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Username, userTypeName)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/app",
		Expires:  time.Now().Add(time.Minute * 20),
		HttpOnly: true,
	}

	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User registered successfully")
	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "data": signedToken})
}

func (h AuthHandlers) Login(context *gin.Context) {
	var data forms.LoginForm
	if err := context.BindJSON(&data); err != nil {
		logger.GetLogger().Error("Invalid login request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.Repo.GetUserByUsername(data.Username)
	if err != nil {
		logger.GetLogger().Error("Failed to get user by username:", err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		logger.GetLogger().Error("Authentication failed for username:", user.Username)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userType, err := h.UserTypeRepo.GetTypeByID(user.UserType)
	if err != nil {
		logger.GetLogger().Error("Account cannot get userTypeID:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Account cannot get userTypeID"})
		return
	}

	token, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Username, userType)
	if err != nil {
		logger.GetLogger().Error("Failed to create token:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/app",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User login successful")
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (h AuthHandlers) Logout(context *gin.Context) {
	logger.GetLogger().Info("User logout")

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/app",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User logged out successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User logged out successfully", "data": nil})
}

func (h AuthHandlers) DeleteAccount(context *gin.Context) {
	username, exists := context.Get("username")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		context.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	usernamem, ok := username.(string)
	if !ok {
		logger.GetLogger().Error("Error while retrieving user ID")
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error while retrieving user ID"})
		return
	}

	user, err := h.Repo.GetUserByUsername(usernamem)
	if err != nil {
		logger.GetLogger().Error("User not found:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	err = h.Repo.DeleteUser(user.ID)
	if err != nil {
		logger.GetLogger().Error("Failed to delete user:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete user"})
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/app",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)

	logger.GetLogger().Info("User account deleted successfully")
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "User account deleted successfully"})
}
