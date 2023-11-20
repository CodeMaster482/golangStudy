package main

import (
	"go-web-bmstu/internal/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	env := dsn.FromEnv()

	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(); err != nil {
		panic("cant migrate db:" + err.Error())
	}
}
