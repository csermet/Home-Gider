package models

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:100;not null"`
	Icon string `json:"icon" gorm:"size:50"`
}
