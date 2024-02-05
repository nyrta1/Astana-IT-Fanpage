package main

import (
	"aitu-funpage/backend/internal/config"
	"aitu-funpage/backend/internal/db"
	"aitu-funpage/backend/internal/db/mongodb"
	"aitu-funpage/backend/internal/repository/nosql"
	"aitu-funpage/backend/internal/repository/sql"
	"aitu-funpage/backend/internal/rest/handlers"
	"aitu-funpage/backend/internal/routers"
	"aitu-funpage/backend/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	if err := godotenv.Load("backend.env"); err != nil {
		log.Fatalf("Error loading backend.env file: %s", err)
	}
}

func initializeDB() config.Database {
	dbConfig := config.Database{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Sslmode:  os.Getenv("DB_SSLMODE"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	return dbConfig
}

func initializeRedis() config.RedisConfig {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to int: %s", err)
	}

	redisConfig := config.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}

	return redisConfig
}

func initializeMongoDB() config.MongoDbConfig {
	mongoDbConfig := config.MongoDbConfig{
		Addr:         os.Getenv("MONGODB_URI"),
		DatabaseName: os.Getenv("MONGODB_DATABASE_NAME"),
	}

	return mongoDbConfig
}

var appConfig config.App

func main() {
	logger.InitLogger()

	appConfig = config.App{
		PORT:          os.Getenv("APP_PORT"),
		DB:            initializeDB(),
		Redis:         initializeRedis(),
		MongoDbConfig: initializeMongoDB(),
	}

	dbInstance, err := db.GetDBInstance(appConfig.DB)
	if err != nil {
		logger.GetLogger().Fatal("Error initializing DB:", err)
	}

	mongoDbInstance, err := mongodb.GetMongoDbInstance(appConfig.MongoDbConfig)
	if err != nil {
		logger.GetLogger().Fatal("Error initializing MongoDB:", err)
	}

	userRepo := sql.NewUserRepository(dbInstance)
	userTypeRepo := sql.NewUserTypeRepository(dbInstance)
	newsRepo := nosql.NewNewsRepository(mongoDbInstance)
	tagRepo := nosql.NewTagRepository(mongoDbInstance)
	commentRepo := nosql.NewCommentRepository(mongoDbInstance)

	authHandlers := handlers.NewAuthHandlers(userRepo, userTypeRepo, appConfig.Redis)
	newsHandlers := handlers.NewNewsHandlers(newsRepo, tagRepo, commentRepo)
	tagHandlers := handlers.NewTagHandlers(tagRepo, newsRepo)
	commentHandlers := handlers.NewCommentHandler(commentRepo)

	r := gin.Default()

	router := routers.NewRouters(*authHandlers, *newsHandlers, *tagHandlers, *commentHandlers)
	router.SetupRoutes(r)
	r.Use(rateLimitMiddleware())

	server := &http.Server{
		Addr:    ":" + appConfig.PORT,
		Handler: r,
	}

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.GetLogger().Info("Server is shutting down...")

		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.GetLogger().Fatal("Server shutdown error:", err)
		}

		logger.GetLogger().Info("Server has gracefully stopped")
		os.Exit(0)
	}()

	logger.GetLogger().Info("Server is running on :" + appConfig.PORT)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetLogger().Fatal("Error starting server:", err)
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	limiter := time.Tick(time.Second)

	return func(c *gin.Context) {
		select {
		case <-limiter:
			c.Next()
		default:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
		}
	}
}
