package main

import (
	"main/internal/ds"
	"main/internal/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	env, err := dsn.FromEnv()
	if err != nil {
		panic("Error from reading env")
	}

	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(
		&ds.User{},
		&ds.Operation{},
		&ds.Banknote{},
		&ds.OperationBanknote{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}

	users := []ds.User{
		{Login: "user1", Password: "password1", IsModerator: true},
		{Login: "user2", Password: "password2", IsModerator: false},
		{Login: "user3", Password: "password3", IsModerator: false},
	}

	money := []ds.Banknote{
		{Nominal: 10, ImageURL: "/static/resources/ten.jpg"},
		{Nominal: 50, ImageURL: "/static/resources/fifty.jpg"},
		{Nominal: 100, ImageURL: "/static/resources/handred.jpg"},
		{Nominal: 500, ImageURL: "/static/resources/fivehandred.jpg"},
		{Nominal: 1000, ImageURL: "/static/resources/thouthend.jpg"},
		{Nominal: 2000, ImageURL: "/static/resources/twothouthend.jpg"},
		{Nominal: 5000, ImageURL: "/static/resources/fivethouthend.jpg"},
	}
	db.Create(&users)
	db.Create(&money)
}
