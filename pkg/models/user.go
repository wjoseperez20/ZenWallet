package models

import "time"

type User struct {
	Username  string    `json:"username" binding:"required" gorm:"type:integer;primaryKey;autoIncrement:true"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
