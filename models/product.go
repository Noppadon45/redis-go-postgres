package models

type Product struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	UserID uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
