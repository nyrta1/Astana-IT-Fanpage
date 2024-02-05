package routers

import (
	"aitu-funpage/backend/internal/rest/handlers"
	"aitu-funpage/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	authHandlers    handlers.AuthHandlers
	newsHandlers    handlers.NewsHandlers
	tagHandlers     handlers.TagHandlers
	commentHandlers handlers.CommentHandlers
}

func NewRouters(authHandlers handlers.AuthHandlers, newsHandlers handlers.NewsHandlers, tagHandlers handlers.TagHandlers, commentHandlers handlers.CommentHandlers) *Routers {
	return &Routers{
		authHandlers:    authHandlers,
		newsHandlers:    newsHandlers,
		tagHandlers:     tagHandlers,
		commentHandlers: commentHandlers,
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
			newsRouter.PUT("/updateByID", middleware.RequireAuthMiddleware, r.newsHandlers.UpdateNewsByObjectID)
			newsRouter.DELETE("/deleteByID", middleware.RequireAuthMiddleware, r.newsHandlers.DeleteNewsByObjectID)
		}
		tagRouter := appRouter.Group("/tag")
		{
			tagRouter.POST("/add", middleware.RequireAuthMiddleware, r.tagHandlers.AddTagToNewsByNewsID)
			tagRouter.GET("/get", r.tagHandlers.GetTagsByNewsID)
			tagRouter.PUT("/update", middleware.RequireAuthMiddleware, r.tagHandlers.UpdateTagByNewsID)
			tagRouter.DELETE("/delete", middleware.RequireAuthMiddleware, r.tagHandlers.DeleteTagByNewsID)
		}
		commentRouter := appRouter.Group("/comment")
		{
			commentRouter.POST("/add", middleware.RequireAuthMiddleware, r.commentHandlers.AddCommentToNews)
			commentRouter.GET("/getByNewsID", r.commentHandlers.GetCommentsByNewsID)
			commentRouter.PUT("/updateByNewsID", middleware.RequireAuthMiddleware, r.commentHandlers.UpdateCommentsByNewsID)
			commentRouter.DELETE("/deleteByNewsID", middleware.RequireAuthMiddleware, r.commentHandlers.DeleteCommentByNewsID)
		}
	}
}
