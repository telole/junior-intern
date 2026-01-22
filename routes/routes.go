package routes

import (
	"backend/controllers/auth"
	"backend/controllers/profile"
	"backend/res/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	authController := &auth.AuthController{DB: db}

	auth := e.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	authProtected := e.Group("/auth", middleware.AuthMiddleware)
	authProtected.POST("/logout", authController.Logout)

	profileController := &profile.ProfileController{DB: db}

	api := e.Group("/api", middleware.AuthMiddleware)
	api.GET("/me", profileController.GetProfile)
}
	