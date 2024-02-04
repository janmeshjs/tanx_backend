package pkg

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open("postgres", "postgres://postgres:Anmol@8907@localhost:5432/tanxficoins?sslmode=disable")
	if err != nil {
		panic("Failed to connect to database:" + err.Error())
	}

	db.LogMode(true)

	migrateTables()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
func migrateTables() {
	db.AutoMigrate(&PriceAlert{})
}
