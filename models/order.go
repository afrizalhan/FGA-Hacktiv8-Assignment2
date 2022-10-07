package models

import "time"

type Order struct {
	OrderID uint `gorm:"primary_key"`
	CustomerName string `gorm:"not null;type:varchar(191)"`
	OrderedAt time.Time
	Items []Item
}