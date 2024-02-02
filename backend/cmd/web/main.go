package main

import (
	"aitu-funpage/backend/internal/config"
	"aitu-funpage/backend/internal/db"
	"aitu-funpage/backend/pkg/logger"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfig config.App

func init() {
	if err := godotenv.Load("backend/.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
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

func main() {
	logger.InitLogger()

	appConfig = config.App{
		PORT: os.Getenv("APP_PORT"),
		DB:   initializeDB(),
	}

	// TODO: replace _ to db
	_, err := db.GetDBInstance(appConfig.DB)
	if err != nil {
		logger.GetLogger().Fatal("Error initializing DB:", err)
	}

	server := &http.Server{
		Addr: ":" + appConfig.PORT,
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
