package models

import (
	"gorm.io/gorm"
)

type FunctionApp struct {
	gorm.Model
	name string `gorm:"unique"`
	path string
}
