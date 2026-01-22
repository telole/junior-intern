package models

import "time"


type User struct { 
	ID uint `gorm:"primarykey"`
	TenantID string `gorm:"type:uuid;not null"`
	Name string
	Email string
	Password string `json:"-"`
	IsActive bool 
	CreatedAt time.Time
	
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID"`
}