package models

import "time"

type Transaction struct {
	ID        uint      `json:"id" gorm:"type:integer;primary_key;autoIncrement:true"`
	Amount    float32   `json:"amount"`
	Date      time.Time `json:"date"`
	Account   uint      `json:"account" gorm:"type:integer;column:account_number"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateTransaction struct {
	Account uint    `json:"account" binding:"required"`
	Date    string  `json:"date" binding:"required"`
	Amount  float32 `json:"amount" binding:"required"`
}

type UpdateTransaction struct {
	Account uint      `json:"account" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
	Amount  float32   `json:"amount" binding:"required"`
}
