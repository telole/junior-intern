package main

import (
	// "log"

	"backend/routes"
	"backend/configs"

	"github.com/labstack/echo/v4"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)

func main() {

	configs.ConnectDB()
	db := configs.DB

	e := echo.New()

	routes.InitRoutes(e, db)


	e.Logger.Fatal(e.Start(":3000"))
}
