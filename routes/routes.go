package routes

import (
	"backend/controllers/auth"
	"backend/controllers/profile"
	"backend/controllers/tickets"
	"backend/res/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	authController := &auth.AuthController{DB: db}
	profileController := &profile.ProfileController{DB: db}
	ticketController := &tickets.TicketController{DB: db}

	authGroup := e.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)

	authProtected := e.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware)
	authProtected.POST("/logout", authController.Logout)

	api := e.Group("/api")
	api.Use(middleware.AuthMiddleware)

	api.GET("/me", profileController.GetProfile)

	tickets := api.Group("/tickets")
	tickets.POST("", ticketController.CreateTicket)
	tickets.GET("", ticketController.GetTickets)
	tickets.GET("/:id", ticketController.GetTicketByID)
	tickets.PUT("/:id/status", ticketController.UpdateTicketStatus, middleware.AdminMiddleware)

	admin := e.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.AdminMiddleware)

	superAdmin := e.Group("/api/superadmin")
	superAdmin.Use(middleware.AuthMiddleware)
	superAdmin.Use(middleware.SuperAdminMiddleware)
}