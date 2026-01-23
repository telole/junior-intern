package tickets

import (
	"backend/models"
	"backend/res/request"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TicketController struct {
	DB *gorm.DB
}

func (t*TicketController) CreateTicket(c echo.Context) error { 
	req := new(request.CreateTicketRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error" : "invalid request"})
	}

	// Get user_id from context
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		fmt.Println("DEBUG CreateTicket: user_id is nil")
		return c.JSON(401, echo.Map{"message" : "Unauthorized: user_id not found"})
	}
	
	var UserID uint
	switch v := userIDValue.(type) {
	case uint:
		UserID = v
	case float64:
		UserID = uint(v)
	case int:
		UserID = uint(v)
	default:
		fmt.Printf("DEBUG CreateTicket: user_id type is %T, value: %v\n", userIDValue, userIDValue)
		return c.JSON(401, echo.Map{"message" : "Unauthorized: invalid user_id type"})
	}

	// Get tenant_id from context
	tenantIDValue := c.Get("tenant_id")
	if tenantIDValue == nil {
		fmt.Println("DEBUG CreateTicket: tenant_id is nil")
		return c.JSON(401, echo.Map{"message" : "Unauthorized: tenant_id not found"})
	}
	
	TenantID, ok := tenantIDValue.(string)
	if !ok {
		fmt.Printf("DEBUG CreateTicket: tenant_id type is %T, value: %v\n", tenantIDValue, tenantIDValue)
		return c.JSON(401, echo.Map{"message" : "Unauthorized: invalid tenant_id type"})
	}

	fmt.Printf("DEBUG CreateTicket: UserID=%v, TenantID=%v\n", UserID, TenantID)
	ticket := models.Ticket{ 
		TenantID: TenantID,
		UserID: UserID,
		Title: req.Title,
		Description: req.Description,
		Status: "open",
	}

	if err := t.DB.Create(&ticket).Error; err != nil {
	fmt.Println("DB ERROR:", err.Error())
	return c.JSON(500, echo.Map{
		"error": err.Error(),
	})
}

	return c.JSON(201, echo.Map{"message" : "ticket created successfully", "status" : ticket})
}   

func (t*TicketController) GetTickets(c echo.Context) error { 
	tenantIDValue := c.Get("tenant_id")
	if tenantIDValue == nil {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: tenant_id not found"})
	}
	
	tenantID, ok := tenantIDValue.(string)
	if !ok {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: invalid tenant_id type"})
	}

	var tickets []models.Ticket

	if err := t.DB.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return c.JSON(500, echo.Map{"error" : "failed to fetch tickets"})
	}

	return c.JSON(200, echo.Map{"tickets" : tickets})
}

func (t*TicketController) GetTicketByID(c echo.Context) error { 
	tenantIDValue := c.Get("tenant_id")
	if tenantIDValue == nil {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: tenant_id not found"})
	}
	
	tenantID, ok := tenantIDValue.(string)
	if !ok {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: invalid tenant_id type"})
	}
	
	ticketID := c.Param("id")

	var ticket models.Ticket

	if err := t.DB.Where("tenant_id = ? AND id = ?", tenantID, ticketID).First(&ticket).Error; err != nil {
		return c.JSON(404, echo.Map{"error" : "ticket not found"})
	}

	return c.JSON(200, echo.Map{"ticket" : ticket})
}

func (t*TicketController) UpdateTicketStatus(c echo.Context) error { 
	req := new(request.UpdateTicketRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error" : "invalid request"})
	}

	tenantIDValue := c.Get("tenant_id")
	if tenantIDValue == nil {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: tenant_id not found"})
	}
	
	tenantID, ok := tenantIDValue.(string)
	if !ok {
		return c.JSON(401, echo.Map{"error" : "Unauthorized: invalid tenant_id type"})
	}
	
	id := c.Param("id")

	if err := t.DB.Model(&models.Ticket{}).Where("id = ? AND tenant_id = ?", id, tenantID).Update("status", req.Status).Error; err != nil {
		return c.JSON(500, echo.Map{"error" : "failed to update ticket status"})
	}
	return c.JSON(200, echo.Map{"message" : "ticket status updated", "status" : req.Status}) 
}