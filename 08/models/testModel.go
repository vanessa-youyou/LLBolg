package models

type UserInfo struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" gorm:"not null"`
	Gender       string `json:"gender" gorm:"not null"`
	Password     string `json:"password" gorm:"not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`
	Introduction string `json:"introduction"`
	Label        string `json:"label"`
}
