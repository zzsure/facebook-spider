package models

import (
	"github.com/op/go-logging"
	"time"
)

type Model struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var logger = logging.MustGetLogger("models")

func CreateTable() {
	//db.DB.DropTableIfExists(&WeatherInfo{})
	//create := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	//create.CreateTable(&WeatherInfo{})
}

func MigrateTable() {
	//create := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	//create.AutoMigrate(&WeatherInfo{})
}
