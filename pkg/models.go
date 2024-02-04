package pkg

import "github.com/jinzhu/gorm"

type PriceAlert struct {
	gorm.Model
	CoinID      string  `json:"coin_id" gorm:"not null"`
	TargetPrice float64 `json:"target_price" gorm:"not null"`
	UserID      uint    `json:"user_id" gorm:"not null"`
	Status      string  `json:"status" gorm:"default:'created'"`
}
