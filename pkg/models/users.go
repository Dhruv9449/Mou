package models

type User struct {
	ID      uint   `gorm:"primary_key"`
	Name    string `gorm:"not null"`
	Email   string `gorm:"unique_index"`
	Role    string `gorm:"default:normal"`
	Picture string `gorm:"not null"`
}
