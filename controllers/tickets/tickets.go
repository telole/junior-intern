package tickets

import (
	"backend/models"
	"backend/res/request"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TicketController struct {
	DB *gorm.DB
}

// CreateTicket - Semua user bisa create ticket
func (t *TicketController) CreateTicket(c echo.Context) error {
	req := new(request.CreateTicketRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid request"})
	}
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid user_id"})
	}

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid tenant_id"})
	}

	ticket := models.Ticket{
		TenantID:    tenantID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
	}

	if err := t.DB.Create(&ticket).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to create ticket"})
	}

	return c.JSON(201, echo.Map{
		"message": "ticket created successfully",
		"ticket":  ticket,
	})
}

func (t *TicketController) GetTickets(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid tenant_id"})
	}

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid user_id"})
	}

	isAdmin := t.isUserAdmin(c)

	var tickets []models.Ticket
	query := t.DB.Where("tenant_id = ?", tenantID)

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&tickets).Error; err != nil {
		return c.JSON(500, echo.Map{"error": "failed to fetch tickets"})
	}

	return c.JSON(200, echo.Map{
		"tickets": tickets,
		"count":   len(tickets),
	})
}

func (t *TicketController) GetTicketByID(c echo.Context) error {
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid tenant_id"})
	}

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid user_id"})
	}

	ticketID := c.Param("id")

	isAdmin := t.isUserAdmin(c)

	var ticket models.Ticket
	query := t.DB.Where("tenant_id = ? AND id = ?", tenantID, ticketID)

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(404, echo.Map{"error": "ticket not found or access denied"})
		}
		return c.JSON(500, echo.Map{"error": "failed to fetch ticket"})
	}

	return c.JSON(200, echo.Map{"ticket": ticket})
}

func (t *TicketController) UpdateTicketStatus(c echo.Context) error {
	req := new(request.UpdateTicketRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid request"})
	}

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok {
		return c.JSON(401, echo.Map{"error": "Unauthorized: invalid tenant_id"})
	}

	ticketID := c.Param("id")

	validStatuses := []string{"open", "in_progress", "resolved", "closed"}
	isValidStatus := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return c.JSON(400, echo.Map{"error": "invalid status value"})
	}

	result := t.DB.Model(&models.Ticket{}).
		Where("id = ? AND tenant_id = ?", ticketID, tenantID).
		Update("status", req.Status)

	if result.Error != nil {
		return c.JSON(500, echo.Map{"error": "failed to update ticket status"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(404, echo.Map{"error": "ticket not found"})
	}

	return c.JSON(200, echo.Map{
		"message": "ticket status updated successfully",
		"status":  req.Status,
	})
}

func (t *TicketController) isUserAdmin(c echo.Context) bool {
	roles, ok := c.Get("roles").([]string)
	if !ok {
		return false
	}

	for _, role := range roles {
		if role == "admin" || role == "Super Admin" {
			return true
		}
	}

	return false
}