package main

import (
	"gorm.io/gorm"
	"time"
)

type Item struct {
	Name  string `gorm:"type:varchar(100);not null"`
	Price int
	URL   string `gorm:"type:varchar(100);uniqueIndex"`
}

type LatestItem struct {
	Item
	CreatedAt time.Time
}

type ItemMaster struct {
	gorm.Model
	Item
}

func (ItemMaster) TableName() string {
	return "item_master"
}
