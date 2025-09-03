package models

import (
	"gorm.io/gorm"
)

type FunctionApp struct {
	gorm.Model
	Name string `gorm:"unique"`
	Path string
	Size int64
}
