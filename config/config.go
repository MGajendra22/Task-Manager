package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DataBaseConfig() {
	var err error
	DB, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db?parseTime=true")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	fmt.Println("Successfully connected to MySQL!")
}
