package main

import (
	"main/internal/app/ds"
	"main/internal/app/dsn"

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
		{Login: "user1", Password: "password1", IsModerator: true, Balance: 9999},
		{Login: "user2", Password: "password2", IsModerator: false, Balance: 1000},
		{Login: "user3", Password: "password3", IsModerator: false, Balance: 90},
	}

	money := []ds.Banknote{
		{Nominal: 10, ImageURL: "/static/resources/ten.jpg", Status: "действует"},
		{Nominal: 50, ImageURL: "/static/resources/fifty.jpg", Status: "действует"},
		{Nominal: 100, ImageURL: "/static/resources/handred.jpg", Status: "действует"},
		{Nominal: 500, ImageURL: "/static/resources/fivehandred.jpg", Status: "действует"},
		{Nominal: 1000, ImageURL: "/static/resources/thouthend.jpg", Status: "действует"},
		{Nominal: 2000, ImageURL: "/static/resources/twothouthend.jpg", Status: "действует"},
		{Nominal: 5000, ImageURL: "/static/resources/fivethouthend.jpg", Status: "действует"},
	}
	db.Create(&users)
	db.Create(&money)
}
