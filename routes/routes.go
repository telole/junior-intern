package routes

import ( 
	"gorm.io/gorm"
	"backend/controllers/auth"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) { 
	authController := &auth.AuthController{DB: db}
	auth := e.Group("/api/auth")
	auth.POST("/register", authController.Register)
}