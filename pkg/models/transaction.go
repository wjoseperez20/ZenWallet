package models

import "time"

type Transaction struct {
	ID        int       `json:"id" gorm:"type:integer;primary_key;autoIncrement:true"`
	Amount    float32   `json:"amount" sql:"type:decimal(10,2);"`
	Date      time.Time `json:"date"`
	Account   int       `json:"account" gorm:"type:integer;column:account_id;references:accounts(account)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateTransaction struct {
	Account int     `json:"account" binding:"required"`
	Date    string  `json:"date" binding:"required"`
	Amount  float32 `json:"amount" binding:"required" sql:"type:decimal(10,2);"`
}

type UpdateTransaction struct {
	Account int     `json:"account" binding:"required"`
	Date    string  `json:"date" binding:"required"`
	Amount  float32 `json:"amount" binding:"required" sql:"type:decimal(10,2);"`
}
