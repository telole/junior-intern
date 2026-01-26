package middleware

import (
	"github.com/labstack/echo/v4"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role")
		if role == nil {
			return c.JSON(401, echo.Map{"error": "Unauthorized: role not found"})
		}

		roles, ok := role.([]string)
		if !ok {
			return c.JSON(401, echo.Map{"error": "Unauthorized: invalid role type"})
		}

		isAdmin := false
		for _, r := range roles {
			if r == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return c.JSON(401, echo.Map{"error": "Unauthorized: admin access required"})
		}

		return next(c)
	}
}

func SuperAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc { 
	return func(c echo.Context) error { 
		roles, ok := c.Get("roles").([]string)
		if !ok { 
			return c.JSON(403, echo.Map{"error" : "Forbidden: invalid role type"})
		}

		tenantID, ok := c.Get("tenant_id").(string)

		if !ok {
			return c.JSON(403, echo.Map{"error" : "Forbidden: invalid user role"})
		}

		if tenantID != "00000000-0000-0000-0000-000000000000" {
			return c.JSON(403, echo.Map{"error" : "Forbidden: not a super admin"})
		}
		isSuperAdmin := false
		for _, r := range roles { 
			if r == "super_admin" { 
				isSuperAdmin = true
			}
		}
		if !isSuperAdmin {
			return c.JSON(403, echo.Map{"error" : "Forbidden: SuperAdmin Access"})
		}

		return next(c)
	}
}