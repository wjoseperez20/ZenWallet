package models

type Email struct {
	Email string
}

type EmailTemplateData struct {
	Username     string
	Message      string
	Transactions []Transaction
}
