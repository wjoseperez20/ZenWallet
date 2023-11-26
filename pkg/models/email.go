package models

type Email struct {
	Client  string
	Email   string
	Message string
}

type EmailTemplateData struct {
	Username     string
	Message      string
	Transactions []Transaction
}
