package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string
	FunctionApps []FunctionApp `gorm:"foreignKey:AccountID; constraint:OnDelete:CASCADE"`
}

type FunctionApp struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Path      string
	Size      int64
	AccountID uint
}
