package models

import "time"

type Account struct {
	ID        uint      `json:"id" gorm:"type:integer;autoIncrement:true"`
	Client    string    `json:"client"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Account   uint      `json:"account"  gorm:"primary_key"`
	Balance   float32   `json:"balance" sql:"type:decimal(10,2);"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateAccount struct {
	Client string `json:"client" binding:"required"`
	Email  string `json:"email" binding:"required"`
}

type UpdateAccount struct {
	Client string `json:"client"`
	Email  string `json:"email"`
}
