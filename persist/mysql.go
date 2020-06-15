package persist

import (
	"github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func InitDB() bool {
	var err error
	dataSourceName := "root:12345@tcp(localhost:3306)/business?parseTime=True"
	Db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		log4go.Error(err)
		return false
	}
	log4go.Info("Connect to Mysql!")
	Db.LogMode(true)
	Db.AutoMigrate(&User{})
	return true
}
