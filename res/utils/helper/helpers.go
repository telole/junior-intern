package helper

import "github.com/labstack/echo/v4"

func isUserAdmin(C echo.Context) bool { 
	roles, ok := C.Get("roles").([]string)
	if !ok {
		return false
	}
	
	for _, role := range roles  {
		if role == "admin" || role == "Super Admin" { 
			return true
		}
	}

	return false
}

func isUserSuperAdmin(C echo.Context) bool { 
	roles, ok := C.Get("roles").([]string)

	if !ok {
		return false
	}

	for _, role := range roles  {
		if role == "Super Admin" {
			return true
		}
	}

	return false
}


func hasRole(c echo.Context, RoleName string) bool {
	roles, ok := c.Get("roles").([]string)
	if !ok {
		return false
	}

	for _, role := range roles {
		if role == RoleName {
			return true
		}
	}

	return false
}

func getUserId( c echo.Context) (uint, bool) { 
	userID, ok := c.Get("user_id").(uint)
	return  userID, ok
}

func getTenantId(c echo.Context) (string, bool) { 
	TenanId, ok := c.Get("tenant_id").(string)
	return TenanId, ok
}


func getRoles(c echo.Context) ([]string, bool) { 
	roles, ok := c.Get("roles").([]string)
	return roles, ok
}