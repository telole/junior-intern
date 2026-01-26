package auth

import (
	"backend/models"
	"backend/res/request"
	"backend/res/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct { 
	DB *gorm.DB
}

func (a *AuthController) Register(c echo.Context) error {
	req := new(request.RegisterRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid request"})
	}

	var tenant models.Tenant
	if err := a.DB.Where("id = ?", req.TenantID).First(&tenant).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "tenant not found"})
	}

	var existing models.User
	if err := a.DB.Where("email = ? AND tenant_id = ?", req.Email, req.TenantID).First(&existing).Error; err == nil {
		return c.JSON(400, echo.Map{"error": "Email already exists"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "failed to hash password"})
	}

	user := models.User{ 
		TenantID: tenant.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		IsActive: true,
	}

	if err := a.DB.Create(&user).Error; err != nil { 
		return c.JSON(500, echo.Map{"error": "failed to create user"})
	}

	var roles []models.Role
	a.DB.Table("roles").
		Select("roles.*").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", user.ID).
		Find(&roles)

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Email, roleNames)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "failed to generate token"})
	}

	user.Password = ""

	return c.JSON(201, echo.Map{
		"message": "user created successfully", 
		"token":   token,
		"user":    user,
	})
}

func (a *AuthController) Login(c echo.Context) error {
	req := new(request.LoginRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid request"})
	}

	var user models.User
	if err := a.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(400, echo.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid credentials"})
	}

	var roles []models.Role
	a.DB.Table("roles").
		Select("roles.*").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", user.ID).
		Find(&roles)

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Email, roleNames)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "failed to generate token"})
	}

	user.Password = ""

	return c.JSON(200, echo.Map{
		"message": "login successful", 
		"token":   token,
		"user":    user,
	})
}

func (a *AuthController) Logout(c echo.Context) error {
	return c.JSON(200, echo.Map{"message": "logout successful"})
}