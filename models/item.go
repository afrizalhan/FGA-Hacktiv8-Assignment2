package models

import (
	// "errors"
	// "fmt"
	// "time"

	// "gorm.io/gorm"
)

type Item struct {
	ItemID uint `gorm:"primary_key"`
	ItemCode string `gorm:"not null;unique;type:varchar(191)"`
	Description string `gorm:"not null;type:varchar(191)"`
	Quantity int `gorm:"not null;type:int"`
	OrderID uint
}

// func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
// 	fmt.Println("Product before create()")

// 	if len(p.Name) < 4 {
// 		err = errors.New("Product name is too short")
// 	}

// 	return
// }