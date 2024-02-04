package routers

import (
	"aitu-funpage/backend/internal/rest/handlers"
	"aitu-funpage/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	authHandlers handlers.AuthHandlers
	newsHandlers handlers.NewsHandlers
}

func NewRouters(authHandlers handlers.AuthHandlers, newsHandlers handlers.NewsHandlers) *Routers {
	return &Routers{
		authHandlers: authHandlers,
		newsHandlers: newsHandlers,
	}
}

func (r *Routers) SetupRoutes(app *gin.Engine) {
	appRouter := app.Group("/app")
	{
		authRouter := appRouter.Group("/auth")
		{
			authRouter.POST("/register", r.authHandlers.Register)
			authRouter.POST("/login", r.authHandlers.Login)
			authRouter.POST("/logout", middleware.RequireAuthMiddleware, r.authHandlers.Logout)
			authRouter.DELETE("/deleteAccount", middleware.RequireAuthMiddleware, r.authHandlers.DeleteAccount)
		}
		newsRouter := appRouter.Group("/news")
		{
			newsRouter.POST("/create", middleware.RequireAuthMiddleware, r.newsHandlers.CreateNews)
			newsRouter.GET("/getByID", r.newsHandlers.GetNewsByObjectID)
			newsRouter.GET("/getByAuthor", r.newsHandlers.GetAllNewsByAuthor)
			newsRouter.GET("/getByTag", r.newsHandlers.GetAllNewsByTag)
			newsRouter.PUT("/updateByID", middleware.RequireAuthMiddleware, r.newsHandlers.UpdateNewsByObjectID)
			newsRouter.DELETE("/deleteByID", middleware.RequireAuthMiddleware, r.newsHandlers.DeleteNewsByObjectID)
		}
	}
}
