package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

const defaultDBAdd = "www.sunsnasserver.top:3306"

func init() {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		addr = defaultDBAdd
	}
	dsn := fmt.Sprintf("root:abcd1234@tcp(%v)/go_r4?charset=utf8mb3&parseTime=True&loc=Local", addr)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
	} else {
		log.Fatalf("Error occoured when connecting to the data base! error: %e", err)
	}
}
