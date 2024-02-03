package routers

import (
	"aitu-funpage/backend/internal/rest/handlers"
	"aitu-funpage/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	authHandlers handlers.AuthHandlers
}

func NewRouters(authHandlers handlers.AuthHandlers) *Routers {
	return &Routers{
		authHandlers: authHandlers,
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
	}
}
