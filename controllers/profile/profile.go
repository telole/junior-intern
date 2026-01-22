package profile

import (
	"backend/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileController struct{
	DB *gorm.DB
}

func (p *ProfileController) GetProfile(c echo.Context) error {
	var user models.User
	if err := p.DB.Where("id = ?", c.Get("user_id")).First(&user).Error; err != nil {
		return c.JSON(404, echo.Map{"error": "user not found"})
	}

	return c.JSON(200, echo.Map{"message": "profile fetched successfully", "user": user})
}